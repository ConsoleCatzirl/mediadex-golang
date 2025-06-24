package walker

import (
	"io/fs"
	"path/filepath"

	"internal/item"
	"internal/mlog"
	"pkg/conf"
)

type Walker interface {
	Walk() error
	ProcessFile(string, fs.DirEntry, error) error
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
	return filepath.WalkDir(w.searchPath, w.ProcessFile)
}

func (w *mediaWalker) ProcessFile(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		mlog.Error("Could not open file", err)
		return err
	}

	dType := dir.Type()
	if dType.IsDir() {
		// Silently skip directories
		return nil
	}

	if !dType.IsRegular() {
		// Skip special file types
		mlog.Info("Skipping non-regular file", "file", path)
		return nil
	}

	// read file stats from disk
	newFileItem, err := item.StatFile(dir, path, w.searchPath)
	if err != nil {
		mlog.Error("Failure creating file item", err)
		return err
	}

	// send file item to a runner
	w.filePipe <- newFileItem
	return nil
}
