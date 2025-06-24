package tconf

import (
	_ "embed"
	"errors"
	"fmt"
	"testing"

	"internal/mlog"
	"pkg/conf"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

//go:embed yaml/empty.yaml
var empty []byte

//go:embed yaml/magicOnly.yaml
var magicOnly []byte

//go:embed yaml/pathsOnly.yaml
var pathsOnly []byte

//go:embed yaml/backendsOnly.yaml
var backendsOnly []byte

//go:embed yaml/minimal.yaml
var minimal []byte

//go:embed yaml/full.yaml
var full []byte

//go:embed yaml/defaults.yaml
var defaults []byte

func TestConf(t *testing.T) {
	mlog.Verbose()

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

func TestConfDefault(t *testing.T) {
	mlog.Verbose()

	testConf := &conf.MagicConf{}
	err := yaml.Unmarshal(defaults, testConf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if testConf.Mediadex == nil {
		t.Error(errors.New("Empty configuration"))
		t.FailNow()
	}

	unalteredConf := *testConf.Mediadex
	alteredConfPtr := testConf.Mediadex

	alteredConfPtr.AddDefaults()
	diff := cmp.Diff(unalteredConf, *alteredConfPtr)
	if diff != "" {
		msg := fmt.Sprintf("AddDefaults alters conf/defaults.yaml:\n%s", diff)
		t.Error(errors.New(msg))
	}
}
