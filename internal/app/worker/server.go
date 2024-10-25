package worker

import (
	"context"

	"booking-event/config"
	"booking-event/internal/common/appcontext"
	"booking-event/internal/infra/asynq"
	"booking-event/internal/infra/redis"
	"booking-event/internal/modules/booking/transport/asyntask"
)

type Server struct {
	appContext    appcontext.AppContext
	config        config.Config
	asynqServer   *asynq.AsynqServer
	asynqHandlers *asyntask.EmailTaskHandler
}

func NewServer(config config.Config) *Server {
	appContext := appcontext.NewAppContext(config)
	redisConfig := redis.Config{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		Prefix:   config.Redis.Prefix,
	}
	asynqServer := asynq.NewAsynqServer(asynq.Config{
		Addr:        redisConfig.Addr(),
		Concurrency: config.Asynq.Concurrency,
		Queues:      config.Asynq.Queues,
	})
	return &Server{config: config, appContext: appContext, asynqServer: asynqServer}
}

func (s *Server) RegisterHandlers() {
	handlers := asyntask.NewEmailTaskHandler(s.appContext.ServiceRegistry().EmailService())
	handlers.Register(s.asynqServer.ServeMux())
	s.asynqHandlers = handlers
}

func (s *Server) Run() error {
	s.RegisterHandlers()
	return s.asynqServer.Start(context.Background())
}
