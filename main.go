package main

import (
	"context"
	"fmt"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"log"
	"os"
	"time"

	"github.com/thomaspoignant/go-feature-flag-relay-proxy/api"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
)

const banner = `█▀▀ █▀█   █▀▀ █▀▀ ▄▀█ ▀█▀ █ █ █▀█ █▀▀   █▀▀ █   ▄▀█ █▀▀
█▄█ █▄█   █▀  ██▄ █▀█  █  █▄█ █▀▄ ██▄   █▀  █▄▄ █▀█ █▄█

█▀█ █▀▀ █   ▄▀█ █▄█   █▀█ █▀█ █▀█ ▀▄▀ █▄█
█▀▄ ██▄ █▄▄ █▀█  █    █▀▀ █▀▄ █▄█ █ █  █ 

GO Feature Flag Relay Proxy
_____________________________________________`

func main() {
	// TODO: should we print this banner ?
	fmt.Println(banner)

	// Init services
	goff, err := ffclient.New(
		ffclient.Config{
			PollingInterval: 10 * time.Second,
			Logger:          log.New(os.Stdout, "", 0),
			Context:         context.Background(),
			Retriever: &ffclient.GithubRetriever{
				RepositorySlug: "thomaspoignant/go-feature-flag",
				Branch:         "main",
				FilePath:       "examples/github/flags.yaml",
				Timeout:        3 * time.Second,
			},
		})

	if err != nil {
		panic(err)
	}

	monitoringService := service.NewMonitoring(goff)

	// Init API server
	apiServer := api.New(api.ServerConfig{}, monitoringService, goff)
	apiServer.Start()
	defer func() { _ = apiServer.Stop }()
}
