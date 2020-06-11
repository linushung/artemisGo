package postgres

import (
	/* Ref: http://go-database-sql.org/index.html */
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/linushung/artemis/internal/pkg/configs"
	/* Ref: http://jmoiron.github.io/sqlx/ */
	"github.com/jmoiron/sqlx"
	// import postgres Driver for database/sql package
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type RDB struct {
	Type string
	Host string
	/*
		1. DB object is a pool of many database connections which contains both 'open' and 'idle' connections. A connection
		is marked open when you are using it to perform a database task, such as executing a SQL statement or querying rows.
		When the task is complete the connection becomes idle.
		2. The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
		Thus, the Open function should be called just once. It is rarely necessary to close a DB.
	*/
	Pool  *sql.DB
	Poolx *sqlx.DB
}

// InitPostgreSQL create an abstraction representing a Database (*sqlx.DB) and verify with a ping
func InitPostgreSQL() RDB {
	dbType := "PostgreSQL"

	host := configs.GetConfigStr("connection.rdb.host")
	username := configs.GetConfigStr("connection.rdb.username")
	password := configs.GetConfigStr("connection.rdb.password")
	db := configs.GetConfigStr("connection.rdb.database")

	connsPool, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, db))
	if err != nil {
		log.Fatalf("***** [DATABASE][FAIL] ***** Failed to create connection to PostgreSQL::%s %s", host, err)
		os.Exit(1)
	}

	/* Ref: https://www.alexedwards.net/blog/configuring-sqldb */
	/* Set the maximum number of concurrently open connections (in-use + idle) to 5. Setting this to less than or equal
	to 0 will mean there is no maximum limit. (Default setting is no limit) */
	connsPool.SetMaxOpenConns(5)
	/* Set the maximum number of concurrently idle connections to 5. Setting this to less than or equal to 0 will mean
	that no idle connections are retained. (Default setting is 2)
	MaxIdleConns should always be less than or equal to MaxOpenConns. Go enforces this and will automatically reduce
	MaxIdleConns if necessary.*/
	connsPool.SetMaxIdleConns(5)
	/* Set the maximum lifetime of a connection to 1 hour. Setting it to 0 means that there is no maximum lifetime and
	the connection is reused forever (Default setting is no limit). */
	connsPool.SetConnMaxLifetime(time.Hour)
	log.Infof("***** [DATABASE:%s] ***** Create connections to PostgreSQL::%s!", dbType, host)

	return RDB{
		Type:  dbType,
		Host:  host,
		Poolx: connsPool,
	}
}

func (rdb *RDB) transactionHandler(ops string, block func(tx *sqlx.Tx)) error {
	tx, err := rdb.Poolx.Beginx()
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute BEGIN(Transaction) operation:: %#v", ops, err)
		return err
	}

	defer recoverTransaction(ops, tx)
	block(tx)

	if err := tx.Commit(); err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute COMMIT(Transaction) operations:: %#v", ops, err)
		return err
	}

	return nil
}

/* Ref:
1. https://blog.golang.org/defer-panic-and-recover
2. https://eli.thegreenplace.net/2018/on-the-uses-and-misuses-of-panics-in-go/
*/
func recoverTransaction(ops string, tx *sqlx.Tx) {
	if p := recover(); p != nil {
		log.Errorf("***** [PANIC:%s] ***** Capture PANIC during DB Transaction:: %#v", ops, p)
		tx.Rollback()
	}
}
