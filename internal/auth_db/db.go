package auth_db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
	"os"
	"time"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

type AccessDetails struct {
	AccessUuid string `json:"access_uuid"`
	UserId     uint64 `json:"user_id"`
}

func NewAuthDB(log *zap.Logger) *AuthDB {
	name := "auth_db"
	db := AuthDB{
		Logger: log.With(zap.String("module", name)),
		bucket: name,
	}
	err := db.Initialize()
	if err != nil {
		db.Logger.Error("failed to initialize auth database", zap.String("subsystem", name), zap.Error(err))
		return nil
	}
	return &db
}

type AuthDB struct {
	Logger *zap.Logger
	DB     *bolt.DB
	bucket string
}

func (db *AuthDB) Initialize() error {
	// Ensure db directory exits
	if _, err := os.Stat(config.GLOBAL_DB_LOCATION); os.IsNotExist(err) {
		err = os.MkdirAll(config.GLOBAL_DB_LOCATION, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create authentication database directory: %v", err)
		}
	}

	// Initialize Database
	var err error
	db.DB, err = bolt.Open(config.DB_NAME, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open authentication database: %v", err)
	}

	// Ensure that the bucket exists
	return db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(db.bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %s", err)
		}
		return nil
	})
}

func (db *AuthDB) UpdateDB(key string, value string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func (db *AuthDB) ViewDB(key string) (value string, err error) {
	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key %s not found in bucket", key)
		}
		// Set value to the string representation of v
		value = string(v)
		// Return no error
		return nil
	})
	return value, err
}
