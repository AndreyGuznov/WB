package db

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"weather/tool/config"
	"weather/tool/log"
)

// DB table names
const (
	LocationsTable = "locations"
	ForecastTable  = "forecast"
)

var (
	sqlUser     = config.Get("USER")
	sqlPassword = config.Get("PASSWORD")
	sqlDatabase = config.Get("NAME")
	sqlHost     = config.Get("HOST")
	conn        *SQLDB
	connMx      = &sync.Mutex{}
)

// GetConn ...
func GetConn() *SQLDB {
	connMx.Lock()
	defer connMx.Unlock()

	if conn == nil {
		var err error

		conn, err = Init(sqlUser, sqlPassword, sqlDatabase, sqlHost)
		if err != nil {
			log.Err("Error initializing DB", err)
			return nil
		}
	}

	return conn
}

// SQLDB ...
type SQLDB struct {
	Instance *sqlx.DB
}

func doInit(user string, password string, database string, host string) (*SQLDB, error) {
	addr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, 5432, user, password, database)

	instance, err := sqlx.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db := SQLDB{
		Instance: instance,
	}

	return &db, nil
}

// Init inits DB
func Init(user string, password string, database string, host string) (*SQLDB, error) {
	var dbsess *SQLDB
	var err error

	dbsess, err = doInit(user, password, database, host)
	if err != nil {
		log.Err("Failed to connect to Postgres", err)
		return nil, err
	}

	log.Info("Connected to Postgres database")
	return dbsess, nil
}

// Close DB connection
func (db *SQLDB) Close() error {
	return db.Instance.Close()
}
