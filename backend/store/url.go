package store

import (
	"fmt"
	"log"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

var urlsBucketname = "Urls"

type UrlStore struct {
	db *bolt.DB
}

func storePath() string {
	storePath := os.Getenv("DB_PATH")
	if storePath == "" {
		return "link-shortener.db"
	}
	return storePath
}

func NewUrlStore() UrlStore {
	store := UrlStore{}
	storePath := storePath()
	db, err := bolt.Open(storePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	store.db = db

	store.createBucket(urlsBucketname)

	return store
}

func (o *UrlStore) createBucket(bucketName string) {
	o.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (o *UrlStore) AddUrl(shortUrl string, fullUrl string) {
	o.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(urlsBucketname))
		err := b.Put([]byte(shortUrl), []byte(fullUrl))
		return err
	})
}

func (o *UrlStore) GetUrl(shortUrl string) (string, error) {
	var fullUrl string
	err := o.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(urlsBucketname))
		fullUrl = string(b.Get([]byte(shortUrl)))
		return nil
	})
	return fullUrl, err
}
