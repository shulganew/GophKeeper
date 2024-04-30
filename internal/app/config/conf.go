package config

import (
	"flag"
	"os"
	"time"

	"github.com/shulganew/GophKeeper/internal/app/validators"
	"go.uber.org/zap"
)

const AuthPrefix = "Bearer "
const TokenExp = time.Hour * 3600
const DataBaseType = "postgres"

type Config struct {
	// flag -a, Server address
	Address string

	// dsn connection string
	DSN string

	PassJWT string

	MasterKey string // Master password for GophKeeper storage.
}

func InitConfig() Config {
	config := Config{}
	// read command line argue
	serverAddress := flag.String("a", "localhost:8443", "Service GKeeper address")
	dsnf := flag.String("d", "", "Data Source Name for DataBase connection")
	authJWT := flag.String("p", "JWTsecret", "JWT private key")
	master := flag.String("m", "MasterKey", "Master password for GophKeeper storage")
	flag.Parse()

	// Check and parse URL.
	startaddr, startport := validators.CheckURL(*serverAddress)

	// Server address.
	config.Address = startaddr + ":" + startport

	// read OS ENVs.
	addr, exist := os.LookupEnv(("RUN_ADDRESS"))

	// JWT password for users auth.
	config.PassJWT = *authJWT

	// Master pass.
	config.MasterKey = *master

	// if env var does not exist  - set def value
	if exist {
		config.Address = addr
		zap.S().Infoln("Set result address from evn RUN_ADDRESS: ", config.Address)
	} else {
		zap.S().Infoln("Env var RUN_ADDRESS not found, use default", config.Address)
	}

	dsn, exist := os.LookupEnv(("DATABASE_URI"))

	// Init shotrage DB from env.
	if exist {
		zap.S().Infoln("Use DataBase DSN from evn DATABASE_URI, use: ", dsn)
		config.DSN = dsn
	} else if *dsnf != "" {
		dsn = *dsnf
		zap.S().Infoln("Use DataBase from -d flag, use: ", dsn)
		config.DSN = dsn
	} else {
		zap.S().Errorf("Can't make config for DB, set -d flag or DATABASE_URI env for DSN!")
		os.Exit(65)
	}

	zap.S().Infoln("Configuration complite")
	return config
}
