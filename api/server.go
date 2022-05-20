package api

import (
	"fmt"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/controller"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
	"go.uber.org/zap"
	"time"
)

// New is used to create a new instance of the API server
func New(config ServerConfig,
	monitoringService service.Monitoring,
	goFF *ffclient.GoFeatureFlag,
	zapLog *zap.Logger) Server {
	s := Server{
		serverConfig:      config,
		monitoringService: monitoringService,
		goFF:              goFF,
		zapLog:            zapLog,
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
	zapLog            *zap.Logger
}

// init initialize the configuration of our API server (using echo)
func (s *Server) init() {
	s.echoInstance = echo.New()
	s.echoInstance.HideBanner = true
	s.echoInstance.HidePort = true
	s.echoInstance.Debug = true

	// Middlewares
	s.echoInstance.Use(echozap.ZapLogger(s.zapLog))
	s.echoInstance.Use(middleware.Recover())

	// TODO: timeout configuration from config
	s.echoInstance.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 5 * time.Second,
	}))

	// Init controllers
	cHealth := controller.NewHealth(s.monitoringService)
	cInfo := controller.NewInfo(s.monitoringService)
	cAllFlags := controller.NewAllFlags(s.goFF)
	cFlagEval := controller.NewFlagEval(s.goFF)

	// health Routes
	s.echoInstance.GET("/health", cHealth.Handler)
	s.echoInstance.GET("/info", cInfo.Handler)
	s.echoInstance.GET("/swagger/*", echoSwagger.WrapHandler)

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
	address := fmt.Sprintf("0.0.0.0:%d", s.serverConfig.Port)

	s.zapLog.Info("Starting go-feature-flag relay proxy ...", zap.String("address", address))

	err := s.echoInstance.Start(address)
	if err != nil {
		s.zapLog.Fatal("impossible to start the proxy", zap.Error(err))
	}
}

// Stop shutdown the API server
func (s *Server) Stop() {
	err := s.echoInstance.Close()
	if err != nil {
		s.zapLog.Fatal("impossible to stop go-feature-flag relay proxy", zap.Error(err))
	}
}
