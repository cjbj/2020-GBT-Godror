package main

// "Paged" query

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

  myoffset     := 0   // do not skip any rows (start at row 1)
  mymaxnumrows := 10  // get 10 rows

  for ; myoffset < 30; myoffset += 10 {
    fmt.Println("Getting rows")
    printrows(db, myoffset, mymaxnumrows)
  }

}

func printrows(db *sql.DB, myoffset int, mymaxnumrows int) {

  query := `select last_name
            from employees
            order by last_name
            offset :offset rows fetch next :maxnumrows rows only`

  rows, err := db.Query(query, myoffset, mymaxnumrows,
    godror.PrefetchCount(mymaxnumrows+1), godror.FetchArraySize(mymaxnumrows))

  if err != nil {
    panic(err)
  }
  defer rows.Close()

  var last_name string
  for rows.Next() {
    err := rows.Scan(&last_name)
    if err != nil {
      panic(err)
    }
    fmt.Println("   " + last_name)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }

}
