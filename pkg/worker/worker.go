package worker

import (
	"errors"
	"internal/backend"
	"internal/item"
	"internal/runner"
	"internal/walker"
	"log"
	"pkg/conf"
)

type Worker struct {
	config *conf.Conf

	featureWalkers []walker.Walker
	musicWalkers   []walker.Walker
	episodeWalkers []walker.Walker

	prePipeFeatures chan *item.FileItem
	prePipeMusic    chan *item.FileItem
	prePipeEpisode  chan *item.FileItem

	featureRunners []runner.Runner
	musicRunners   []runner.Runner
	episodeRunners []runner.Runner

	postPipeArango     chan item.Item
	postPipeOpenSearch chan item.Item

	arangoBackend     backend.Client
	openSearchBackend backend.Client
}

func MakeWorker(config *conf.Conf) (*Worker, error) {
	if config == nil {
		return nil, errors.New("Nil configuration")
	}

	err := config.Validate()
	if err != nil {
		return nil, err
	}
	log.Println("Trace: valid configuration; creating worker")

	pipeSize := 64   // channel buffer size
	runnerCount := 2 // todo: conf setting

	// create a new worker
	newWorker := &Worker{
		config: config,

		prePipeFeatures: make(chan *item.FileItem, pipeSize),
		prePipeMusic:    make(chan *item.FileItem, pipeSize),
		prePipeEpisode:  make(chan *item.FileItem, pipeSize),

		postPipeArango:     make(chan item.Item, pipeSize),
		postPipeOpenSearch: make(chan item.Item, pipeSize),
	}

	// create backends first to pass them to runners

	// add arango backend
	if config.Backend.ArangoDB != nil {
		log.Println("Trace: creating ArangoDB client")
		newArangoClient := backend.MakeArangoClient(
			config.Backend.ArangoDB,
			newWorker.postPipeArango,
		)
		newWorker.arangoBackend = newArangoClient
	}

	// add opensearch backend
	if config.Backend.OpenSearch != nil {
		log.Println("Trace: creating OpenSearch client")
		newOpenSearchClient := backend.MakeOpenSearchClient(
			config.Backend.OpenSearch,
			newWorker.postPipeOpenSearch,
		)
		newWorker.openSearchBackend = newOpenSearchClient
	}
	// add feature walkers
	if len(config.Paths.Features) > 0 {
		log.Println("Trace: creating movie walkers")
		newFeaturesWalkers := walker.MakeWalkers(
			config.Paths.Features,
			&config.Actions,
			newWorker.prePipeFeatures,
		)
		for _, walker := range newFeaturesWalkers {
			newWorker.featureWalkers = append(newWorker.featureWalkers, walker)
		}

		// add feature runners
		log.Println("Trace: creating movie runners")

		newFeatureRunners := make([]runner.Runner, 0)
		for _ = range runnerCount {
			newRunner := runner.NewRunner(
				&config.Actions,
				item.FeatureFamily,
				newWorker.prePipeFeatures,
				newWorker.postPipeArango,
				newWorker.postPipeOpenSearch,
				newWorker.arangoBackend,
				newWorker.openSearchBackend,
			)
			newFeatureRunners = append(newFeatureRunners, newRunner)
		}
		newWorker.featureRunners = newFeatureRunners
	}

	// add music walkers
	if len(config.Paths.Music) > 0 {
		log.Println("Trace: creating music walkers")
		newMusicWalkers := walker.MakeWalkers(
			config.Paths.Music,
			&config.Actions,
			newWorker.prePipeMusic,
		)
		for _, walker := range newMusicWalkers {
			newWorker.musicWalkers = append(newWorker.musicWalkers, walker)
		}

		// add music runners
		log.Println("Trace: creating music runners")
		newMusicRunners := make([]runner.Runner, 0)
		for _ = range runnerCount {
			newRunner := runner.NewRunner(
				&config.Actions,
				item.MusicFamily,
				newWorker.prePipeMusic,
				newWorker.postPipeArango,
				newWorker.postPipeOpenSearch,
				newWorker.arangoBackend,
				newWorker.openSearchBackend,
			)
			newMusicRunners = append(newMusicRunners, newRunner)
		}
		newWorker.musicRunners = newMusicRunners
	}

	// add episode walkers
	if len(config.Paths.Episodes) > 0 {
		log.Println("Trace: creating series walkers")
		newEpisodeWalkers := walker.MakeWalkers(
			config.Paths.Episodes,
			&config.Actions,
			newWorker.prePipeEpisode,
		)
		for _, walker := range newEpisodeWalkers {
			newWorker.episodeWalkers = append(newWorker.episodeWalkers, walker)
		}

		// add episode runners
		log.Println("Trace: creating series runners")
		newEpisodeRunners := make([]runner.Runner, 0)
		for _ = range runnerCount {
			newRunner := runner.NewRunner(
				&config.Actions,
				item.EpisodeFamily,
				newWorker.prePipeEpisode,
				newWorker.postPipeArango,
				newWorker.postPipeOpenSearch,
				newWorker.arangoBackend,
				newWorker.openSearchBackend,
			)
			newEpisodeRunners = append(newEpisodeRunners, newRunner)
		}
		newWorker.episodeRunners = newEpisodeRunners
	}

	log.Println("Trace: new worker created")
	return newWorker, nil
}

func (w *Worker) Work() error {
	//
	// synchronize goroutines
	//

	// todo: if we're cleaning, do it first

	log.Println("Trace: starting backend indexers")

	log.Println("Trace: starting core runners")

	log.Println("Trace: core runners finished")

	log.Println("Trace: starting file walkers")

	log.Println("Trace: file walkers finished")

	log.Println("Trace: runners finished")

	return nil
}
