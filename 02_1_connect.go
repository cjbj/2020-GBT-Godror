package main

/*
 Install godror:

   go get github.com/godror/godror

 Install Oracle Client, e.g. on macOS use Oracle Instant Client:

   cd $HOME/Downloads
   curl -O https://download.oracle.com/otn_software/mac/instantclient/198000/instantclient-basic-macos.x64-19.8.0.0.0dbru.dmg
   hdiutil mount instantclient-basic-macos.x64-19.8.0.0.0dbru.dmg
   /Volumes/instantclient-basic-macos.x64-19.8.0.0.0dbru/install_ic.sh
   hdiutil unmount /Volumes/instantclient-basic-macos.x64-19.8.0.0.0dbru

*/

import (
  "database/sql"
  _ "github.com/godror/godror"
)

func main() {

  // libDir is usable on macOS and Windows
  
  const dsn = `user="cj"
               password="cj"
               connectString="localhost/orclpdb1"
               libDir="/Users/cjones/Downloads/instantclient_19_8"`

  db, err := sql.Open("godror", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  print("OK\n")
}
