package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_migrations "github.com/golang-migrate/migrate/v4"

	"booking-event/config"
	"booking-event/internal/common/appcontext"
	authhttphandler "booking-event/internal/modules/auth/transport/http"
	"booking-event/migrations"
)

type Server struct {
	appContext appcontext.AppContext
	router     *gin.Engine
	config     config.Config
	server     *http.Server
}

func NewServer(config config.Config) *Server {
	appContext := appcontext.NewAppContext(config)
	return &Server{config: config, router: gin.New(), appContext: appContext}
}

func (s *Server) RegisterMiddlewares() {
	s.router.Use(gin.Recovery())
}

func (s *Server) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (s *Server) RegisterRoutes() {
	s.router.GET("/health", s.HealthCheck)

	userRoutes := s.router.Group("/api/v1")
	authHttpHandler := authhttphandler.NewAuthHandler(s.appContext.ServiceRegistry().AuthService())
	authHttpHandler.RegisterRoutes(userRoutes)
}

func (s *Server) Run() error {
	if err := migrations.RunMigrations(s.appContext.InfraRegistry().DBUrl()); err != nil && err != _migrations.ErrNoChange {
		fmt.Println("Failed to run migrations:", err)
	}
	s.RegisterMiddlewares()
	s.RegisterRoutes()

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.Server.Port),
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		s.server.Shutdown(ctx)
	}

	if err := s.appContext.InfraRegistry().DB().Close(); err != nil {
		return err
	}

	if err := s.appContext.InfraRegistry().Redis().Close(); err != nil {
		return err
	}

	return nil
}
