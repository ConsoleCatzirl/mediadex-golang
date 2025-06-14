package walker

import (
	"internal/item"
	"pkg/conf"
	"testing"
)

func TestMakeWalkers(t *testing.T) {
	paths := []string{"foo", "bar", "baz"}
	walkers := MakeWalkers(
		paths,
		&conf.ActionConf{},
		make(chan *item.FileItem, 2),
	)

	if len(walkers) != len(paths) {
		t.Errorf("Incorrect number of workers")
	}
}
