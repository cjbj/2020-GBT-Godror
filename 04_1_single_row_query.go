package main

// Single row query

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

  query := `select last_name
            from employees
            where employee_id = :id`

  var last_name string
  err = db.QueryRow(query, 100, godror.FetchArraySize(1)).Scan(&last_name)
  if err != nil {
    panic(err)
  }
  fmt.Println(last_name)

}
