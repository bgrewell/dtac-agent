package authndb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
	"os"
	"time"
)

// User is the struct for a user
type User struct {
	ID       uint64   `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}

// TokenDetails is the struct for the token details
type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

// AccessDetails is the struct for the access details
type AccessDetails struct {
	AccessUUID string `json:"access_uuid"`
	UserID     uint64 `json:"user_id"`
}

// NewAuthDB creates a new authn database
func NewAuthDB(log *zap.Logger) *AuthDB {
	name := "authndb"
	db := AuthDB{
		Logger: log.With(zap.String("module", name)),
		bucket: name,
	}
	err := db.Initialize()
	if err != nil {
		db.Logger.Error("failed to initialize authn database", zap.String("subsystem", name), zap.Error(err))
		return nil
	}
	return &db
}

// AuthDB is the struct for the authn database
type AuthDB struct {
	Logger *zap.Logger
	DB     *bolt.DB
	bucket string
}

// Initialize initializes the authn database
func (db *AuthDB) Initialize() error {
	// Ensure db directory exits
	if _, err := os.Stat(config.GlobalDBLocation); os.IsNotExist(err) {
		err = os.MkdirAll(config.GlobalDBLocation, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create authentication database directory: %v", err)
		}
	}

	// Initialize Database
	var err error
	db.DB, err = bolt.Open(config.DBName, 0600, &bolt.Options{Timeout: 1 * time.Second})
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

// UpdateDB updates the authn database
func (db *AuthDB) UpdateDB(key string, value string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

// ViewDB views the authn database
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
