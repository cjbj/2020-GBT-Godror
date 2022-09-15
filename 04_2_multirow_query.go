package main

// Multi-row query for a large number of rows (many hundreds? thousands?)

// Tuning goal: reduce round-trips between Go and Oracle Database

import (
  "os"
  "fmt"
  "database/sql"
  godror "github.com/godror/godror"
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

  query := `select employee_id, last_name 
            from employees 
            where employee_id <= :id 
            order by employee_id`

  rows, err := db.Query(query, 120,
    godror.PrefetchCount(1000), godror.FetchArraySize(1000))  // internal tuning parameters
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  var employee_id uint
  var last_name   string
  for rows.Next() {
    err := rows.Scan(&employee_id, &last_name)
    if err != nil {
      panic(err)
    }
    fmt.Println(employee_id, last_name)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }

}
