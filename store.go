package main

import (
	"errors"
	bolt "go.etcd.io/bbolt"
	"log"
	"strconv"
)

var SizeBucket = []byte("SizeBucket")

var HashBucket = []byte("HashBucket")

type DBI struct {
	db        *bolt.DB
	tx        *bolt.Tx
	filesizes *bolt.Bucket
	hashes    *bolt.Bucket
}

func dbInit(dbFileName string) (*DBI, error) {
	if len(dbFileName) == 0 {
		return nil, errors.New("Bad file name:[" + dbFileName + "]")
	}
	var err error
	dbi := new(DBI)

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbi.db, err = bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		log.Println("Error opening db file:", dbFileName)
		log.Println(err)
		return nil, err
	}

	dbi.tx, err = begin(dbi.db)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dbi.hashes, err = newBucket(dbi.tx, HashBucket)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dbi.filesizes, err = newBucket(dbi.tx, SizeBucket)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Write the buckets
	err = commitAndBeginAndBuckets(dbi)
	if err != nil {
		return nil, err
	}
	return dbi, nil
	//return nil, nil
}

func commitAndBeginAndBuckets(dbi *DBI) error {
	err := commit(dbi.tx)
	if err != nil {
		return err
	}

	dbi.tx = nil
	dbi.filesizes = nil
	dbi.hashes = nil

	// Start new transaction
	dbi.tx, err = begin(dbi.db)
	if err != nil {
		log.Println(err)
		return err
	}

	// Get the buckets back
	dbi.hashes = dbi.tx.Bucket(HashBucket)
	if dbi.hashes == nil {
		err := errors.New("Bucket does not exist")
		log.Println(err)
		return err
	}

	dbi.filesizes = dbi.tx.Bucket(SizeBucket)
	if dbi.filesizes == nil {
		err := errors.New("Bucket does not exist")
		log.Println(err)
		return err
	}
	return nil
}

func begin(db *bolt.DB) (*bolt.Tx, error) {
	log.Println("TXTXT BEGIN  ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	if db == nil {
		return nil, errors.New("DB is nil")
	}
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		log.Println("Error opening tx")
		log.Println(err)
	}
	return tx, err
}
func commit(tx *bolt.Tx) error {
	log.Println("TXTXT COMMIT  ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	if tx == nil {
		return errors.New("Tx is nil")
	}
	err := tx.Commit()
	if err != nil {
		log.Println("Error in tx commit")
		log.Println(err)
	}
	return err

}

func newBucket(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	if tx == nil {
		return nil, errors.New("Tx is nil")
	}
	if len(bucketName) == 0 {
		return nil, errors.New("Bucket is zero length")
	}

	b, err := tx.CreateBucket(bucketName)
	if err != nil {
		log.Println("Error creating bucket", bucketName)
		log.Println(err)
	}
	return b, err
}

var temp []byte

func storeSize(b *bolt.Bucket, filesize int64) error {
	bs := []byte(strconv.FormatInt(filesize, 10))
	return store(b, bs)
}

func store(b *bolt.Bucket, bs []byte) error {
	if b == nil {
		return errors.New("Bucket is nil")
	}
	if bs == nil {
		return errors.New("Content []byte is nil")
	}

	if len(bs) == 0 {
		return errors.New("Content []byte is zero length")
	}

	err := b.Put(bs, temp)
	if err != nil {
		log.Println(err)
	}
	return err
}

func pushTx(dbi *DBI) error {
	return commitAndBeginAndBuckets(dbi)
}
