package main

import (
	"flag"
	"fmt"
	"github.com/uxff/cronhubot/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	slog "log"
	"os"
	"github.com/spf13/cobra"
)

var (
	version     string
	versionFlag bool
)
var (
	level    = flag.Int("l", 0, "log level, -1:debug, 0:info, 1:warn, 2:error")
	confFile = flag.String("conf", "conf/files/base.json", "config file path")
)
func main() {
	versionString := "Cron Service v" + version
	cobra.OnInitialize(func() {
		if versionFlag {
			fmt.Println(versionString)
			os.Exit(0)
		}
	})

	flag.Parse()

	lcf := zap.NewDevelopmentConfig()
	lcf.Level.SetLevel(zapcore.Level(*level))
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
