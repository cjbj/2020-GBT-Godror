package main

// Concurrent DML

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
        ` poolMinSessions=5
          poolMaxSessions=5
          poolIncrement=0
          libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  db.SetMaxIdleConns(0)  // don't use database/sql connection pool
  
  setup(db)

  const numBatches = 5
  c := make(chan int)
  for i := 0; i < numBatches; i++ {
    go insert(db, i, c)
  }

  for i := 0; i < numBatches; i++ {
    <- c
  }


  query(db)

}

// Insert a row using a channel

func insert(db *sql.DB, i int, c chan int) {

  s := "Angela-" + strconv.Itoa(i)
  if _, err := db.Exec(`insert into mytest values (:1, :2)`, i, s); err != nil {
    panic(err)
  }
  c <- 1

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
