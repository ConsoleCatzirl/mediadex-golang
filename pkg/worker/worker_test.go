package worker

import (
	"errors"
	"testing"

	"internal/item"
	"internal/runner"
	"internal/walker"
	"pkg/conf"
	"test"
)

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
	prePipe2 := make(chan *item.FileItem, 2)
	prePipe3 := make(chan *item.FileItem, 3)

	postPipe1 := make(chan item.Item, 1)
	postPipe2 := make(chan item.Item, 2)

	walkers1 := make([]walker.Walker, 0)
	walkers1 = append(walkers1, &test.FakeWalker{Pipe: prePipe1})

	walkers2 := make([]walker.Walker, 0)
	walkers2 = append(walkers2, &test.FakeWalker{Pipe: prePipe2})
	walkers2 = append(walkers2, &test.FakeWalker{Pipe: prePipe2})

	walkers3 := make([]walker.Walker, 0)
	walkers3 = append(walkers3, &test.FakeWalker{Pipe: prePipe3})
	walkers3 = append(walkers3, &test.FakeWalker{Pipe: prePipe3})
	walkers3 = append(walkers3, &test.FakeWalker{Pipe: prePipe3})

	runners1 := make([]runner.Runner, 0)
	runners1 = append(runners1, &test.FakeRunner{
		InPipe:   prePipe1,
		OutPipe1: postPipe1,
		OutPipe2: postPipe2,

		InBackRespBool: true,
		InBackRespErr:  nil,
	})

	runners2 := make([]runner.Runner, 0)
	runners2 = append(runners2, &test.FakeRunner{
		InPipe:   prePipe2,
		OutPipe1: postPipe1,
		OutPipe2: postPipe2,

		InBackRespBool: false,
		InBackRespErr:  nil,
	})

	runners3 := make([]runner.Runner, 0)
	runners3 = append(runners3, &test.FakeRunner{
		InPipe:   prePipe3,
		OutPipe1: postPipe1,
		OutPipe2: postPipe2,

		InBackRespBool: false,
		InBackRespErr:  errors.New("InBackend"),
	})

	backend1 := &test.FakeClient{Pipe: postPipe1}
	backend2 := &test.FakeClient{Pipe: postPipe2}

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
