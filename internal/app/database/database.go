package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// RDB represents a relational database
type RDBx struct {
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
