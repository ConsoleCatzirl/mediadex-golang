package test

import (
	"errors"
	"io/fs"
	"math/rand/v2"
	"time"

	"internal/backend"
	"internal/item"
)

// Item
type FakeItem struct{}

func (f *FakeItem) JsonDoc() (*item.JsonDocument, error) {
	return &item.JsonDocument{FileStats: &item.FileItem{}}, nil
}

func (f *FakeItem) AddMetadata() error {
	return nil
}

// Backend Client
type FakeClient struct {
	Pipe chan item.Item

	HasLookup    bool
	HasLookupErr bool
	HasUpsertErr bool
}

func (f *FakeClient) Connect() error {
	return nil
}

func (f *FakeClient) Index() {
	for {
		_, more := <-f.Pipe
		if !more {
			break
		}
		r := rand.IntN(100) + 100
		t := time.Duration(r)
		time.Sleep(t * time.Millisecond)
	}
}

func (f *FakeClient) LookupItem(_ string) (item.Item, error) {
	var it item.Item
	var err error
	if f.HasLookupErr {
		err = errors.New("LookupItem")
	}
	if f.HasLookup {
		it = &FakeItem{}
	}
	return it, err
}

func (f *FakeClient) UpsertItem(_ item.Item) error {
	var err error
	if f.HasUpsertErr {
		err = errors.New("UpsertItem")
	}
	return err
}

// Walker
type FakeWalker struct {
	Pipe chan *item.FileItem
}

func (f *FakeWalker) Walk() error {
	f.ProcessFile("foo", nil, nil)
	f.ProcessFile("bar", nil, nil)
	f.ProcessFile("baz", nil, nil)
	return nil
}

func (f *FakeWalker) ProcessFile(_ string, _ fs.DirEntry, _ error) error {
	f.Pipe <- &item.FileItem{}
	return nil
}

// Runner
type FakeRunner struct {
	InPipe   chan *item.FileItem
	OutPipe1 chan item.Item
	OutPipe2 chan item.Item

	InBackRespBool bool
	InBackRespErr  error
}

func (f *FakeRunner) Run() error {
	for {
		_, more := <-f.InPipe
		if !more {
			return nil
		}

		r := rand.IntN(100) + 100
		t := time.Duration(r)
		time.Sleep(t * time.Millisecond)

		f.OutPipe1 <- &FakeItem{}
		f.OutPipe2 <- &FakeItem{}
	}
}

func (f *FakeRunner) InBackend(_ *item.FileItem, _ backend.Client) (bool, error) {
	return f.InBackRespBool, f.InBackRespErr
}
