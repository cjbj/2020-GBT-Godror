package main

// DML

import (
  "os"
  "fmt"
  "time"
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

// Create a table

func setup(db *sql.DB) {
  db.Exec(`drop table mytest`)
  db.Exec(`create table mytest (k number, n varchar2(20), d date)`)
}

// Insert using "Array DML"

func insert(db *sql.DB) {

  const numRows = 10
  intVals := make([]int, numRows)
  strVals := make([]string, numRows)
  dateVals := make([]time.Time, numRows)
  for i := range intVals {
    intVals[i] = i
    strVals[i] = "Chris-" + strconv.Itoa(i)
    dateVals[i] = time.Now()
  }

  // Note: autocommits
  if _, err := db.Exec(`insert into mytest values (:1, :2, :3)`, intVals, strVals, dateVals); err != nil {
    panic(err)
  }

}

// Display the table data

func query(db *sql.DB) {
  sql := `select k, n, d
          from mytest
          order by k`

  rows, err := db.Query(sql)
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  var key uint
  var name string
  var date time.Time
  for rows.Next() {
    err := rows.Scan(&key, &name, &date)
    if err != nil {
      panic(err)
    }
    fmt.Println(key, name, date)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }

}
