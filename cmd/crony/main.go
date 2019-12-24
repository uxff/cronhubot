/**
  运行方式：
 	APPENV=beta SERVICE_PORT=9001 DATASTORE_URL='mysql://yourusername:yourpwd(yourmysqlhost)/yourdbname?charset=utf8mb4&parseTime=True&loc=Local' ./main


*/
package main

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/uxff/cronhubot/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	slog "log"
	"os"
)

var (
	version     string
	versionFlag bool
	logLevel = 0
)


func main() {
	versionString := "Cron Service v" + version
	cobra.OnInitialize(func() {
		if versionFlag {
			fmt.Println(versionString)
			os.Exit(0)
		}
	})

	flag.IntVar(&logLevel, "l", logLevel, "log logLevel, -1:debug, 0:info, 1:warn, 2:error")
	flag.Parse()

	lcf := zap.NewDevelopmentConfig()
	lcf.Level.SetLevel(zapcore.Level(logLevel))
	lcf.Development = false
	lcf.DisableStacktrace = true
	logger, err := lcf.Build(zap.AddCallerSkip(1))
	if err != nil {
		slog.Fatalln("new log err:", err.Error())
	}


	log.SetLogger(logger.Sugar())

	var rootCmd = &cobra.Command{
		Use:   "cron-srv",
		Short: "Cron Service",
		Long:  versionString,
		Run:   Serve,
	}

	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print application version")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
