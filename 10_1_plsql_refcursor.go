package main

// REF CURSOR

import (
  "os"
  "fmt"
  "time"
  "context"
  "database/sql"
  "database/sql/driver"
  godror "github.com/godror/godror"
)

func main() {

  dsn := `user="cj"
          password="cj"
          connectString=` + os.Getenv("GODROR_CONNECTSTRING") +   // like "example.com/orclpdb1"
        ` libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  setup(db)
  call(db)
}

func setup(db *sql.DB) {

  _, err := db.Exec(
    `create or replace procedure myproc(id in number, rc out sys_refcursor)
     as
     begin
       open rc for
         select first_name, last_name
         from employees
         where employee_id <= id
         order by employee_id;
     end;`)
  if err != nil {
    panic(err)
  }

}

func call(db *sql.DB) {

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  const plsql = "begin myproc(:1, :2); end;"
  const maxEmpId = 110
  var   rc driver.Rows
  if _, err := db.Exec(plsql, maxEmpId, sql.Out{Dest: &rc}); err != nil {
    panic(err)
  }
  defer rc.Close()

  sub, err := godror.WrapRows(ctx, db, rc)  // transform driver.Rows into *sql.Rows
  if err != nil {
    panic(err)
  }
  defer sub.Close()

  for sub.Next() {
    var firstName string
    var lastName string
    if err := sub.Scan(&firstName, &lastName); err != nil {
      panic(err)
    }
    fmt.Println(firstName, lastName)
  }

}
