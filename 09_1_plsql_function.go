package main

// PL/SQL Function call

import (
  "os"
  "fmt"
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

  _, err = db.Exec(
    `create or replace function myfunc(p1 in varchar2, p2 in out varchar2)
       return varchar2 as
     begin
       p2 := p1 || ' ' || p2;
       return p1 || ' Jane';
     end;`)
  if err != nil {
    panic(err)
  }

  defer db.Exec(`drop function myfunc`)

  var res string
  plsql := `begin :1 := myfunc(:2, :3); end;`

  p2 := "Smith"
  db.Exec(plsql, sql.Out{Dest: &res}, "Kylie", sql.Out{Dest: &p2, In: true})
  fmt.Println(res)
  fmt.Println(p2)

}
