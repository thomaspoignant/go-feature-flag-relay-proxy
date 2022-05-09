package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/controller"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
	"log"
)

// New is used to create a new instance of the API server
func New(config ServerConfig, monitoringService service.Monitoring, goFF *ffclient.GoFeatureFlag) Server {
	s := Server{
		serverConfig:      config,
		monitoringService: monitoringService,
		goFF:              goFF,
	}
	s.init()
	return s
}

// ServerConfig is the configuration struct for the API server
type ServerConfig struct {
	Port int
}

// Server is the struct that represent the API server
type Server struct {
	serverConfig      ServerConfig
	echoInstance      *echo.Echo
	monitoringService service.Monitoring
	goFF              *ffclient.GoFeatureFlag
}

// init initialize the configuration of our API server (using echo)
func (s *Server) init() {
	s.echoInstance = echo.New()
	s.echoInstance.HideBanner = true

	// Middlewares
	s.echoInstance.Use(middleware.Logger())
	s.echoInstance.Use(middleware.Recover())

	// Init controllers
	cHealth := controller.NewHealth(s.monitoringService)
	cInfo := controller.NewInfo(s.monitoringService)
	cAllFlags := controller.NewAllFlags(s.goFF)
	cFlagEval := controller.NewFlagEval(s.goFF)

	// health Routes
	s.echoInstance.GET("/health", cHealth.Handler)
	s.echoInstance.GET("/info", cInfo.Handler)

	// GO feature flags routes
	v1 := s.echoInstance.Group("/v1")
	v1.POST("/allflags", cAllFlags.Handler)
	v1.POST("/feature/:flagKey/eval", cFlagEval.Handler)
}

// Start launch the API server
func (s *Server) Start() {
	if s.serverConfig.Port == 0 {
		s.serverConfig.Port = 3000
	}
	s.echoInstance.Logger.Fatal(s.echoInstance.Start(fmt.Sprintf(":%d", s.serverConfig.Port)))
}

// Stop shutdown the API server
func (s *Server) Stop() {
	err := s.echoInstance.Close()
	if err != nil {
		log.Fatalf("impossible to stop go-feature-flag relay proxy %e", err)
	}
}
