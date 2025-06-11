package runner

import (
	"internal/backend"
	"internal/item"
	"log"
	"pkg/conf"
)

type Runner interface {
	Run() error

	inArango(*item.FileItem) (bool, error)
	inOpenSearch(*item.FileItem) (bool, error)
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
	var skipArango, skipOpenSearch bool

	if !r.action.Force {
		log.Println("Trace: searching for existing document")
		var fItem *item.FileItem

		foundArango, err := r.inArango(fItem)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else if foundArango {
			// todo: check for unchanged document in each backend
			log.Println("Trace: found item in ArangoDB")
			skipArango = true
		}

		foundOpenSearch, err := r.inOpenSearch(fItem)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else if foundOpenSearch {
			// todo: check for unchanged document in each backend
			log.Println("Trace: found item in OpenSearch")
			skipOpenSearch = true
		}
	}

	// create media item from file item
	log.Println("Trace: parsing more media metadata")

	// index into backends
	log.Println("Trace: sending item to backends")

	if !skipArango && r.arangoBackend != nil {
		log.Println("Trace: sending item to arangodb")
	}

	if !skipOpenSearch && r.openSearchBackend != nil {
		log.Println("Trace: sending item to opensearch")
	}

	return nil
}

func (r *mediaRunner) inArango(fItem *item.FileItem) (bool, error) {
	return false, nil
}

func (r *mediaRunner) inOpenSearch(fItem *item.FileItem) (bool, error) {
	return false, nil
}
