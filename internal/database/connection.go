package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func (c *Config) AdminConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.SSLMode,
	)
}

func Connect(cfg *Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		if strings.Contains(err.Error(), "no existe la base de datos") || strings.Contains(err.Error(), "does not exist") {
			if err := createDatabase(cfg); err != nil {
				return nil, fmt.Errorf("could not create database: %w", err)
			}
			return Connect(cfg)
		}
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &DB{db}, nil
}

func createDatabase(cfg *Config) error {
	db, err := sql.Open("postgres", cfg.AdminConnectionString())
	if err != nil {
		return fmt.Errorf("error connecting to admin db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return fmt.Errorf("error creating database: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
