package test

import (
	_ "embed"

	"pkg/conf"
	"pkg/worker"
	"testing"

	"gopkg.in/yaml.v2"
)

//go:embed conf/docker.yaml
var onlineConf []byte

func TestOnlineWorker(t *testing.T) {

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
