// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"goodrain.com/cloud-adaptor/internal/biz"
	"goodrain.com/cloud-adaptor/internal/handler"
	"goodrain.com/cloud-adaptor/internal/middleware"
	"goodrain.com/cloud-adaptor/internal/nsqc/producer"
	"goodrain.com/cloud-adaptor/internal/repo"
	"goodrain.com/cloud-adaptor/internal/repo/appstore"
	"goodrain.com/cloud-adaptor/internal/repo/dao"
	"goodrain.com/cloud-adaptor/internal/task"
	"goodrain.com/cloud-adaptor/internal/types"
	"gorm.io/gorm"
)

// Injectors from wire.go:

// initApp init the application.
func initApp(contextContext context.Context, db *gorm.DB, arg chan types.KubernetesConfigMessage, arg2 chan types.InitRainbondConfigMessage, arg3 chan types.UpdateKubernetesConfigMessage) (*gin.Engine, error) {
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
	clusterUsecase := biz.NewClusterUsecase(db, taskProducer, cloudAccesskeyRepository, createKubernetesTaskRepository, initRainbondTaskRepository, updateKubernetesTaskRepository, taskEventRepository, rainbondClusterConfigRepository)
	clusterHandler := handler.NewClusterHandler(clusterUsecase)
	appStoreUsecase := biz.NewAppStoreUsecase(appStoreRepo)
	appStoreHandler := handler.NewAppStoreHandler(appStoreUsecase)
	systemHandler := handler.NewSystemHandler(db)
	router := handler.NewRouter(middlewareMiddleware, clusterHandler, appStoreHandler, systemHandler)
	createKubernetesTaskHandler := task.NewCreateKubernetesTaskHandler(clusterUsecase)
	cloudInitTaskHandler := task.NewCloudInitTaskHandler(clusterUsecase)
	updateKubernetesTaskHandler := task.NewCloudUpdateTaskHandler(clusterUsecase)
	engine := newApp(contextContext, router, arg, arg2, arg3, createKubernetesTaskHandler, cloudInitTaskHandler, updateKubernetesTaskHandler)
	return engine, nil
}
