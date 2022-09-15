package main

/*

Grants need before running:

  grant aq_administrator_role to cj;
  grant aq_user_role to cj;

*/

import (
  "os"
  "context"
  "database/sql"
  _ "errors"
  "fmt"
  "github.com/godror/godror"
  "time"
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

  ctx := context.Background()
  tx, err := db.BeginTx(ctx, nil)
  if err != nil {
    panic(err)
  }
  defer tx.Rollback()

  stmt := `DECLARE
             tbl CONSTANT VARCHAR2(61) := USER||'.TEST_Q_TBL';
             q   CONSTANT VARCHAR2(61) := USER||'.TEST_Q';
           BEGIN
             BEGIN SYS.DBMS_AQADM.stop_queue(q); EXCEPTION WHEN OTHERS THEN NULL; END;
             BEGIN SYS.DBMS_AQADM.drop_queue(q); EXCEPTION WHEN OTHERS THEN NULL; END;
             BEGIN SYS.DBMS_AQADM.drop_queue_table(tbl); EXCEPTION WHEN OTHERS THEN NULL; END;

             SYS.DBMS_AQADM.CREATE_QUEUE_TABLE(tbl, 'RAW');
             SYS.DBMS_AQADM.CREATE_QUEUE(q, tbl);

             SYS.DBMS_AQADM.grant_queue_privilege('ENQUEUE', q, USER);
             SYS.DBMS_AQADM.grant_queue_privilege('DEQUEUE', q, USER);

             SYS.DBMS_AQADM.start_queue(q);
           END;`

  if _, err = tx.ExecContext(ctx, stmt); err != nil {
    panic(err)
  }

  defer func() {
    stmt := `DECLARE
               tbl CONSTANT VARCHAR2(61) := USER||'.TEST_Q_TBL';
               q   CONSTANT VARCHAR2(61) := USER||'.TEST_Q';
             BEGIN
               BEGIN SYS.DBMS_AQADM.stop_queue(q); EXCEPTION WHEN OTHERS THEN NULL; END;
               BEGIN SYS.DBMS_AQADM.drop_queue(q); EXCEPTION WHEN OTHERS THEN NULL; END;
               BEGIN SYS.DBMS_AQADM.drop_queue_table(tbl); EXCEPTION WHEN OTHERS THEN NULL; END;
             END;`
    db.ExecContext(context.Background(), stmt)
  }()

  q, err := godror.NewQueue(ctx, tx, "TEST_Q", "", godror.WithEnqOptions(godror.EnqOptions{
    Visibility:   godror.VisibleOnCommit,
    DeliveryMode: godror.DeliverPersistent,
  }))
  if err != nil {
    panic(err)
  }
  defer q.Close()

  if err := putmessages(ctx, db, "TEST_Q"); err != nil {
    panic(err)
  }
  if err = readmessages(ctx, db, "TEST_Q"); err != nil {
    panic(err)
  }

}

func putmessages(ctx context.Context, db *sql.DB, dbQueue string) error {

  tx, err := db.BeginTx(ctx, nil)
  if err != nil {
    return err
  }
  q, err := godror.NewQueue(ctx, tx, dbQueue, "")   // empty string for RAW queue.
                                                    // Alternatively an object name
  if err != nil {
    return err
  }
  defer q.Close()
  defer tx.Commit()

  // Put some messages into the queue
  msgs := make([]godror.Message, 2)
  msg1 := "Hello, world!"
  msg2 := "How are you?"

  msgs[0] = godror.Message{
    Expiration: 10 * time.Second,
    Raw:        []byte(msg1),
  }
  msgs[1] = godror.Message{
    Expiration: 10 * time.Second,
    Raw:        []byte(msg2),
  }

  if err := q.Enqueue(msgs); err != nil {
    panic(err)
  }

  fmt.Printf("Enqueued message 1: %+v\n", msgs[0])
  fmt.Printf("Enqueued message 2: %+v\n", msgs[1])

  return nil
}

func readmessages(ctx context.Context, db *sql.DB, dbQueue string) error {
  
  tx, err := db.BeginTx(ctx, nil)
  if err != nil {
    return err
  }
  q, err := godror.NewQueue(ctx, tx, dbQueue, "")
  if err != nil {
    return err
  }
  defer q.Close()
  defer tx.Commit()

  msgs := make([]godror.Message, 2)
  n, err := q.Dequeue(msgs)
  if err != nil {
    panic(err)
  }

  for _, m := range msgs[:n] {
    fmt.Println("Dequeued message:", string(m.Raw))
    err = m.Object.Close()
    if err != nil {
      return err
    }
  }
  return nil
}
