package app

import (
	"context"
	"cynxhostagent/internal/dependencies/param"
	"log"
)

type App struct {
	Dependencies *Dependencies
	Repos        *Repos
	Usecases     *UseCases
}

func NewApp(ctx context.Context, configPath string) (*App, error) {

	log.Println("Initializing Dependencies")
	dependencies := NewDependencies(configPath)

	logger := dependencies.Logger

	logger.Infoln("Initializing Repositories")
	repos := NewRepos(dependencies)

	logger.Infoln("Initializing Param")
	go param.SetupStaticParam(repos.TblParameter, logger)

	logger.Infoln("Initializing Usecases")
	usecases := NewUseCases(repos, dependencies)

	if !dependencies.Config.App.Debug {
		logger.Infoln("Loading Config")
		newConfig := dependencies.Config.LazyLoadConfig(repos.TblInstance, repos.TblPersistentNode)
		dependencies.Config = &newConfig
	}

	logger.Infoln("App initialized")
	return &App{
		Dependencies: dependencies,
		Repos:        repos,
		Usecases:     usecases,
	}, nil
}
