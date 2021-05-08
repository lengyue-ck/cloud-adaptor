// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"goodrain.com/cloud-adaptor/cmd/cloud-adaptor/config"
	"goodrain.com/cloud-adaptor/internal/handler"
	"goodrain.com/cloud-adaptor/internal/middleware"
	"goodrain.com/cloud-adaptor/internal/nsqc/producer"
	"goodrain.com/cloud-adaptor/internal/repo"
	"goodrain.com/cloud-adaptor/internal/repo/appstore"
	"goodrain.com/cloud-adaptor/internal/repo/dao"
	"goodrain.com/cloud-adaptor/internal/task"
	"goodrain.com/cloud-adaptor/internal/types"
	"goodrain.com/cloud-adaptor/internal/usecase"
	"gorm.io/gorm"
)

import (
	_ "github.com/helm/helm/pkg/repo"
	_ "k8s.io/helm/pkg/proto/hapi/chart"
)

// Injectors from wire.go:

// initApp init the application.
func initApp(contextContext context.Context, db *gorm.DB, configConfig *config.Config, arg chan types.KubernetesConfigMessage, arg2 chan types.InitRainbondConfigMessage, arg3 chan types.UpdateKubernetesConfigMessage) (*gin.Engine, error) {
	appStoreDao := dao.NewAppStoreDao(db)
	appTemplater := appstore.NewAppTemplater()
	storer := appstore.NewStorer(appTemplater)
	appStoreRepo := repo.NewAppStoreRepo(appStoreDao, storer, appTemplater)
	middlewareMiddleware := middleware.NewMiddleware(appStoreRepo)
	taskProducer := producer.NewTaskChannelProducer(arg, arg2, arg3)
	cloudAccesskeyRepository := repo.NewCloudAccessKeyRepo(db)
	createKubernetesTaskRepository := repo.NewCreateKubernetesTaskRepo(db)
	initRainbondTaskRepository := repo.NewInitRainbondRegionTaskRepo(db)
	updateKubernetesTaskRepository := repo.NewUpdateKubernetesTaskRepo(db)
	taskEventRepository := repo.NewTaskEventRepo(db)
	rainbondClusterConfigRepository := repo.NewRainbondClusterConfigRepo(db)
	clusterUsecase := usecase.NewClusterUsecase(db, taskProducer, cloudAccesskeyRepository, createKubernetesTaskRepository, initRainbondTaskRepository, updateKubernetesTaskRepository, taskEventRepository, rainbondClusterConfigRepository)
	clusterHandler := handler.NewClusterHandler(clusterUsecase)
	appStoreUsecase := usecase.NewAppStoreUsecase(appStoreRepo)
	templateVersioner := appstore.NewTemplateVersioner(configConfig)
	templateVersionRepo := repo.NewTemplateVersionRepo(templateVersioner)
	appTemplate := usecase.NewAppTemplate(templateVersionRepo)
	appStoreHandler := handler.NewAppStoreHandler(appStoreUsecase, appTemplate)
	systemHandler := handler.NewSystemHandler(db)
	router := handler.NewRouter(middlewareMiddleware, clusterHandler, appStoreHandler, systemHandler)
	createKubernetesTaskHandler := task.NewCreateKubernetesTaskHandler(clusterUsecase)
	cloudInitTaskHandler := task.NewCloudInitTaskHandler(clusterUsecase)
	updateKubernetesTaskHandler := task.NewCloudUpdateTaskHandler(clusterUsecase)
	engine := newApp(contextContext, router, arg, arg2, arg3, createKubernetesTaskHandler, cloudInitTaskHandler, updateKubernetesTaskHandler)
	return engine, nil
}
