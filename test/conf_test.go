package test

import (
	_ "embed"
	"errors"
	"pkg/conf"
	"testing"

	"gopkg.in/yaml.v2"
)

//go:embed conf/empty.yaml
var empty []byte

//go:embed conf/magicOnly.yaml
var magicOnly []byte

//go:embed conf/pathsOnly.yaml
var pathsOnly []byte

//go:embed conf/backendsOnly.yaml
var backendsOnly []byte

//go:embed conf/minimal.yaml
var minimal []byte

//go:embed conf/full.yaml
var full []byte

func TestConf(t *testing.T) {

	testInput := []struct {
		name      string
		yamlBytes []byte
		wantErr   bool
	}{
		{
			name:      "emptyConf",
			yamlBytes: empty,
			wantErr:   true,
		}, {
			name:      "magicOnly",
			yamlBytes: magicOnly,
			wantErr:   true,
		}, {
			name:      "pathsOnly",
			yamlBytes: pathsOnly,
			wantErr:   true,
		}, {
			name:      "backendsOnly",
			yamlBytes: backendsOnly,
			wantErr:   true,
		}, {
			name:      "minimalSuccess",
			yamlBytes: minimal,
			wantErr:   false,
		}, {
			name:      "fullExample",
			yamlBytes: full,
			wantErr:   false,
		},
	}

	for _, input := range testInput {
		t.Run(input.name, func(t *testing.T) {
			var err error
			var haveErr bool
			testConf := &conf.MagicConf{}

			err = yaml.Unmarshal(input.yamlBytes, testConf)
			haveErr = (err != nil)
			if !haveErr {
				if testConf.Mediadex != nil {
					err = testConf.Mediadex.Validate()
					haveErr = (err != nil)
				} else {
					haveErr = true
					err = errors.New("Empty configuration")
				}
			}

			if !input.wantErr {
				if haveErr {
					// we have an error we don't want
					t.Logf("%v", err)
					t.Error(err)
				}
			} else {
				if !haveErr {
					// we want an error we don't have
					err = errors.New("Error expected but not found")
					t.Error(err)
				} else {
					t.Logf("%v", err)
				}
			}
		})
	}
}
