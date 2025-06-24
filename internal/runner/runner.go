package runner

import (
	"errors"

	"internal/backend"
	"internal/item"
	"internal/mlog"
	"pkg/conf"
)

type Runner interface {
	Run() error

	InBackend(*item.FileItem, backend.Client) (bool, error)
}

type mediaRunner struct {
	action *conf.ActionConf

	itemFamily item.Family
	filePipe   chan *item.FileItem

	arangoPipe     chan item.Item
	openSearchPipe chan item.Item

	arangoBackend     backend.Client
	openSearchBackend backend.Client
}

func NewRunner(
	aConf *conf.ActionConf,
	family item.Family,
	fPipe chan *item.FileItem,
	aPipe chan item.Item,
	osPipe chan item.Item,
	aBack backend.Client,
	osBack backend.Client,
) Runner {
	return &mediaRunner{
		action:            aConf,
		itemFamily:        family,
		filePipe:          fPipe,
		arangoPipe:        aPipe,
		openSearchPipe:    osPipe,
		arangoBackend:     aBack,
		openSearchBackend: osBack,
	}
}

func (r *mediaRunner) Run() error {
	mlog.Trace("runner.mediaRunner.Run", "Runner running",
		"type", r.itemFamily)

	for {
		fItem, more := <-r.filePipe
		if !more {
			mlog.Trace("runner.mediaRunner.Run", "file pipe is empty and closed",
				"type", r.itemFamily,
			)
			return nil
		}
		var skipArango, skipOpenSearch bool

		if !r.action.Force {
			mlog.Trace("runner.mediaRunner.Run", "searching for existing document",
				"file", fItem.FileName, "type", r.itemFamily,
			)

			foundArango, err := r.InBackend(fItem, r.arangoBackend)
			if err != nil {
				mlog.Error("Failure finding item", err,
					"backend", "arango",
					"file", fItem.FileName,
					"type", r.itemFamily,
				)
			} else if foundArango {
				mlog.Trace("runner.mediaRunner.Run", "Found item",
					"backend", "arango",
					"file", fItem.FileName,
					"type", r.itemFamily,
				)
				skipArango = true
			}

			foundOpenSearch, err := r.InBackend(fItem, r.openSearchBackend)
			if err != nil {
				mlog.Error("Failure finding item", err,
					"backend", "opensearch",
					"file", fItem.FileName,
					"type", r.itemFamily,
				)
			} else if foundOpenSearch {
				mlog.Trace("runner.mediaRunner.Run", "Found item",
					"backend", "opensearch",
					"file", fItem.FileName,
					"type", r.itemFamily,
				)
				skipOpenSearch = true
			}
		}

		if skipArango && skipOpenSearch {
			mlog.Trace("runner.mediaRunner.Run", "Skipping unchanged item",
				"file", fItem.FileName, "type", r.itemFamily,
			)
			continue
		}

		mlog.Info("Processing", "file", fItem.FileName, "type", r.itemFamily)

		var mediaItem item.Item
		switch r.itemFamily {
		case item.EpisodeFamily:
			mlog.Trace("runner.mediaRunner.Run", "Creating series item")
			mediaItem = item.FileToEpisode(fItem)
		case item.FeatureFamily:
			mlog.Trace("runner.mediaRunner.Run", "Creating movie item")
			mediaItem = item.FileToFeature(fItem)
		case item.MusicFamily:
			mlog.Trace("runner.mediaRunner.Run", "Creating music item")
			mediaItem = item.FileToMusic(fItem)
		default:
			mlog.Error(
				"Unknown item type", errors.New("Unknown item type"),
				"file", fItem.FileName, "type", r.itemFamily,
			)
			continue // skip to next item from pipe
		}

		mediaItem.AddMetadata()

		if !r.action.DryRun {
			if !skipArango && r.arangoBackend != nil {
				mlog.Trace("runner.mediaRunner.Run", "Sending item to ArangoDB",
					"file", fItem.FileName, "type", r.itemFamily,
				)
				r.arangoPipe <- mediaItem
			}

			if !skipOpenSearch && r.openSearchBackend != nil {
				mlog.Trace("runner.mediaRunner.Run", "Sending item to OpenSearch",
					"file", fItem.FileName, "type", r.itemFamily,
				)
				r.openSearchPipe <- mediaItem
			}
		}
	}
}

func (r *mediaRunner) InBackend(fItem *item.FileItem, client backend.Client) (bool, error) {
	if fItem == nil {
		mlog.Trace("runner.mediaRunner.InBackend", "nil item, no match")
		return false, nil
	}
	key := fItem.Checksum
	path := fItem.FileName

	found, err := client.LookupItem(key)
	if err != nil && found != nil {
		return true, err
	} else if err != nil {
		return false, err
	}
	if found != nil {
		mlog.Trace(
			"runner.mediaRunner.InBackend", "found item",
			"key", key, "file", path,
		)
		jDoc, err := found.JsonDoc()
		if err != nil {
			mlog.Error("Unable to get json doc", err, "key", key)
			return false, err
		}
		if jDoc != nil {
			same, err := fItem.SameFile(jDoc.FileStats)
			if err != nil {
				mlog.Error(
					"checksum match; file stats mismatch", err,
					"key", key,
				)
				return false, err
			}
			if same {
				mlog.Trace(
					"runner.mediaRunner.InBackend",
					"file item exists in backend",
					"file", path,
				)
				return true, nil
			} else {
				msg := "file stats mismatch"
				mlog.Error(msg, errors.New(msg), "key", key, "file", path)
			}
		} else {
			msg := "empty json doc"
			mlog.Error(msg, errors.New(msg), "key", key)
		}
	}
	return false, nil
}
