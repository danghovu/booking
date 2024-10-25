package appcontext

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"

	"booking-event/config"
	"booking-event/internal/infra/asynq"
	"booking-event/internal/infra/emailsender"
	"booking-event/internal/infra/paymentgateway"
	postgresql "booking-event/internal/infra/posgresql"
	"booking-event/internal/infra/redis"
)

type InfraRegistry interface {
	DB() *sqlx.DB
	DBUrl() string
	Redis() redis.Redis

	EmailService() emailsender.EmailService
	PaymentService() paymentgateway.PaymentGateway
	AsyncTaskEnqueueClient() asynq.AsyncTaskEnqueueClient
}

type infraRegistry struct {
	db                     *sqlx.DB
	redis                  redis.Redis
	kafkaSyncProducer      sarama.SyncProducer
	emailService           emailsender.EmailService
	paymentGateway         paymentgateway.PaymentGateway
	asyncTaskEnqueueClient asynq.AsyncTaskEnqueueClient
	dbUrl                  string
}

func NewInfraRegistry(config config.Config) InfraRegistry {
	dbConfig := postgresql.Config{
		Host:     config.Postgres.Host,
		Port:     config.Postgres.Port,
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		DBName:   config.Postgres.DbName,
		SSLMode:  config.Postgres.SSLMode,
		IdleConn: config.Postgres.IdleConn,
		MaxOpen:  config.Postgres.OpenConn,
	}

	db := postgresql.NewClient(dbConfig)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisConfig := redis.Config{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		Prefix:   config.Redis.Prefix,
	}

	redis := redis.NewClient(redisConfig)

	if _, err := redis.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	emailService := emailsender.NewNoopEmailService()

	paymentGateway := paymentgateway.NewNoopPaymentGateway()

	return &infraRegistry{
		db:                     db,
		redis:                  redis,
		emailService:           emailService,
		paymentGateway:         paymentGateway,
		asyncTaskEnqueueClient: asynq.NewEnqueueClient(asynq.Config{Addr: redisConfig.Addr()}),
		dbUrl:                  dbConfig.URL(),
	}
}

func (r *infraRegistry) DB() *sqlx.DB {
	return r.db
}

func (r *infraRegistry) Redis() redis.Redis {
	return r.redis
}
func (r *infraRegistry) EmailService() emailsender.EmailService {
	return r.emailService
}

func (r *infraRegistry) PaymentService() paymentgateway.PaymentGateway {
	return r.paymentGateway
}

func (r *infraRegistry) DBUrl() string {
	return r.dbUrl
}

func (r *infraRegistry) AsyncTaskEnqueueClient() asynq.AsyncTaskEnqueueClient {
	return r.asyncTaskEnqueueClient
}
