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
		Password: "rootpw",
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

	err = txHandler.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare("INSERT INTO test1(name) values(?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec("Test")
		if err != nil {
			return err
		}

		stmt2, err := tx.Prepare("INSERT INTO test2(name) values(?)")
		if err != nil {
			return err
		}

		_, err = stmt2.Exec("Test")
		if err != nil {
			return err
		}

		stmt3, err := tx.Prepare("INSERT INTO test3(name) values(?)")
		if err != nil {
			return err
		}

		_, err = stmt3.Exec("Test")
		if err != nil {
			return err
		}

		return nil
	})

	fmt.Println("error", err)

	// hyperlocalUserId := rand.Intn(50)
	// areaID := rand.Intn(10)
	// price := rand.Intn(10000)
	// condition := rand.Intn(1)
	// categoryId := rand.Intn(100)
	// productState := rand.Intn(5)
	// // row, err := self.DB.Exec(QueryHyperlocalProducts["insert_product"], hyperlocalUserId, areaId, form.Title, form.Price, form.Description, model.ProductConditionState(form.Condition), form.CategoryId, 1, model.ProductStateInt("pending_review"), 0, "")
	// sql := "INSERT INTO hyperlocal_products(`hyperlocal_user_id`, `hyperlocal_id`, `title`, `price`, `description`, `condition`, `category_id`, `active`, `state`, `rejected_id`, `rejected_remark`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
	// db.Exec(sql, hyperlocalUserId, areaID, "Test", price, "Description", condition, categoryId, 1, productState, "")

	// a, b := strconv.ParseUint("a", 10, 32)
	// fmt.Println(a)
	// fmt.Println(reflect.TypeOf(a))
	// fmt.Println(b)
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
