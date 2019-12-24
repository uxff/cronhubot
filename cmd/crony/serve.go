package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/uxff/cronhubot/pkg/checker"
	"github.com/uxff/cronhubot/pkg/config"
	"github.com/uxff/cronhubot/pkg/datastore"
	"github.com/uxff/cronhubot/pkg/handlers"
	"github.com/uxff/cronhubot/pkg/repos"
	"github.com/uxff/cronhubot/pkg/scheduler"
	"github.com/uxff/cronhubot/pkg/log"
)

func Serve(cmd *cobra.Command, args []string) {
	log.Infof("cron Service starting...")

	env, err := config.LoadEnv()
	failOnError(err, "Failed to load config!")
	//if strings.ToLower(env.APPENV) == "pro" {
	//	log.SetLevel(log.LevelINFO)
	//} else {
	//	log.SetLevel(log.LevelDEBUG)
	//}
	log.Debugf("环境变量：%+v", env)

	ds, err := datastore.New(env.DatastoreURL)
	failOnError(err, "Failed to init dababase connection!")
	defer ds.Close()

	checkers := map[string]checker.Checker{
		"api": checker.NewApi(),
		// "postgres": checker.NewPostgres(env.DatastoreURL),
		"mysql": checker.NewMysql(env.DatastoreURL),
	}
	healthzHandler := handlers.NewHealthzHandler(checkers)

	eventRepo := repos.NewEvent(ds)
	authRepo := repos.NewAuth(ds)

	sc := scheduler.New(eventRepo)
	go sc.ScheduleAll()

	eventsHandler := handlers.NewEventsHandler(authRepo, eventRepo, sc)

	router := gin.New()

	router.GET("/health", healthzHandler.HealthzIndex)
	group := router.Group("/", eventsHandler.BasicMiddleware)
	// group.GET("/events", eventsHandler.EventsIndex)
	group.POST("/events", eventsHandler.EventsCreate)
	// group.GET("/events/:id", eventsHandler.EventsShow)
	// group.PUT("/events/:id", eventsHandler.EventsUpdate)
	// group.DELETE("/events/:id", eventsHandler.EventsDelete)

	addr := fmt.Sprintf(":%d", env.Port)
	log.Fatalf("%v", router.Run(addr))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, err:%s", msg, err)
	}
}
