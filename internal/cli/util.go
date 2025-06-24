package cli

import (
	"errors"
	"io/ioutil"
	"os/user"

	"internal/mlog"
	"pkg/conf"

	"test/tconf"
	"test/tworker"

	"gopkg.in/yaml.v2"
)

// reference test module so that `go mod tidy` doesn't
// remove test-only indirect dependencies
type tconf_hax tconf.Hax
type tworker_hax tworker.Hax

func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		mlog.Error("Failure getting user home dir", err)
		return "", err
	}
	return usr.HomeDir, nil
}

func openFile(path string, cfg *conf.MagicConf) error {
	mlog.Trace("cli.openFile", "Reading configuration", "file", path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		mlog.Error("Failure reading configuration", err, "file", path)
		return err
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		mlog.Error("Failure unmarshalling configuration", err, "file", path)
		return err
	}

	return nil
}

func readConfigFile(path string) (*conf.Conf, error) {
	var err error

	mlog.Trace("cli.readConfigFile", "Reading configuration", "file", path)

	newConf := &conf.MagicConf{
		Mediadex: &conf.Conf{},
	}

	if path != "" {
		// If a path is set then it's required.
		err = openFile(path, newConf)
		if err != nil {
			mlog.Error("Could not read configuration", err, "file", path)
			return nil, err
		}
	} else {
		// If no path is given, try the local directory first,
		// then fall back to the user's home directory.
		localFile := "mediadex.yaml"

		err = openFile(localFile, newConf)
		if err != nil {
			mlog.Warn("Could not read configuration", err, "file", localFile)

			hDir, err := homeDir()
			if err != nil {
				mlog.Error("Could not get home directory", err)
				return nil, err
			}

			homeDirFile := hDir + "/" + localFile
			err = openFile(homeDirFile, newConf)
			if err != nil {
				mlog.Error("Could not read configuration", err, "file", homeDirFile)
				return nil, err
			}
		}
	}

	if newConf == nil || newConf.Mediadex == nil {
		return nil, errors.New("Empty configuration")
	}

	return newConf.Mediadex, nil
}
