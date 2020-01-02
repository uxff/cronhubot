/**
  分布式应用入口 (Distributed)
  运行方式：
 	APPENV=beta SERVICE_PORT=9001 DATASTORE_URL='mysql://yourusername:yourpwd@tcp(yourmysqlhost)/yourdbname?charset=utf8mb4&parseTime=True&loc=Local' ./main


*/
package main

import (
	"flag"
	"github.com/uxff/cronhubot/pkg/config"
	"github.com/uxff/cronhubot/pkg/log"
	"github.com/uxff/electio/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	slog "log"
	"os"
	"strings"
	"sync"
)

var (
	version         = "0.1"
	showVersion 	bool
	logLevel     	= 0
	clusterName     = "c1"
	clusterMembers  = "127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003"
	clusterNodeAddr = "127.0.0.1:9001"
)


func main() {

	flag.IntVar(&logLevel, "l", logLevel, "log logLevel, -1:debug, 0:info, 1:warn, 2:error")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&clusterName, "c", clusterName, "cluster name on it")
	flag.StringVar(&clusterMembers, "m", clusterMembers, "cluster members")
	flag.StringVar(&clusterNodeAddr, "a", clusterNodeAddr, "cluster node addr that you will listen as serve as a cluster node")
	flag.Parse()

	if showVersion {
		flag.Usage()
		os.Exit(0)
	}

	lcf := zap.NewDevelopmentConfig()
	lcf.Level.SetLevel(zapcore.Level(logLevel))
	lcf.Development = false
	lcf.DisableStacktrace = true
	logger, err := lcf.Build(zap.AddCallerSkip(1))
	if err != nil {
		slog.Fatalln("new log err:", err.Error())
	}

	log.SetLogger(logger.Sugar())

	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load env:%v", err)
	}


	workerNode := worker.NewWorker(clusterNodeAddr, clusterName)
	workerNode.AddMates(strings.Split(clusterMembers, ","))

	errCh := make(chan error, 1)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// 启动cron服务
	go func() {
		errCh<-Serve(env)
		wg.Done()
	}()

	// 启动分布式集群节点
	go func() {
		errCh <- workerNode.Start()
		wg.Done()
	}()

	err = <-errCh
	if err != nil {
		log.Fatalf("main error:%v", err)
	}
}
