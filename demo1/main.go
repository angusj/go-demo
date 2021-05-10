package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"os"
)

func main() {
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		fmt.Printf("Original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("Stack trace:\n%+v\n", err)
		os.Exit(1)
	}
	type row struct {
		age  int
		name string
	}
	var r row
	err = db.QueryRow("select age, name from users where age = 2").Scan(&r.age, &r.name)
	if err != nil {
		fmt.Printf("Original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("Stack trace:\n%+v\n", err)
		os.Exit(1)
	}
}
