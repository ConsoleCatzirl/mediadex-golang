package worker

import (
	"internal/backend"
	"internal/item"
	"internal/runner"
	"internal/walker"
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

	moviesRunners []runner.Runner
	musicRunners  []runner.Runner
	seriesRunners []runner.Runner

	postPipeArango     chan item.Item
	postPipeOpenSearch chan item.Item

	arangoBackend     backend.Client
	openSearchBackend backend.Client
}

func (w *Worker) Work() error {
	return nil
}

func MakeWorker(cfg *conf.Conf) (*Worker, error) {
	newWorker := &Worker{config: cfg}
	return newWorker, nil
}
