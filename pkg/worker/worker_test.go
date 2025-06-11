package worker

import (
	"internal/item"
	"internal/runner"
	"internal/walker"
	"io/fs"
	"math/rand/v2"
	"pkg/conf"
	"testing"
	"time"
)

type fakeItem struct{}

func (f *fakeItem) JsonDoc() (*item.JsonDocument, error) {
	return &item.JsonDocument{}, nil
}

func (f *fakeItem) AddMetadata() error {
	return nil
}

type fakeClient struct {
	pipe chan item.Item
}

func (f *fakeClient) Connect() error {
	return nil
}

func (f *fakeClient) Index() error {
	for {
		_, more := <-f.pipe
		if !more {
			return nil
		}
		r := rand.IntN(100) + 100
		t := time.Duration(r)
		time.Sleep(t * time.Millisecond)
	}
}

func (f *fakeClient) LookupItem(x string) (item.Item, error) {
	return nil, nil
}

type fakeRunner struct {
	inPipe   chan *item.FileItem
	outPipe1 chan item.Item
	outPipe2 chan item.Item
}

func (f *fakeRunner) Run() error {
	for {
		_, more := <-f.inPipe
		if !more {
			return nil
		}
		r := rand.IntN(100) + 100
		t := time.Duration(r)
		time.Sleep(t * time.Millisecond)
		f.outPipe1 <- &fakeItem{}
		f.outPipe2 <- &fakeItem{}
	}
}

func (f *fakeRunner) InArango(fItem *item.FileItem) (bool, error) {
	return false, nil
}

func (f *fakeRunner) InOpenSearch(fItem *item.FileItem) (bool, error) {
	return false, nil
}

type fakeWalker struct {
	pipe chan *item.FileItem
}

func (f *fakeWalker) Walk() error {
	f.ProcessFile("", nil, nil)
	f.ProcessFile("", nil, nil)
	return nil
}

func (f *fakeWalker) ProcessFile(_ string, _ fs.DirEntry, _ error) error {
	f.pipe <- &item.FileItem{}
	return nil
}

type fakeWorker struct {
}

func TestFakeWorker(t *testing.T) {
	fakeArangoConf := &conf.ArangoConf{}
	fakeArangoConf.Auth.User = "username"
	fakeArangoConf.Auth.Pass = "password"

	fakeOpenSearchConf := &conf.OpenSearchConf{}
	fakeOpenSearchConf.Auth.User = "username"
	fakeOpenSearchConf.Auth.Pass = "password"

	fakeBackendConf := conf.BackendConf{
		ArangoDB:   fakeArangoConf,
		OpenSearch: fakeOpenSearchConf,
	}

	fakePathsConf := conf.PathsConf{
		Features: []string{"/foo"},
		Music:    []string{"/bar"},
		Episodes: []string{"/baz"},
	}

	fakeConf := &conf.Conf{
		Backend: fakeBackendConf,
		Paths:   fakePathsConf,
	}

	prePipe1 := make(chan *item.FileItem, 1)
	prePipe2 := make(chan *item.FileItem, 1)
	prePipe3 := make(chan *item.FileItem, 1)

	postPipe1 := make(chan item.Item, 1)
	postPipe2 := make(chan item.Item, 1)

	walkers1 := make([]walker.Walker, 0)
	walkers1 = append(walkers1, &fakeWalker{pipe: prePipe1})

	walkers2 := make([]walker.Walker, 0)
	walkers2 = append(walkers2, &fakeWalker{pipe: prePipe2})
	walkers2 = append(walkers2, &fakeWalker{pipe: prePipe2})

	walkers3 := make([]walker.Walker, 0)
	walkers3 = append(walkers3, &fakeWalker{pipe: prePipe3})
	walkers3 = append(walkers3, &fakeWalker{pipe: prePipe3})
	walkers3 = append(walkers3, &fakeWalker{pipe: prePipe3})

	runners1 := make([]runner.Runner, 0)
	runners1 = append(runners1, &fakeRunner{
		inPipe:   prePipe1,
		outPipe1: postPipe1,
		outPipe2: postPipe2,
	})

	runners2 := make([]runner.Runner, 0)
	runners2 = append(runners2, &fakeRunner{
		inPipe:   prePipe2,
		outPipe1: postPipe1,
		outPipe2: postPipe2,
	})

	runners3 := make([]runner.Runner, 0)
	runners3 = append(runners3, &fakeRunner{
		inPipe:   prePipe3,
		outPipe1: postPipe1,
		outPipe2: postPipe2,
	})

	backend1 := &fakeClient{pipe: postPipe1}
	backend2 := &fakeClient{pipe: postPipe2}

	fakeWorker := &Worker{
		fakeConf,

		walkers1,
		walkers2,
		walkers3,

		prePipe1,
		prePipe2,
		prePipe3,

		runners1,
		runners2,
		runners3,

		postPipe1,
		postPipe2,

		backend1,
		backend2,
	}

	err := fakeWorker.Work()
	if err != nil {
		t.Errorf(err.Error())
	}
}
