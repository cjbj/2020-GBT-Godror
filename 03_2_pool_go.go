package main

// database/sql Connection Pool (standaloneConnection=1 MaxIdleConns > 0)

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

        ` standaloneConnection=1

          libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  db.SetMaxOpenConns(4)
  db.SetMaxIdleConns(4)

  err = db.Ping()
  if err != nil {
    fmt.Println("Ooops! Ping failed")
    panic(err)
  }

  fmt.Println("Pinged DB OK")
}
