package main

// Transactions

import (
  "os"
  "fmt"
  "strconv"
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
  query(db)

}

// "Array DML"

func insert(db *sql.DB) {

  tx, err := db.Begin()
  if err != nil {
    panic(err)
  }

  for i := 0; i < 5; i++ {
    s := "Chris-" + strconv.Itoa(i)
    _, err := tx.Exec("insert into mytest (k, n) values (:1, :2)", i, s)
    if err != nil {
      panic(err)
    }
  }

  err = tx.Commit()
  if err != nil {
    panic(err)
  }

}

// Create a table

func setup(db *sql.DB) {
  db.Exec("drop table mytest")
  db.Exec("create table mytest (k number, n varchar2(20))")
}

// Display the table data

func query(db *sql.DB) {

  sql := `select k, n
          from mytest
          order by k`

  rows, err := db.Query(sql)
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  var key uint
  var name string
  for rows.Next() {
    err := rows.Scan(&key, &name)
    if err != nil {
      panic(err)
    }
    fmt.Println(key, name)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }

}
