//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"answer/internal/base/conf"
	"answer/internal/base/data"
	"answer/internal/base/middleware"
	"answer/internal/base/server"
	"answer/internal/base/translator"
	"answer/internal/controller"
	"answer/internal/controller_backyard"
	"answer/internal/repo"
	"answer/internal/router"
	"answer/internal/service"
	"answer/internal/service/service_config"
	"github.com/google/wire"
	"github.com/segmentfault/pacman"
	"github.com/segmentfault/pacman/log"
)

// initApplication init application.
func initApplication(
	debug bool,
	serverConf *conf.Server,
	dbConf *data.Database,
	cacheConf *data.CacheConf,
	i18nConf *translator.I18n,
	swaggerConf *router.SwaggerConfig,
	serviceConf *service_config.ServiceConfig,
	logConf log.Logger) (*pacman.Application, func(), error) {
	panic(wire.Build(
		server.ProviderSetServer,
		router.ProviderSetRouter,
		controller.ProviderSetController,
		controller_backyard.ProviderSetController,
		service.ProviderSetService,
		repo.ProviderSetRepo,
		translator.ProviderSet,
		middleware.ProviderSetMiddleware,
		newApplication,
	))
}
