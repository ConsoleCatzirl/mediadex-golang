package tworker

import (
	_ "embed"
	"testing"

	"internal/mlog"
	"pkg/conf"
	"pkg/worker"

	"gopkg.in/yaml.v2"
)

//go:embed conf/docker.yaml
var onlineConf []byte

func TestLocalhostWorker(t *testing.T) {
	mlog.Verbose()

	if testing.Short() {
		t.Skip("Skipping online tests")
	}

	magic := &conf.MagicConf{}
	err := yaml.Unmarshal(onlineConf, magic)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	magic.Mediadex.AddDefaults()

	onlineWorker, err := worker.MakeWorker(magic.Mediadex)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}

	err = onlineWorker.Work()
	if err != nil {
		t.Errorf(err.Error())
	}
}
