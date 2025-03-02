package main

import (
	"Webook/cronjob/domain"
	"Webook/cronjob/ioc"
	"Webook/cronjob/job"
	"Webook/cronjob/repository"
	"Webook/cronjob/repository/dao"
	"Webook/cronjob/service"
	"Webook/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"net"
)

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/dev.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViperV2Watch()
	app := Init()
	listen, err := net.Listen("tcp", ":"+"8080")
	err = app.server.Server.Serve(listen)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	l := ioc.InitLogger()
	rankingClient := ioc.InitRankingRpcClient()
	//启动的时候自动添加执行器，调度器，任务
	rankingJob := domain.CronJob{
		Name:       "ranking_job",
		Executor:   "local",
		Cfg:        "",
		Expression: "0 12 * * *",
		CancelFunc: func() {},
	}

	//这里可以做成自动注入
	jobRepository := repository.NewPreemptCronJobRepository(dao.NewGORMJobDAO(ioc.InitDB(l)))
	err = jobRepository.AddJob(ctx, rankingJob)
	if err != nil {
		panic("热榜定时初始化失败")
	}
	cronService := service.NewCronJobService(jobRepository, l)
	scheduler := job.NewScheduler(cronService, l)

	executor := job.NewLocalFuncExecutor(rankingClient)
	scheduler.RegisterExecutor(executor)
	executor.RegisterFunc("ranking_job", executor.Ranking)
	if err := scheduler.Schedule(ctx); err != nil {
		l.Error("调度器运行失败", logger.Error(err))
	}

}
