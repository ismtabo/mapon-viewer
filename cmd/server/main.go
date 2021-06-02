package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ismtabo/mapon-viewer/pkg/cfg"
	"github.com/ismtabo/mapon-viewer/pkg/controller"
	"github.com/ismtabo/mapon-viewer/pkg/repository"
	"github.com/ismtabo/mapon-viewer/pkg/routes"
	"github.com/ismtabo/mapon-viewer/pkg/service"
	"github.com/ismtabo/mapon-viewer/pkg/template"
	"github.com/kataras/go-sessions/v3"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	log := getLogger()
	var config Config
	if err := cfg.Load("config.yml", &config); err != nil {
		log.Fatal().Msgf("Error loading configuration. %s", err)
	}
	log.Debug().Msgf("Configuration loaded: %+v", config)
	if err := configLogger(&config); err != nil {
		log.Fatal().Msgf("Error configuring the logger. %s", err)
	}
	ctx := context.Background()
	connection := getMongoConnection(ctx, &config, log)
	database := connection.Database(config.Mongo.Database)
	collection := database.Collection(config.Mongo.Collections.Users)
	ssns := getSessionsFactory(&config)
	secSvc := service.NewSecurityService(ssns)
	userRepo := repository.NewMongoUserRepository(collection)
	userSvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSvc, secSvc)
	tmplMngr := template.NewTemplateManager()
	pagesCtrl := controller.NewPagesController(ssns, tmplMngr)
	maponRepo := repository.NewMaponRespository(&config.MaponConfig)
	maponCtrl := controller.NewMaponController(maponRepo)
	routes := routes.NewRoutes(userCtrl, pagesCtrl, maponCtrl, secSvc, log)
	routes.AddRoutes()
	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	if err := http.ListenAndServe(addr, routes); err != nil {
		log.Fatal().Msg("Error starting server")
	}
}

func getLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}

func configLogger(config *Config) error {
	lvl, err := zerolog.ParseLevel(strings.ToLower(config.Log.Level))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)
	return nil
}

func getMongoConnection(ctx context.Context, config *Config, log *zerolog.Logger) *mongo.Client {
	clientOpts := options.Client()
	clientOpts.SetHosts([]string{config.Mongo.Host})
	if config.Mongo.User != "" {
		clientOpts.SetAuth(options.Credential{
			Username: config.Mongo.User,
			Password: string(config.Mongo.Password),
		})
	}
	conn, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting mongo server.")
	}
	if err := conn.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal().Err(err).Msg("Error connecting to mongo.")
	}
	return conn
}

func getSessionsFactory(config *Config) *sessions.Sessions {
	return sessions.New(sessions.Config{
		Cookie:  config.Session.Cookie,
		Expires: time.Duration(config.Session.Expires) * time.Second,
	})
}
