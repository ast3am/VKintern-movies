package db

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/jackc/pgx/v4"
	"time"
)

//go:generate mockery --name logger
type logger interface {
	DebugMsg(msg string)
	ErrorMsg(msg string, err error)
}

type DB struct {
	dbConnect *pgx.Conn
	log       logger
}

func NewClient(ctx context.Context, cfg *models.Config, log logger) (*DB, error) {
	DB := DB{dbConnect: nil, log: log}
	var err error
	posgresURL := "postgresql://" + cfg.SqlConfig.UsernameDB + ":" + cfg.SqlConfig.PasswordDB + "@" + cfg.SqlConfig.HostDB + ":" + cfg.SqlConfig.PortDB + "/" + cfg.SqlConfig.DBName
	for i := 0; i < cfg.SqlConfig.DelayTime; i++ {
		DB.dbConnect, err = pgx.Connect(ctx, posgresURL)
		if err != nil {
			time.Sleep(1 * time.Second)
			log.DebugMsg("fail connect to DB, try again")
			continue
		}
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	err = DB.dbConnect.Ping(ctx)
	if err != nil {
		return nil, err
	}
	log.DebugMsg("connection to DB is OK")
	return &DB, nil
}

func (db *DB) Close(ctx context.Context) {
	err := db.dbConnect.Close(ctx)
	if err != nil {
		db.log.ErrorMsg("error closing connection to BD", err)
	} else {
		db.log.DebugMsg("closing connection to BD is OK")
	}
}
