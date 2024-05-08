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

	PathJWT string

	MasterKey string // Master password for GophKeeper storage.

	Backet string // MINIO backet
}

func InitConfig() Config {
	config := Config{}
	// read command line argue
	serverAddress := flag.String("a", "localhost:8443", "Service GKeeper address")
	dsnf := flag.String("d", "", "Data Source Name for DataBase connection")
	jwtPath := flag.String("p", "cert/jwtpkey.pem", "path to JWT private key file, ex cert/jwtpkey.pem")
	master := flag.String("m", "1NewMasterKey", "Master password for GophKeeper storage")
	mb := flag.String("b", "gohpkeeper", "MINIO backet for files")

	flag.Parse()

	// Check and parse URL.
	startaddr, startport := validators.CheckURL(*serverAddress)

	// Server address.
	config.Address = startaddr + ":" + startport

	// read OS ENVs.
	addr, exist := os.LookupEnv(("RUN_ADDRESS"))

	// JWT password for users auth.
	config.PathJWT = *jwtPath

	// Master pass.
	config.MasterKey = *master

	// MINIO.
	config.Backet = *mb

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
