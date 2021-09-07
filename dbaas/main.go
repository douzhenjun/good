/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: Dou
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: Dou
 * @LastEditTime: 2021-02-07 16:32:07
 */

package main

import (
	"DBaas/config"
	"DBaas/controller"
	"DBaas/datasource"
	"DBaas/grpc/dbaasgrpcservice"
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/robfig/cron/v3"
)

func main() {
	app := newApp()
	//应用App设置
	configuration(app)
	c := config.GetConfig()
	//路由设置
	mvcHandle(app, c)
	addr := ":" + c.Port
	_ = app.Run(
		iris.Addr(addr),
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

//构建App
func newApp() *iris.Application {
	app := iris.New()
	app.Use(irisMiddleware)
	// 设置日志级别  开发阶段为debug
	app.Logger().SetLevel("debug")
	return app
}

func irisMiddleware(ctx iris.Context) {
	utils.LoggerInfo(ctx.Path())
	ctx.Next()
}

/**
 * MVC 架构模式处理
 */
func mvcHandle(app *iris.Application, c *config.AppConfig) {
	//实例化pg数据库引擎
	engine := datasource.NewPgEngine(c)
	influxdbClient, err := datasource.NewInfluxdbEngine(c)
	utils.LoggerError(err)
	corn := cron.New()
	// 初始化k8s
	k8sConfig := models.Sysparameter{ParamKey: "kubernetes_config"}
	_, _ = engine.Get(&k8sConfig)
	kConfig, clientSet, ctx, err := service.InitK8s(k8sConfig.ParamValue)
	utils.LoggerError(err)
	logService := service.NewLogService(engine)
	commonService := service.NewCommonService(engine, kConfig, clientSet, ctx, err, logService, influxdbClient)
	commonService.SetNameSpace()
	clusterService := service.NewClusterService(engine, commonService)
	podService := service.NewPodService(engine, commonService)
	storageService := service.NewStorageService(engine, commonService)
	nodeService := service.NewNodeService(engine, commonService)
	userService := service.NewUserService(engine)
	parameterService := service.NewParameterService(engine)
	imageService := service.NewImageService(engine, commonService)
	homeService := service.NewHomeService(engine, commonService)
	initService := service.NewInitService(engine)
	backupService := service.NewBackupService(engine, commonService)
	externalService := service.NewExternalService(engine, clusterService, commonService)
	quotaService := service.NewQuotaService(engine)
	comboService := service.NewComboService(engine)

	testService := service.NewTestService()

	// test模块
	test := mvc.New(app.Party("/api/dbaas/test"))
	test.Register(testService)
	test.Handle(new(controller.TestController))

	// 公共模块
	common := mvc.New(app.Party("/api/dbaas/common"))
	common.Register(commonService)
	common.Handle(new(controller.CommonController))

	// 存储模块
	storage := mvc.New(app.Party("/api/dbaas/storage"))
	storage.Register(storageService, commonService)
	storage.Handle(new(controller.StorageController))

	// node节点
	node := mvc.New(app.Party("/api/dbaas/host"))
	node.Register(nodeService, commonService)
	node.Handle(new(controller.NodeController))

	// 集群实例模块
	cluster := mvc.New(app.Party("/api/dbaas/cluster"))
	cluster.Register(clusterService, commonService, userService)
	cluster.Handle(new(controller.ClusterController))

	// pod模块
	pod := mvc.New(app.Party("/api/dbaas/pod"))
	pod.Register(podService, commonService)
	pod.Handle(new(controller.PodController))

	// 用户模块
	user := mvc.New(app.Party("/api/dbaas/user"))
	user.Register(userService, commonService, storageService)
	user.Handle(new(controller.UserController))

	// 系统参数模块
	parameter := mvc.New(app.Party("/api/dbaas/parameter"))
	parameter.Register(parameterService, commonService)
	parameter.Handle(new(controller.ParameterController))

	// 系统日志模块
	log := mvc.New(app.Party("/api/dbaas/log"))
	log.Register(logService, commonService)
	log.Handle(new(controller.LogController))

	// 镜像模块
	image := mvc.New(app.Party("api/dbaas/image"))
	image.Register(imageService, commonService)
	image.Handle(new(controller.ImageController))

	// 首页模块
	home := mvc.New(app.Party("api/dbaas/home"))
	home.Register(commonService, nodeService, userService, clusterService, homeService)
	home.Handle(new(controller.HomeController))

	// 初始化模块
	init := mvc.New(app.Party("api/dbaas/init"))
	init.Register(initService, commonService, parameterService, nodeService, storageService, imageService)
	init.Handle(new(controller.InitController))

	// 外部调用api
	external := mvc.New(app.Party("api/dbaas/external"))
	external.Register(externalService, commonService)
	external.Handle(new(controller.ExternalController))

	// 备份模块
	backup := mvc.New(app.Party("api/dbaas/backup"))
	backup.Register(backupService, commonService)
	backup.Handle(new(controller.BackupController))

	// 接口配额
	quota := mvc.New(app.Party("api/dbaas/quota"))
	quota.Register(quotaService)
	quota.Handle(new(controller.QuotaController))

	// 统计模块
	service.InitStatistics(engine)
	statistics := mvc.New(app.Party("api/dbaas/statistics"))
	statistics.Handle(new(controller.StatisticsController))

	// 实例套餐
	combo := mvc.New(app.Party("api/dbaas/combo"))
	combo.Register(comboService)
	combo.Handle(new(controller.ComboController))

	go cronTask(corn, commonService)
	go dbaasgrpcservice.RungGRPCServer(podService, c)
}

/**
 * 项目设置
 */
func configuration(app *iris.Application) {
	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))
	//错误配置
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		_, _ = context.JSON(iris.Map{
			"error_msg_en": " not found ",
			"error_msg_zh": iris.StatusNotFound,
			"data":         iris.Map{},
		})
	})
	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		_, _ = context.JSON(iris.Map{
			"data":         iris.Map{},
			"error_msg_en": " internal error ",
			"error_msg_zh": iris.StatusInternalServerError,
		})
	})
}

//添加定时任务
func cronTask(cc *cron.Cron, cs service.CommonService) {
	every10s, every20s, every600s := "@every 10s", "@every 20s", "@every 600s"

	_, err := cc.AddFunc(every10s, cs.AsyncMetricsPods)
	utils.LoggerError(err)

	_, err = cc.AddFunc(every10s, cs.AsyncCommonInfo)
	utils.LoggerError(err)

	_, err = cc.AddFunc(every20s, cs.AsyncNodeInfo)
	utils.LoggerError(err)

	_, err = cc.AddFunc(every20s, cs.AsyncOperatorLog)
	utils.LoggerError(err)

	_, err = cc.AddFunc(every20s, cs.AsyncBackupJob)
	utils.LoggerError(err)

	_, err = cc.AddFunc(every600s, cs.AsyncImageStatus)
	utils.LoggerError(err)

	_, err = cc.AddFunc(every600s, cs.AsyncPVStatus)
	utils.LoggerError(err)

	every12h := "@every 12h"
	_, err = cc.AddFunc(every12h, cs.ClearUselessEvents)
	utils.LoggerError(err)

	cc.Start()
}
