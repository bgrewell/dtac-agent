package authndb

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
	"os"
	"time"
)

// UserArgs is the struct for the user arguments validation
type UserArgs struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User is the struct for a user
type User struct {
	ID       int      `json:"id"`       // User ID
	Username string   `json:"username"` // Username
	Password string   `json:"password"` // Password stored as sha256 hash
	Groups   []string `json:"groups"`   // Groups user belongs to
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
	UserID     int    `json:"user_id"`
}

// NewAuthDB creates a new authn database
func NewAuthDB(log *zap.Logger) *AuthDB {
	userBucketName := "users"
	tokenBucketName := "tokens"
	db := AuthDB{
		Logger:      log.With(zap.String("module", userBucketName)),
		userBucket:  userBucketName,
		tokenBucket: tokenBucketName,
	}
	err := db.Initialize()
	if err != nil {
		db.Logger.Error("failed to initialize authn database", zap.String("subsystem", userBucketName), zap.Error(err))
		return nil
	}
	return &db
}

// AuthDB is the struct for the authn database
type AuthDB struct {
	Logger      *zap.Logger
	DB          *bolt.DB
	userBucket  string
	tokenBucket string
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

	// Ensure that the user bucket exists
	err = db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(db.userBucket))
		if err != nil {
			return fmt.Errorf("failed to create userBucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create userBucket: %s", err)
	}

	// Ensure that the token bucket exists
	return db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(db.tokenBucket))
		if err != nil {
			return fmt.Errorf("failed to create tokenBucket: %s", err)
		}
		return nil
	})
}

// UpdateToken updates the token in the authn database
func (db *AuthDB) UpdateToken(key string, value string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.tokenBucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

// ViewToken views the token in the authn database
func (db *AuthDB) ViewToken(key string) (value string, err error) {
	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.tokenBucket))
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

// CreateUser creates a new user in the authn database
func (db *AuthDB) CreateUser(user *User) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.userBucket))

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		user.ID = int(id)

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Store the user in the users userBucket
		return b.Put(itob(user.ID), buf)
	})
}

func (db *AuthDB) CreateUserWithID(user *User) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.userBucket))

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Store the user in the users userBucket
		return b.Put(itob(user.ID), buf)
	})
}

// UpdateUser updates the authn database
func (db *AuthDB) UpdateUser(user *User) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.userBucket))

		id := itob(user.ID)
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put(id, buf)
	})
}

// DeleteUser deletes a user from the database
func (db *AuthDB) DeleteUser(userID int) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.userBucket))

		key := itob(userID)

		return b.Delete(key)
	})
}

// ViewUser views the specified user in the authn database
func (db *AuthDB) ViewUser(userID int) (user *User, err error) {
	var u User
	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.userBucket))

		key := itob(userID)

		v := b.Get(key)
		if v == nil {
			return fmt.Errorf("key %s not found in userBucket", key)
		}

		// Unmarshal the user
		return json.Unmarshal(v, &u)
	})
	return &u, err
}

// ViewUsers views the users in the authn database
func (db *AuthDB) ViewUsers() (users []*User, err error) {
	users = make([]*User, 0)
	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.userBucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u User
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}
			users = append(users, &u)
		}
		return nil
	})
	return users, err
}

// SafeViewUser views the specified user in the authn database without the password hashes
func (db *AuthDB) SafeViewUser(userID int) (user *User, err error) {
	u, err := db.ViewUser(userID)
	if err != nil {
		return nil, err
	}
	u.Password = "**********"
	return u, nil
}

// SafeViewUsers views the users in the authn database without the password hashes
func (db *AuthDB) SafeViewUsers() (users []*User, err error) {
	users, err = db.ViewUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		u.Password = "**********"
	}
	return users, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
