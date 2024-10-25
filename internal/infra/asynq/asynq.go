package asynq

import (
	"context"

	"github.com/hibiken/asynq"
)

type Config struct {
	Addr        string
	Concurrency int
	Queues      map[string]int
}

type AsynqServer struct {
	server   *asynq.Server
	handlers map[string]asynq.Handler
	mux      *asynq.ServeMux
}

func NewAsynqServer(config Config) *AsynqServer {
	queues := make(map[string]int)
	for queue, concurrency := range config.Queues {
		queues[queue] = concurrency
	}
	if queues["default"] == 0 {
		queues["default"] = 3
	}
	mux := asynq.NewServeMux()
	return &AsynqServer{
		server: asynq.NewServer(
			asynq.RedisClientOpt{Addr: config.Addr},
			asynq.Config{
				Concurrency: config.Concurrency,
				Queues:      config.Queues,
			},
		),
		handlers: make(map[string]asynq.Handler),
		mux:      mux,
	}
}

func (s *AsynqServer) ServeMux() *asynq.ServeMux {
	return s.mux
}

func (s *AsynqServer) RegisterHandler(taskName string, handler asynq.Handler) {
	s.handlers[taskName] = handler
}

func (s *AsynqServer) Start(ctx context.Context) error {
	for taskName, handler := range s.handlers {
		s.mux.HandleFunc(taskName, handler.ProcessTask)
	}
	return s.server.Run(s.mux)
}
