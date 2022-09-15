
package main

// DML

import (
  "os"
  "fmt"
  "bytes"
  "strings"
  "io"
  "io/ioutil"
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

  setup(db)
  simple(db)
  reader(db)

}

// Create a table

func setup(db *sql.DB) {
  db.Exec("drop table mytest")
  db.Exec("create table mytest (k number, c clob, b blob)")
}


// Insert and query LOBs as direct strings or bytes

func simple(db *sql.DB) {

  fmt.Println("Simple")

  insSql := `insert into mytest (k, c, b) values (:1, :2, :3)`

  bIn := []byte{0, 1, 2, 3, 4}
  sIn := "abcdef"

  if _, err := db.Exec(insSql, 1, sIn, bIn); err != nil {
    panic(err)
  }
  
  rows, _ := db.Query("select c, b from mytest where k = 1")
  defer rows.Close()

  var c string
  var b []byte
  for rows.Next() {
    rows.Scan(&c, &b)
    fmt.Printf("  CLOB=%s BLOB=%v\n", c, b)
  }

}

// Insert and query LOBs using Streaming

func reader(db *sql.DB) {

  fmt.Println("Streaming")

  // Streaming insert

  insSql := `insert into mytest (k, c, b) values (:1, :2, :3)`

  b := []byte{5, 6, 7, 8, 9}
  s := "ghijkl"

  if _, err := db.Exec(insSql,
    2,
    godror.Lob{Reader: strings.NewReader(s), IsClob: true},
    godror.Lob{Reader: bytes.NewReader(b)},
  ); err != nil {
    panic(err)
  }

  // Streaming query
  
  rows, _ := db.Query("select c, b from mytest where k = 2", godror.LobAsReader())
  defer rows.Close()
  
  var clobCol interface{}
	var blobCol interface{}
  for rows.Next() {
    rows.Scan(&clobCol, &blobCol)

		clob := clobCol.(*godror.Lob)
		clobdata, _ := ioutil.ReadAll(clob)
    
		blob := blobCol.(*godror.Lob)
		blobdata := make([]byte, 5)   // 5 is length of byte array inserted
		io.ReadFull(blob, blobdata)
    
    fmt.Printf("  CLOB=%s BLOB=%v\n", clobdata, blobdata)

  }
}
