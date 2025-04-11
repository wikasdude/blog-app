package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	DB     *sql.DB
	Config AppConfig
)

type AppConfig struct {
	Environments map[string]EnvironmentConfig `mapstructure:"environments"`
}

type EnvironmentConfig struct {
	Hostname string `mapstructure:"hostname"`
	DB       struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"db"`
}

func LoadConfig() (EnvironmentConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	hostname, err := os.Hostname()
	fmt.Println(hostname)
	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}
	hostname = strings.ToLower(hostname)
	log.Println("Detected Hostname:", hostname)

	for env, cfg := range Config.Environments {
		if strings.ToLower(cfg.Hostname) == hostname {
			log.Printf("Using environment: %s", env)
			return cfg, err
		}
	}

	log.Fatalf("No matching environment found for hostname: %s", hostname)
	return EnvironmentConfig{}, err
}

func ConnectDB() (*sql.DB, error) {
	// envConfig, err := LoadConfig()
	// fmt.Println("Db name is : ", envConfig.DB.Name)

	// dsn := fmt.Sprintf(
	// 	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	envConfig.DB.Host,
	// 	envConfig.DB.Port,
	// 	envConfig.DB.User,
	// 	envConfig.DB.Password,
	// 	envConfig.DB.Name,
	// )

	// db, err := sql.Open("postgres", dsn)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to DB: %v", err)
	// }
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	//Connection Pooling Config
	db.SetMaxOpenConns(25)                  // Max total open connections
	db.SetMaxIdleConns(10)                  // Max idle connections
	db.SetConnMaxIdleTime(5 * time.Minute)  // Max idle time
	db.SetConnMaxLifetime(30 * time.Minute) // Max lifetime of a conne¸ction

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	log.Println("Database connection established")
	DB = db
	return db, err
}
