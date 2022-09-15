package main

// DML RETURNING aka RETURNING INTO

import (
  "os"
  "fmt"
  "time"
  "database/sql"
  _ "github.com/godror/godror"
)

func main() {

  dsn := `user="cj"
          password="cj"
          connectString=` + os.Getenv("GODROR_CONNECTSTRING") +
        ` libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  setup(db)
  insert(db)

}

// Create a table

func setup(db *sql.DB) {
  db.Exec("drop table mytest")
  db.Exec("create table mytest (k number, n varchar2(20), d date)")
}

// "Array DML"

func insert(db *sql.DB) {

  sqlText := `insert into mytest (k, n, d)
              values (:k_bv, :n_bv, sysdate)
              returning d into :TimeInserted`

  var timeInserted time.Time

  if _, err := db.Exec(sqlText,
    sql.Named("k_bv", 101),
    sql.Named("n_bv", "Petra"),
    sql.Named("TimeInserted", sql.Out{Dest: &timeInserted})); err != nil {
    panic(err)
  }
  fmt.Printf("Time inserted was %v\n", timeInserted)

}
