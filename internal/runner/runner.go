package runner

import (
	"internal/backend"
	"internal/item"
	"log"
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
	log.Printf("Trace: %s runner running", r.itemFamily)

	for {
		fItem, more := <-r.filePipe
		if !more {
			//log.Printf("Trace: file pipe is empty and closed")
			return nil
		}
		log.Printf("Trace: processing %s", fItem.FileName)

		var skipArango, skipOpenSearch bool

		if !r.action.Force {
			//log.Printf("Trace: searching for existing document: %s", fItem.FileName)

			foundArango, err := r.InBackend(fItem, r.arangoBackend)
			if err != nil {
				log.Printf("Error finding item: %s", err.Error())
				log.Printf("Trace: found error and item in ArangoDB: %s", fItem.FileName)
			} else if foundArango {
				//log.Printf("Trace: found item in ArangoDB: %s", fItem.FileName)
				skipArango = true
			}

			foundOpenSearch, err := r.InBackend(fItem, r.openSearchBackend)
			if err != nil {
				log.Printf("Error finding item: %s", err.Error())
			} else if foundOpenSearch {
				//log.Printf("Trace: found item in OpenSearch: %s", fItem.FileName)
				skipOpenSearch = true
			}
		}

		if skipArango && skipOpenSearch {
			log.Printf("Trace: skipping unchanged %s: %s", r.itemFamily, fItem.FileName)
			continue
		}

		var mediaItem item.Item
		switch r.itemFamily {
		case item.FeatureFamily:
			//log.Printf("Trace: creating movie item")
			mediaItem = item.FileToFeature(fItem)
		case item.EpisodeFamily:
			//log.Printf("Trace: creating series item")
			mediaItem = item.FileToEpisode(fItem)
		case item.MusicFamily:
			//log.Printf("Trace: creating music item")
			mediaItem = item.FileToMusic(fItem)
		default:
			log.Printf("Error: unknown item family: '%s'", r.itemFamily)
			continue // skip to next item from pipe
		}

		mediaItem.AddMetadata()

		if !r.action.DryRun {
			if !skipArango && r.arangoBackend != nil {
				//log.Printf("Trace: sending item to arangodb: %s", fItem.FileName)
				r.arangoPipe <- mediaItem
			}

			if !skipOpenSearch && r.openSearchBackend != nil {
				//log.Printf("Trace: sending item to opensearch: %s", fItem.FileName)
				r.openSearchPipe <- mediaItem
			}
		}
	}
}

func (r *mediaRunner) InBackend(fItem *item.FileItem, client backend.Client) (bool, error) {
	if fItem == nil {
		log.Printf("Trace: InBackned: nil item not in backend")
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
		//log.Printf("Trace: InBackned: found item %s: %s", key, path)
		jDoc, err := found.JsonDoc()
		if err != nil {
			log.Printf("Error: InBackned: unable to get json doc: %s", key)
			return false, err
		}
		if jDoc != nil {
			same, err := fItem.SameFile(jDoc.FileStats)
			if err != nil {
				log.Printf("Error: checksum match; file stats mismatch: %s", key)
				return false, err
			}
			if same {
				//log.Printf("Trace: InBackned: file item exists in backend: %s", path)
				return true, nil
			} else {
				log.Printf("Trace: InBackned: file stats mismatch: %s: %s", key, path)
			}
		} else {
			log.Printf("Trace: InBackned: empty json doc: %s", key)
		}
	}
	return false, nil
}
