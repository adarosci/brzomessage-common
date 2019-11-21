package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repository db connectionn
type Repository struct {
	Conn *sqlx.DB
	Tx   *sqlx.Tx
}

// EventHandler call execute in transaction
type EventHandler func() error

var _db = "postgres"
var _conn *sqlx.DB

func newConnection() *sqlx.DB {
	db, err := sqlx.Connect(_db, os.Getenv("CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Minute * 5)
	return db
}

func open() *sqlx.DB {
	if _conn == nil {
		_conn = newConnection()
	}
	return _conn
}

// ExecInTransaction executa conjunto de comandos em transacao
func (r *Repository) ExecInTransaction(eh EventHandler) error {
	//if r.Conn == nil {
	r.Conn = open()
	//}
	//defer r.Conn.Close()
	r.Tx, _ = r.Conn.Beginx()
	err := eh()
	if err != nil {
		r.Tx.Rollback()
		return err
	}
	err = r.Tx.Commit()
	r.Tx = nil
	return err
}

// Query return select query db
func (r *Repository) Query(rs interface{}, query string, param ...interface{}) error {
	str := reflect.TypeOf(rs).String()
	if strings.Contains(str, "[]") {
		return getMany(rs, query, param...)
	}
	return getOne(rs, query, param...)
}

// GetOne return select query db multiples rows
func getMany(rs interface{}, query string, param ...interface{}) error {
	conn := open()
	// defer conn.Close()
	if len(param) == 0 {
		return conn.Select(rs, query)
	}
	rows, err := conn.NamedQuery(query, param[0])
	if err != nil {
		fmt.Println(err)
	}
	if rows != nil {
		return sqlx.StructScan(rows, rs)
	}
	return nil
}

// GetOne return select query db single rows
func getOne(rs interface{}, query string, param ...interface{}) (er error) {
	conn := open()
	// defer conn.Close()
	if param == nil {
		return conn.Get(rs, query)
	}
	rows, err := conn.NamedQuery(query, param[0])
	if err != nil {
		fmt.Println(err)
		return err
	}
	for rows.Next() {
		er = rows.StructScan(rs)
	}
	return
}

// Exec exec sql insert update delete
func (r *Repository) Exec(rs interface{}, query string, param interface{}) error {
	if r.Tx == nil {
		return errors.New("execute query, INSERT, UPDATE OR DELETE in transaction!")
	}
	if rs == nil {
		return r.execNoResult(query, param)
	}
	return r.execResult(rs, query, param)
}

func (r *Repository) execResult(rs interface{}, query string, param interface{}) error {
	smtp, err := r.Tx.PrepareNamed(query)
	if err != nil {
		return err
	}
	return smtp.Get(rs, param)
}

func (r *Repository) execNoResult(query string, param interface{}) error {
	_, err := r.Tx.NamedExec(query, param)
	return err
}
