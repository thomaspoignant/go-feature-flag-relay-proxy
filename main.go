package main

import (
	"context"
	"fmt"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/log"
	"go.uber.org/zap"
	"time"

	"github.com/thomaspoignant/go-feature-flag-relay-proxy/api"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/docs"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
)

// version, releaseDate are override by the makefile during the build.
var version = "localdev"

const banner = `█▀▀ █▀█   █▀▀ █▀▀ ▄▀█ ▀█▀ █ █ █▀█ █▀▀   █▀▀ █   ▄▀█ █▀▀
█▄█ █▄█   █▀  ██▄ █▀█  █  █▄█ █▀▄ ██▄   █▀  █▄▄ █▀█ █▄█

     █▀█ █▀▀ █   ▄▀█ █▄█   █▀█ █▀█ █▀█ ▀▄▀ █▄█
     █▀▄ ██▄ █▄▄ █▀█  █    █▀▀ █▀▄ █▄█ █ █  █ 

GO Feature Flag Relay Proxy
_____________________________________________`

// @title go-feature-flag relay proxy
// @description Swagger for the go-feature-flag relay proxy.
// @description
// @description go-feature-flag relay proxy is using thomaspoignant/go-feature-flag to handle your feature flagging.
// @description It is a proxy to your flags, you can get the values of your flags using APIs.
// @contact.name GO feature flag relay proxy
// @contact.url https://github.com/thomaspoignant/go-feature-flag-relay-proxy
// @license.name MIT
// @license.url https://github.com/thomaspoignant/go-feature-flag-relay-proxy/blob/main/LICENSE
// @BasePath /
func main() {
	// TODO: options to implement
	// - hideBanner
	// - enableSwagger - default is false
	// - debug mode for echo
	// - fail on startup
	// - HTTP port

	zapLog := log.InitLogger()
	defer func() { _ = zapLog.Sync() }()

	// TODO: should we print this banner ?
	fmt.Println(banner)

	// Init swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", "3000") // TODO: change host by the config

	// Init services
	goff, err := ffclient.New(
		ffclient.Config{
			StartWithRetrieverError: true,
			PollingInterval:         10 * time.Second,
			Logger:                  zap.NewStdLog(zapLog),
			Context:                 context.Background(),
			Retriever: &ffclient.GithubRetriever{
				RepositorySlug: "thomaspoignant/go-feature-flag",
				Branch:         "main",
				FilePath:       "testdata/flag-config.yaml",
				Timeout:        3 * time.Second,
			},
		})

	if err != nil {
		panic(err)
	}

	monitoringService := service.NewMonitoring(goff)

	// Init API server
	apiServer := api.New(api.ServerConfig{}, monitoringService, goff, zapLog)
	apiServer.Start()
	defer func() { _ = apiServer.Stop }()
}
