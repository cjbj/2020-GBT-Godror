package main

// Nested Cursor

import (
  "os"
  "time"
  "fmt"
  "context"
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

  query(db)
}


// Query a nested cursor

func query(db *sql.DB) {

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  stmt := `select
             cursor(select first_name
             from employees
             where employee_id <= 110
             order by employee_id)
           from dual`

  rows, err := db.QueryContext(ctx, stmt)
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  for rows.Next() {

    var nestedRows sql.Rows
    if err := rows.Scan(&nestedRows); err != nil {
      panic(err);
    }
    defer nestedRows.Close()

    for nestedRows.Next() {
      var first_name string
      if err := nestedRows.Scan(&first_name); err != nil {
        panic(err);
      }
      fmt.Println(first_name)
    }
  }
}
