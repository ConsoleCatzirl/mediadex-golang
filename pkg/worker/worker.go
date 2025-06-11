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

	moviesWalkers []walker.Walker
	musicWalkers  []walker.Walker
	seriesWalkers []walker.Walker

	prePipeMovies chan *item.FileItem
	prePipeMusic  chan *item.FileItem
	prePipeSeries chan *item.FileItem

	movieRunners  []runner.Runner
	musicRunners  []runner.Runner
	seriesRunners []runner.Runner

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

		prePipeMovies: make(chan *item.FileItem, pipeSize),
		prePipeMusic:  make(chan *item.FileItem, pipeSize),
		prePipeSeries: make(chan *item.FileItem, pipeSize),

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
	// add movies walkers
	if len(config.Paths.Movies) > 0 {
		log.Println("Trace: creating movie walkers")
		newMoviesWalkers := walker.MakeWalkers(
			config.Paths.Movies,
			&config.Actions,
			newWorker.prePipeMovies,
		)
		for _, walker := range newMoviesWalkers {
			newWorker.moviesWalkers = append(newWorker.moviesWalkers, walker)
		}

		// add movies runners
		log.Println("Trace: creating movie runners")

		newMovieRunners := make([]runner.Runner, 0)
		for _ = range runnerCount {
			newRunner := runner.NewRunner(
				&config.Actions,
				item.MovieFamily,
				newWorker.prePipeMovies,
				newWorker.postPipeArango,
				newWorker.postPipeOpenSearch,
				newWorker.arangoBackend,
				newWorker.openSearchBackend,
			)
			newMovieRunners = append(newMovieRunners, newRunner)
		}
		newWorker.movieRunners = newMovieRunners
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

	// add series walkers
	if len(config.Paths.Series) > 0 {
		log.Println("Trace: creating series walkers")
		newSeriesWalkers := walker.MakeWalkers(
			config.Paths.Series,
			&config.Actions,
			newWorker.prePipeSeries,
		)
		for _, walker := range newSeriesWalkers {
			newWorker.seriesWalkers = append(newWorker.seriesWalkers, walker)
		}

		// add series runners
		log.Println("Trace: creating series runners")
		newSeriesRunners := make([]runner.Runner, 0)
		for _ = range runnerCount {
			newRunner := runner.NewRunner(
				&config.Actions,
				item.SeriesFamily,
				newWorker.prePipeSeries,
				newWorker.postPipeArango,
				newWorker.postPipeOpenSearch,
				newWorker.arangoBackend,
				newWorker.openSearchBackend,
			)
			newSeriesRunners = append(newSeriesRunners, newRunner)
		}
		newWorker.seriesRunners = newSeriesRunners
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
