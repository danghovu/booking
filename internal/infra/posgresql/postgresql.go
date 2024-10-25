package postgresql

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	IdleConn int
	MaxOpen  int
}

func (cfg Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func (cfg Config) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}

func NewClient(cfg Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db.SetMaxIdleConns(cfg.IdleConn)
	db.SetMaxOpenConns(cfg.MaxOpen)
	return db
}
