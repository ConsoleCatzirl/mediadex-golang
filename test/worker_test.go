package test

import (
	_ "embed"

	"pkg/conf"
	"pkg/worker"
	"testing"

	"gopkg.in/yaml.v2"
)

//go:embed conf/full.yaml
var onlineConf []byte

func TestOnlineWorker(t *testing.T) {

	magic := &conf.MagicConf{}
	err := yaml.Unmarshal(onlineConf, magic)
	if err != nil {
		t.Errorf(err.Error())
	}

	onlineWorker, err := worker.MakeWorker(magic.Mediadex)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = onlineWorker.Work()
	if err != nil {
		t.Errorf(err.Error())
	}
}
