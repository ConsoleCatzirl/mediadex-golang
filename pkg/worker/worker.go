package worker

import (
	"errors"
	"sync"

	"internal/backend"
	"internal/item"
	"internal/mlog"
	"internal/runner"
	"internal/walker"
	"pkg/conf"
)

type Worker struct {
	config *conf.Conf

	episodeWalkers []walker.Walker
	featureWalkers []walker.Walker
	musicWalkers   []walker.Walker

	prePipeEpisode chan *item.FileItem
	prePipeFeature chan *item.FileItem
	prePipeMusic   chan *item.FileItem

	episodeRunners []runner.Runner
	featureRunners []runner.Runner
	musicRunners   []runner.Runner

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
	mlog.Trace("worker.MakeWorker", "Valid configuration; creating worker")

	runnerCount := config.Threads.Workers
	filePipeSize := config.Threads.FileBuffer
	backendPipeSize := config.Threads.BackendBuffer

	// create a new worker
	newWorker := &Worker{
		config: config,

		prePipeEpisode: make(chan *item.FileItem, filePipeSize),
		prePipeFeature: make(chan *item.FileItem, filePipeSize),
		prePipeMusic:   make(chan *item.FileItem, filePipeSize),

		postPipeArango:     make(chan item.Item, backendPipeSize),
		postPipeOpenSearch: make(chan item.Item, backendPipeSize),
	}

	// create backends first to pass them to runners

	// add arango backend
	if config.Backend.ArangoDB != nil {
		mlog.Trace("worker.MakeWorker", "Creating ArangoDB client")
		newArangoClient := backend.MakeArangoClient(
			config.Backend.ArangoDB,
			newWorker.postPipeArango,
		)
		err = newArangoClient.Connect()
		if err != nil {
			mlog.Error("ArangoDB client failed to connect", err)
			return nil, err
		} else {
			newWorker.arangoBackend = newArangoClient
		}
	}

	// add opensearch backend
	if config.Backend.OpenSearch != nil {
		mlog.Trace("worker.MakeWorker", "Creating OpenSearch client")
		newOpenSearchClient := backend.MakeOpenSearchClient(
			config.Backend.OpenSearch,
			newWorker.postPipeOpenSearch,
		)
		err = newOpenSearchClient.Connect()
		if err != nil {
			mlog.Error("OpenSearch client failed to connect", err)
			return nil, err
		} else {
			newWorker.openSearchBackend = newOpenSearchClient
		}
	}

	if newWorker.arangoBackend == nil && newWorker.openSearchBackend == nil {
		return nil, errors.New("Error: no backend connected")
	}

	if len(config.Paths.Episodes) > 0 {
		// add episode walkers
		mlog.Trace("worker.MakeWorker", "Creating series walkers")
		newEpisodeWalkers := walker.MakeWalkers(
			config.Paths.Episodes,
			&config.Actions,
			newWorker.prePipeEpisode,
		)
		for _, walker := range newEpisodeWalkers {
			newWorker.episodeWalkers = append(newWorker.episodeWalkers, walker)
		}

		// add episode runners
		mlog.Trace("worker.MakeWorker", "Creating series runners")
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

	if len(config.Paths.Features) > 0 {
		// add feature walkers
		mlog.Trace("worker.MakeWorker", "Creating movie walkers")
		newFeaturesWalkers := walker.MakeWalkers(
			config.Paths.Features,
			&config.Actions,
			newWorker.prePipeFeature,
		)
		for _, walker := range newFeaturesWalkers {
			newWorker.featureWalkers = append(newWorker.featureWalkers, walker)
		}

		// add feature runners
		mlog.Trace("worker.MakeWorker", "Creating movie runners")

		newFeatureRunners := make([]runner.Runner, 0)
		for _ = range runnerCount {
			newRunner := runner.NewRunner(
				&config.Actions,
				item.FeatureFamily,
				newWorker.prePipeFeature,
				newWorker.postPipeArango,
				newWorker.postPipeOpenSearch,
				newWorker.arangoBackend,
				newWorker.openSearchBackend,
			)
			newFeatureRunners = append(newFeatureRunners, newRunner)
		}
		newWorker.featureRunners = newFeatureRunners
	}

	if len(config.Paths.Music) > 0 {
		// add music walkers
		mlog.Trace("worker.MakeWorker", "Creating music walkers")
		newMusicWalkers := walker.MakeWalkers(
			config.Paths.Music,
			&config.Actions,
			newWorker.prePipeMusic,
		)
		for _, walker := range newMusicWalkers {
			newWorker.musicWalkers = append(newWorker.musicWalkers, walker)
		}

		// add music runners
		mlog.Trace("worker.MakeWorker", "Creating music runners")
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

	mlog.Trace("worker.MakeWorker", "Worker created")
	return newWorker, nil
}

func (w *Worker) Work() error {
	//
	// synchronize goroutines
	//

	// todo: if we're cleaning, do it first

	// spawn worker components in reverse: backends, runners, walkers
	// so that the consumers are listening before the producers begin

	mlog.Trace("worker.Worker.Work", "Starting backend indexers")

	var wgArango, wgOpenSearch, wgBackends sync.WaitGroup

	if w.arangoBackend != nil {
		wgArango.Add(1) // todo: conf setting
		wgBackends.Add(1)

		// close backend pipe after indexing finishes
		go func() {
			defer wgArango.Done()
			mlog.Trace("worker.Worker.Work", "ArangoDB index starting")
			w.arangoBackend.Index()
			mlog.Trace("worker.Worker.Work", "ArangoDB index finished")
		}()

		// signal backend complete
		go func() {
			defer wgBackends.Done()
			wgArango.Wait()
		}()
	}

	if w.openSearchBackend != nil {
		wgOpenSearch.Add(1) // todo: conf setting
		wgBackends.Add(1)

		// close backend pipe after indexing finishes
		go func() {
			defer wgOpenSearch.Done()
			mlog.Trace("worker.Worker.Work", "OpenSearch index starting")
			w.openSearchBackend.Index()
			mlog.Trace("worker.Worker.Work", "OpenSearch index finished")
		}()

		// signal backend complete
		go func() {
			defer wgBackends.Done()
			wgOpenSearch.Wait()
		}()
	}

	mlog.Trace("worker.Worker.Work", "Starting runners")

	var wgRunners sync.WaitGroup
	wgRunners.Add(len(w.episodeRunners))
	wgRunners.Add(len(w.featureRunners))
	wgRunners.Add(len(w.musicRunners))

	for _, runner := range w.episodeRunners {
		go func() {
			defer wgRunners.Done()
			//mlog.Trace("worker.Worker.Work", "series runner starting")
			runner.Run()
			//mlog.Trace("worker.Worker.Work", "series runner finished")
		}()
	}

	for _, runner := range w.featureRunners {
		go func() {
			defer wgRunners.Done()
			//mlog.Trace("worker.Worker.Work", "movie runner starting")
			runner.Run()
			//mlog.Trace("worker.Worker.Work", "movie runner finished")
		}()
	}

	for _, runner := range w.musicRunners {
		go func() {
			defer wgRunners.Done()
			//mlog.Trace("worker.Worker.Work", "music runner starting")
			runner.Run()
			//mlog.Trace("worker.Worker.Work", "music runner finished")
		}()
	}

	// close all backend pipes together after all runners finish
	go func() {
		defer close(w.postPipeArango)
		defer close(w.postPipeOpenSearch)
		wgRunners.Wait()
		mlog.Trace("worker.Worker.Work", "All runners finished; closing backend pipes")
	}()

	mlog.Trace("worker.Worker.Work", "Starting file walkers")

	var wgFeatureWalkers, wgMusicWalkers, wgEpisodeWalkers sync.WaitGroup
	wgEpisodeWalkers.Add(len(w.episodeWalkers))
	wgFeatureWalkers.Add(len(w.featureWalkers))
	wgMusicWalkers.Add(len(w.musicWalkers))

	for _, walker := range w.episodeWalkers {
		go func() {
			defer wgEpisodeWalkers.Done()
			//mlog.Trace("worker.Worker.Work", "series walker starting")
			walker.Walk()
			//mlog.Trace("worker.Worker.Work", "series walker finished")
		}()
	}
	for _, walker := range w.featureWalkers {
		go func() {
			defer wgFeatureWalkers.Done()
			//mlog.Trace("worker.Worker.Work", "movie walker starting")
			walker.Walk()
			//mlog.Trace("worker.Worker.Work", "movie walker finished")
		}()
	}
	for _, walker := range w.musicWalkers {
		go func() {
			defer wgMusicWalkers.Done()
			//mlog.Trace("worker.Worker.Work", "music walker starting")
			walker.Walk()
			//mlog.Trace("worker.Worker.Work", "music walker finished")
		}()
	}

	// close each file pipe after each walker type finishes
	go func() {
		defer close(w.prePipeEpisode)
		wgEpisodeWalkers.Wait()
		mlog.Trace("worker.Worker.Work", "All series walkers finished; closing pipe")
	}()
	go func() {
		defer close(w.prePipeFeature)
		wgFeatureWalkers.Wait()
		mlog.Trace("worker.Worker.Work", "All movie walkers finished; closing pipe")

	}()
	go func() {
		defer close(w.prePipeMusic)
		wgMusicWalkers.Wait()
		mlog.Trace("worker.Worker.Work", "All music walkers finished; closing pipe")
	}()

	// finally, wait for backends to finish
	wgBackends.Wait()
	mlog.Trace("worker.Worker.Work", "Backends finished")

	return nil
}
