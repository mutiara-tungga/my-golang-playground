package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type MySQLDBOption struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
	Env      string
}

func main() {
	mySQLOpt := MySQLDBOption{
		Username: "root",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "test_tx_mutia",
		Charset:  "utf8mb4",
		Env:      "",
	}
	dbURL := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True", mySQLOpt.Username, mySQLOpt.Password, mySQLOpt.Host, mySQLOpt.Port, mySQLOpt.Database, mySQLOpt.Charset)
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	txHandler := NewTxHandler(db)
	test1 := NewTest1(db)
	test2 := NewTest2(db)
	test3 := NewTest3(db)

	err = txHandler.WithTransaction(func(tx *sql.Tx) error {
		name1 := "Test 1 1"
		err = test1.InsertTx(tx, name1)
		if err != nil {
			return err
		}

		name2 := "Test 2 1"
		err = test2.InsertTx(tx, name2)
		if err != nil {
			return err
		}

		name3 := "Test 3 1"
		err = test3.InsertTx(tx, name3)
		if err != nil {
			return err
		}

		return nil
	})

	fmt.Println("error", err)
}

type Test1 struct {
	db *sql.DB
}

func NewTest1(db *sql.DB) *Test1 {
	return &Test1{db: db}
}

func (t *Test1) InsertTx(tx *sql.Tx, name string) error {
	stmt, err := tx.Prepare("INSERT INTO test1(name) values(?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

type Test2 struct {
	db *sql.DB
}

func NewTest2(db *sql.DB) *Test2 {
	return &Test2{db: db}
}

func (t *Test2) InsertTx(tx *sql.Tx, name string) error {
	stmt, err := tx.Prepare("INSERT INTO test2(name) values(?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

type Test3 struct {
	db *sql.DB
}

func NewTest3(db *sql.DB) *Test3 {
	return &Test3{db: db}
}

func (t *Test3) InsertTx(tx *sql.Tx, name string) error {
	stmt, err := tx.Prepare("INSERT INTO test3(name) values(?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

type TxFn func(*sql.Tx) error

type TxHandler interface {
	WithTransaction(TxFn) error
}

type txHandler struct {
	*sql.DB
}

func NewTxHandler(db *sql.DB) TxHandler {
	return &txHandler{db}
}

func (th *txHandler) WithTransaction(fn TxFn) (err error) {
	tx, err := th.Begin()

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrapf(err, "fail rollback panic: %+v", rollbackErr)
			}
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrapf(err, "fail rollback error: %+v", rollbackErr)
			}
			return
		}

		// all good, commit
		err = errors.Wrap(tx.Commit(), "fail commit")
	}()

	if err != nil {
		return
	}

	err = fn(tx)
	return
}
