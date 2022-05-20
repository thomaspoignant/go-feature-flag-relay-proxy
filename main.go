package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/api"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/config"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/docs"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/log"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
	"go.uber.org/zap"
	"net/http"
	"time"
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
	// https://docs.aws.amazon.com/sdkref/latest/guide/creds-config-files.html
	// doc specifier que tous les temps sont en millisecond

	// Init logger
	zapLog := log.InitLogger()
	defer func() { _ = zapLog.Sync() }()

	// Loading the configuration in viper
	proxyConf, err := parseConfig()
	if err != nil {
		zapLog.Fatal("error while reading configuration", zap.Error(err))
	}

	if err := proxyConf.IsValid(); err != nil {
		zapLog.Fatal("configuration error", zap.Error(err))
	}

	if !proxyConf.HideBanner {
		fmt.Println(banner)
	}

	// Init swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", proxyConf.Host, proxyConf.ListenPort)

	// Init services
	goff, err := service.NewGoFeatureFlagClient(proxyConf, zapLog)

	if err != nil {
		panic(err)
	}

	monitoringService := service.NewMonitoring(goff)

	// Init API server
	apiServer := api.New(api.ServerConfig{
		Port: proxyConf.ListenPort,
	}, monitoringService, goff, zapLog)
	apiServer.Start()
	defer func() { _ = apiServer.Stop }()
}

// parseConfig is reading the configuration file
func parseConfig() (*config.Config, error) {
	// default values
	viper.SetDefault("listen", "3000")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("fileFormat", "yaml")

	viper.SetDefault("retriever.timeout", int64(10*time.Second/time.Millisecond))
	viper.SetDefault("retriever.method", http.MethodGet)
	viper.SetDefault("retriever.body", "")

	// TODO add more location + add flags from command line
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./testdata/config/")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	proxyConf := &config.Config{}
	err = viper.Unmarshal(proxyConf)
	if err != nil {
		return nil, err
	}
	return proxyConf, nil
}
