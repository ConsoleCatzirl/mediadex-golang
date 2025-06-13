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
			log.Printf("Trace: file pipe is empty and closed")
			return nil
		}
		log.Printf("Trace: processing %s", fItem.FullPath)

		var skipArango, skipOpenSearch bool

		if !r.action.Force {
			log.Println("Trace: searching for existing document")

			foundArango, err := r.InBackend(fItem, r.arangoBackend)
			if err != nil {
				log.Printf("Error: %s", err.Error())
			} else if foundArango {
				log.Println("Trace: found item in ArangoDB")
				skipArango = true
			}

			foundOpenSearch, err := r.InBackend(fItem, r.openSearchBackend)
			if err != nil {
				log.Printf("Error: %s", err.Error())
			} else if foundOpenSearch {
				log.Println("Trace: found item in OpenSearch")
				skipOpenSearch = true
			}
		}

		var mediaItem item.Item
		switch r.itemFamily {
		case item.FeatureFamily:
			log.Println("Trace: creating movie item")
			mediaItem = item.FileToFeature(fItem)
		case item.EpisodeFamily:
			log.Println("Trace: creating series item")
			mediaItem = item.FileToEpisode(fItem)
		case item.MusicFamily:
			log.Println("Trace: creating music item")
			mediaItem = item.FileToMusic(fItem)
		default:
			log.Printf("Error: unknown item family: '%s'", r.itemFamily)
			continue // skip to next item from pipe
		}

		mediaItem.AddMetadata()

		if !r.action.DryRun {
			if !skipArango && r.arangoBackend != nil {
				log.Println("Trace: sending item to arangodb")
				r.arangoPipe <- mediaItem
			}

			if !skipOpenSearch && r.openSearchBackend != nil {
				log.Println("Trace: sending item to opensearch")
				r.openSearchPipe <- mediaItem
			}
		}
	}
}

func (r *mediaRunner) InBackend(fItem *item.FileItem, client backend.Client) (bool, error) {
	if fItem == nil {
		return false, nil
	}

	found, err := client.LookupItem(fItem.Checksum)
	if err != nil {
		return false, err
	}
	if found != nil {
		jDoc, err := found.JsonDoc()
		if err != nil {
			return false, err
		}
		if jDoc != nil {
			same, err := fItem.SameFile(jDoc.FileStats)
			if err != nil {
				return false, err
			}
			if same {
				return true, nil
			}
		}
	}
	return false, nil
}
