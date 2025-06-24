package runner

import (
	"testing"

	"internal/item"
	"pkg/conf"
	"test/fakes"
)

func TestRunner(t *testing.T) {
	fPipe := make(chan *item.FileItem)
	iPipe1 := make(chan item.Item)
	iPipe2 := make(chan item.Item)

	foundBack := &fakes.FakeClient{
		Pipe:      iPipe1,
		HasLookup: true,
	}
	notFoundBack := &fakes.FakeClient{
		Pipe:      iPipe2,
		HasLookup: false,
	}
	errFoundBack := &fakes.FakeClient{
		Pipe:         iPipe1,
		HasLookup:    true,
		HasLookupErr: true,
	}
	errNotFoundBack := &fakes.FakeClient{
		Pipe:         iPipe2,
		HasLookup:    false,
		HasLookupErr: true,
	}

	fakeRunner := mediaRunner{
		action:         &conf.ActionConf{},
		itemFamily:     item.FeatureFamily,
		filePipe:       fPipe,
		arangoPipe:     iPipe1,
		openSearchPipe: iPipe2,
	}

	fItem := &item.FileItem{}

	exist, err := fakeRunner.InBackend(fItem, foundBack)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !exist {
		t.Errorf("Fake Item should exist")
	}

	exist, err = fakeRunner.InBackend(fItem, notFoundBack)
	if err != nil {
		t.Errorf("%v", err)
	}
	if exist {
		t.Errorf("Fake Item should not exist")

	}

	exist, err = fakeRunner.InBackend(fItem, errFoundBack)
	if err == nil {
		t.Errorf("%v", err)
	}
	if !exist {
		t.Errorf("Fake Item should exist")
	}

	exist, err = fakeRunner.InBackend(fItem, errNotFoundBack)
	if err == nil {
		t.Errorf("%v", err)
	}
	if exist {
		t.Errorf("Fake Item should not exist")

	}

}
