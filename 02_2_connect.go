package main

import (
  "fmt"
  "database/sql"
  _ "github.com/godror/godror"
)

func main() {

  // https://download.oracle.com/ocomdocs/global/Oracle-Net-19c-Easy-Connect-Plus.pdf
  
  const dsn = `user="cj"
               password="cj"
               connectString="localhost/orclpdb1?connect_timeout=2&expire_time=2"
               libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()     // gets, uses & releases a connection
  if err != nil {
    fmt.Println("Ooops! Ping failed")
    panic(err)
  }

  fmt.Println("Pinged DB OK")
}
