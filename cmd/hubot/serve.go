package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/uxff/cronhubot/pkg/checker"
	"github.com/uxff/cronhubot/pkg/config"
	"github.com/uxff/cronhubot/pkg/datastore"
	"github.com/uxff/cronhubot/pkg/handlers"
	"github.com/uxff/cronhubot/pkg/log"
	"github.com/uxff/cronhubot/pkg/repos"
	"github.com/uxff/cronhubot/pkg/scheduler"
)

func Serve(env *config.Config) error {
	log.Infof("cron Service starting...")

	ds, err := datastore.New(env.DatastoreURL)
	failOnError(err, "Failed to init dababase connection!")
	defer ds.Close()

	checkers := map[string]checker.Checker{
		"api": checker.NewApi(),
		// "postgres": checker.NewPostgres(env.DatastoreURL),
		"mysql": checker.NewMysql(env.DatastoreURL),
	}
	healthzHandler := handlers.NewHealthzHandler(checkers)

	eventRepo := repos.NewCronJob(ds)
	authRepo := repos.NewAuth(ds)

	sc := scheduler.New(eventRepo)
	go sc.ScheduleAll()

	eventsHandler := handlers.NewJobsHandler(authRepo, eventRepo, sc)

	router := gin.New()

	router.GET("/health", healthzHandler.HealthzIndex)
	group := router.Group("/", eventsHandler.BasicMiddleware)
	group.GET("/events", eventsHandler.JobsIndex)
	group.POST("/events", eventsHandler.JobsCreate)
	group.GET("/events/:id", eventsHandler.JobsDetail)
	// group.PUT("/events/:id", eventsHandler.JobsUpdate)
	// group.DELETE("/events/:id", eventsHandler.JobsDelete)

	addr := fmt.Sprintf(":%d", env.Port)
	log.Debugf("service will start at %s", addr)

	return router.Run(addr)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, err:%s", msg, err)
	}
}