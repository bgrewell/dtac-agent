package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	bucket string = "auth"
)

var (
	DB *bolt.DB
	// Test User
	testUser = User{
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

type AccessDetails struct {
	AccessUuid string
	UserId uint64
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
	log.Printf("key: %s value: %s bucket %s\n", key, value, bucket)
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

	// check the testUser
	if testUser.Username != u.Username || testUser.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "invalid login credentials")
		return
	}

	token, err := CreateToken(testUser.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := CreateAuth(testUser.ID, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("CreateAuth: %s", saveErr.Error()))
		return
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
		RefreshUuid: uuid.NewV4().String(),
	}

	os.Setenv("ACCESS_SECRET", "FAKESECRETDONTUSEME") //todo: set this from an env file or externally
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessSecret := base64.URLEncoding.EncodeToString([]byte(os.Getenv("ACCESS_SECRET")))
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	os.Setenv("REFRESH_SECRET", "FAKEREFRESHSECRETDONOTUSEME") //todo: same as above
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshSecret := base64.URLEncoding.EncodeToString([]byte(os.Getenv("REFRESH_SECRET")))
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
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

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	if bearToken == "" {
		return ""
	}
	tokenArr := strings.Split(bearToken, " ")
	if len(tokenArr) == 2 {
		return tokenArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok :=token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(base64.URLEncoding.EncodeToString([]byte("FAKESECRETDONTUSEME"))), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok :=token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("token is invalid")
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId: userId,
		}, nil
	}

	return nil, err
}

func FetchAuth(authD *AccessDetails) (userId uint64, err error) {
	userIdStr := ViewDB(authD.AccessUuid)
	if userIdStr == "" {
		return 0, fmt.Errorf("unable to find testUser auth details in database")
	}
	userId, _ = strconv.ParseUint(userIdStr, 10, 64)
	return userId, nil
}

func AuthorizeUser(r *http.Request) (user *User, err error) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		log.Printf("error getting token metadata: %s", err.Error())
		return nil, fmt.Errorf("unauthorized")
	}
	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		return nil, err
	}
	//todo: lookup testUser when we have a users database
	if userId == testUser.ID {
		return &testUser, nil
	} else {
		return nil, fmt.Errorf("unable to find testUser")
	}
}