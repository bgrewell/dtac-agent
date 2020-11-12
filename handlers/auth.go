package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	bucket string = "auth"
)

var (
	DB *bolt.DB
	// Test User
	user = User{
		ID:       1,
		Username: "intel",
		Password: "intel123",
	}
)

type User struct {
	ID uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AccessUuid string
	RefreshUuid string
	AtExpires int64
	RtExpires int64
}

func init() {
	// Initialize Database
	var err error
	DB, err = bolt.Open("db/auth.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	//todo: need to have a finalizer for the whole program. this should close DB
	DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %s", err)
		}
		return nil
	})
}

func UpdateDB(key string, value string) (err error) {
	err = DB.Update(func(tx *bolt.Tx) error {
		b :=tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return err
}

func ViewDB(key string) (value string) {
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		value = string(v)
		return nil
	})
	return value
}

func LoginHandler(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	// check the user
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "invalid login credentials")
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := CreateAuth(user.ID, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token": token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func CreateToken(userid uint64) (token *TokenDetails, err error) {

	td := &TokenDetails{
		AtExpires: time.Now().Add(time.Minute * 15).Unix(),
		AccessUuid: uuid.NewV4().String(),
		RtExpires: time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshToken: uuid.NewV4().String(),
	}

	os.Setenv("ACCESS_SECRET", "FAKESECRETDONTUSEME") //todo: set this from an env file or externally
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	os.Setenv("REFRESH_SECRET", "FAKEREFRESHSECRETDONOTUSEME") //todo: same as above
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userid uint64, td *TokenDetails) (err error) {
	//at := time.Unix(td.AtExpires, 0)
	//rt := time.Unix(td.RtExpires, 0)
	//now := time.Now()

	errAccess := UpdateDB(td.AccessUuid, strconv.Itoa(int(userid))) //todo: need to look into how to time-expire these entries
	if errAccess != nil {
		return errAccess
	}

	errRefresh := UpdateDB(td.RefreshUuid, strconv.Itoa(int(userid)))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}