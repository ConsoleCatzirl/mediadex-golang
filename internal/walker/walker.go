package walker

import (
	"internal/item"
	"io/fs"
	"path/filepath"
	"pkg/conf"
)

type Walker interface {
	Walk() error
	processFile(string, fs.DirEntry, error) error
}

type mediaWalker struct {
	action     *conf.ActionConf
	searchPath string
	filePipe   chan *item.FileItem
}

func MakeWalkers(paths []string, aConf *conf.ActionConf, fPipe chan *item.FileItem) []Walker {
	newWalkers := make([]Walker, 0)
	for _, path := range paths {
		nw := &mediaWalker{
			action:     aConf,
			searchPath: path,
			filePipe:   fPipe,
		}
		newWalkers = append(newWalkers, nw)
	}
	return newWalkers
}

func (w *mediaWalker) Walk() error {
	return filepath.WalkDir(w.searchPath, w.processFile)
}

func (w *mediaWalker) processFile(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// todo: skip non-regular files

	// todo: create item.FileItem

	// todo: send file item down the pipe

	return nil
}
