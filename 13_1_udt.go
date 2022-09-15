package main

// Oracle Named Objects aka User Defined Types

import (
  "os"
  "time"
  "fmt"
  "context"
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
  insert(db)
  fetch(db)
}

func setup(db *sql.DB) {

  db.Exec("drop table books")

  if _, err := db.Exec(
    `create or replace type book_t as object (
      title        varchar2(100),
      authors      varchar2(100),
      price        number(5,2))`); err != nil {
    panic(err)
  }

  if _, err :=  db.Exec(
    `create table books (
     id           number(9) not null,
     book         book_t not null)`); err != nil {
    panic(err)
  }

}

func insert(db *sql.DB) {

  tx, _ := db.Begin()
  ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
  objType, _ := godror.GetObjectType(ctx, tx, "BOOK_T")
  obj, _ := objType.NewObject()
  defer obj.Close()
  obj.Set("TITLE", "The Fellowship of the Ring")
  obj.Set("AUTHORS", "J.R.R. Tolkien")
  obj.Set("PRICE", 12.50)
  tx.Exec("insert into books values (:1, :2)", 1, obj)
  tx.Commit()

}

func fetch(db *sql.DB) {

  var idVal int
  var obj *godror.Object
  row := db.QueryRow("select id, book from books")
  row.Scan(&idVal, &obj)
  defer obj.Close()
  authors, _ := obj.Get("AUTHORS")
  title, _   := obj.Get("TITLE")
  price, _   := obj.Get("PRICE")
  fmt.Printf("Authors: %s\nTitle: %s\nPrice: %.2f\n", authors, title, price)

}
