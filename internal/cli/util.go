package cli

import (
	"errors"
	"io/ioutil"
	"log"
	"os/user"
	"pkg/conf"

	"gopkg.in/yaml.v2"
)

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Panicf("%v", err)
	}
	return usr.HomeDir
}

func openFile(path string, cfg *conf.MagicConf) error {
	log.Printf("Reading configuration from %s", path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Could not read file: %s", path)
		return err
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Printf("Could not unmarshal file: %s", path)
		return err
	}

	return nil
}

func ReadConfigFile(path string) (*conf.Conf, error) {
	var err error

	log.Printf("Reading config file: %s", path)

	newConf := &conf.MagicConf{
		Mediadex: &conf.Conf{},
	}

	if path != "" {
		// If a path is given, it's required.
		log.Printf("Reading configuration from %s", path)
		err = openFile(path, newConf)
		if err != nil {
			log.Printf("Could not read config file: %s", path)
			return nil, err
		}
	} else {
		// If no path is given, try the local directory first,
		// then fall back to the user's home directory.
		localFile := "mediadex.yaml"
		homeDirFile := homeDir() + "/" + localFile

		err = openFile(localFile, newConf)
		if err != nil {
			log.Printf("Could not read config file: %s", localFile)
			err = openFile(homeDirFile, newConf)
			if err != nil {
				log.Printf("Could not read config file: %s", homeDirFile)
				return nil, err
			}
		}
	}

	if newConf == nil || newConf.Mediadex == nil {
		return nil, errors.New("Empty configuration")
	}

	err = newConf.Mediadex.AddDefaults()
	if err != nil {
		return nil, err
	}

	return newConf.Mediadex, nil
}
