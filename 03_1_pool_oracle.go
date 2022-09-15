package main

// Oracle Session "Connection" Pool (MaxIdleConns=0)

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

        ` poolMinSessions=4
          poolMaxSessions=4
          poolIncrement=0

          libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  db.SetMaxIdleConns(0)  // don't also use database/sql connection pool

  err = db.Ping()
  if err != nil {
    fmt.Println("Ooops! Ping failed")
    panic(err)
  }

  fmt.Println("Pinged DB OK")
}
