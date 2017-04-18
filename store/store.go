package store

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"log"
)

var Db *bolt.DB 

func InitDB() (*bolt.DB) {
// Open the blog.db data file in your current directory.
// It will be created if it doesn't exist.
  Db, err := bolt.Open("blog.db", 0600, nil)
  if err != nil {
    log.Fatal(err)
  }
  return Db
}

// itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
  b := make([]byte, 8)
  binary.BigEndian.PutUint64(b, uint64(v))
return b
}

func GetNextId (bucketName string) (int,error) {
  // Open a writable transaction
  tx, err := Db.Begin(true)
  if err != nil {
    return -1,err
  }
  defer tx.Rollback() 

  b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
  seq, _ := b.NextSequence()

  // Commit changes 
  if err:= tx.Commit(); err != nil {
    return  -1,err 
  }
  return int(seq),nil
}

func ReadUnique (bucketName string, key int) ([]byte,error) {
  var buf[]byte

  Db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
      // Get only one key
      buf = b.Get(Itob(key))
      return nil
  })
  return buf,nil
}

func ReadRangeFromSet (parentBucketName string,detailBucketName string) ([][]byte,error) {
  var buf [][]byte
  var elem []byte

  Db.View(func(tx *bolt.Tx) error {
    p := tx.Bucket([]byte(parentBucketName))
    d := tx.Bucket([]byte(detailBucketName))
      // Iterate results
      c := p.Cursor()

      for k, v := c.First(); k != nil; k, v = c.Next() {
        elem = d.Get(v)
        buf= append(buf,elem)
      }
      return nil
  })
  return buf,nil
}

func ReadRange (bucketName string) ([][]byte,error) {
  var buf [][]byte

  Db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
      // Iterate results
      c := b.Cursor()

      for k, v := c.First(); k != nil; k, v = c.Next() {
        buf= append(buf,v)
      }
      return nil
  })
  return buf,nil
}

func Remove (bucketName string, key int) (err error) {
  
  Db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    err = b.Delete(Itob(key))
    return err
  })
  return
}

func Create (bucketName string, buf []byte, key int) error {
  // Open a writable transaction
  tx, err := Db.Begin(true)
  if err != nil {
    return err
  }
  defer tx.Rollback() 

  // Get bucket
  b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
  // Persist bytes to bucket.
  err = b.Put(Itob(key), buf)
  if  err != nil {
    return err
  }
  // Commit changes 
  if err:= tx.Commit(); err != nil {
    return  err 
  }
  return nil
}

