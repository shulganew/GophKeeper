package config

import (
	"fmt"
	"os"
	"time"

	"github.com/shulganew/GophKeeper/internal/app/validators"
)

const AuthPrefix = "Bearer "
const TokenExp = time.Hour * 3600
const DataBaseType = "postgres"

type Config struct {
	Address string // Server address and post, ex localhost:8080.

	DSN string // DB connection.

	PathJWT string

	MasterKey string // Master password for GophKeeper storage.

	ZapPath  string // Logging path.
	ZapLevel string // Logging level.

	Backetmi string // MINIO backet
	IDmi     string
	Secretmi string

	TLSPK  string // Private key
	TLSPub string // Self-signed public servtificate

	RunLocal bool // Run for testing on local host, not ptint logs to file.
}

func InitConfig() Config {
	config := Config{}

	// Read OS ENVs.
	// Master password.
	mp, exist := os.LookupEnv(("GK_ZAP_MASTER"))
	if !exist {
		fmt.Println("ENV GK_ZAP_MASTER not set")
	}
	config.MasterKey = mp

	// Logging level.
	level, exist := os.LookupEnv(("GK_ZAP_LEVEL"))
	if !exist {
		fmt.Println("ENV GK_ZAP_LEVEL not set")
	}
	config.ZapLevel = level

	// Logging path.
	lp, exist := os.LookupEnv(("GK_ZAP_PATH"))
	if !exist {
		fmt.Println("ENV GK_ZAP_PATH not set")
	}
	config.ZapPath = lp

	// Check and parse server URL.
	addr, exist := os.LookupEnv(("GK_ADDRESS"))
	if !exist {
		fmt.Println("ENV GK_ADDRESS not set")
	}
	startaddr, startport := validators.CheckURL(addr)

	// Server address.
	config.Address = startaddr + ":" + startport

	// JWT private key for users auth.
	jwt, exist := os.LookupEnv(("GK_JWT"))
	if !exist {
		fmt.Println("ENV GK_JWT not set")
	}
	config.PathJWT = jwt

	// TLS Private key.
	tlsPk, exist := os.LookupEnv(("GK_TLS_PK"))
	if !exist {
		fmt.Println("ENV GK_TLS_PK not set")
	}
	config.TLSPK = tlsPk

	// GK_TLS_PUB_SERV Public key.
	tlsPub, exist := os.LookupEnv(("GK_TLS_PUB_SERV"))
	if !exist {
		fmt.Println("ENV GK_TLS_PUB_SERV not set")
	}
	config.TLSPub = tlsPub

	// MINIO Config properties.
	bk, exist := os.LookupEnv(("GK_MINIO_BACKET"))
	if !exist {
		fmt.Println("ENV GK_MINIO_BACKET not set")
	}

	id, exist := os.LookupEnv(("GK_MINIO_ID"))
	if !exist {
		fmt.Println("ENV GK_MINIO_ID not set")
	}

	sc, exist := os.LookupEnv(("GK_MINIO_SECRET"))
	if !exist {
		fmt.Println("ENV GK_MINIO_SECRET not set")
	}

	config.Backetmi = bk
	config.IDmi = id
	config.Secretmi = sc

	// DSN postgres URI.
	dsn, exist := os.LookupEnv(("GK_DSN_URI"))
	if !exist {
		fmt.Println("ENV GK_DSN_URI not set")
	}
	config.DSN = dsn

	// Run local, no printing logs to file.
	_, exist = os.LookupEnv(("GK_ZAP_LOCAL"))
	if exist {
		config.RunLocal = true
	}

	return config
}
