package authn

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authndb"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// NewSubsystem creates a new authn subsystem
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "auth"
	as := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
		admin: authndb.User{ // This is all stubbed in until we get the authentication database up and running
			ID:       1,
			Username: c.Config.Auth.User,
			Password: c.Config.Auth.Pass,
			Groups:   []string{"admin"},
		},
		guest: authndb.User{
			ID:       2,
			Username: "guest",
			Password: "guest",
			Groups:   []string{"guest"},
		},
	}
	as.register()
	return &as
}

// Subsystem is the subsystem for authentication
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	admin      authndb.User
	guest      authndb.User
	endpoints  []endpoint.Endpoint
}

// register registers the authn subsystem
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	// Create a group for this subsystem
	base := s.Controller.Router.Group(s.name)

	// Endpoints
	secure := s.Controller.Config.Auth.DefaultSecure
	s.endpoints = []endpoint.Endpoint{
		{fmt.Sprintf("%s/login", base), endpoint.ActionRead, s.loginHandler, secure},
	}
}

// Enabled returns true if the subsystem is enabled
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the subsystem
func (s *Subsystem) Name() string {
	return s.name
}

// Endpoints returns an array of endpoints that this Subsystem handles
func (s *Subsystem) Endpoints() []endpoint.Endpoint {
	return s.endpoints
}

// TODO: Will need to figure out how this fits into the new decoupled API frontend architecture
// AuthenticationHandler is the middleware function that is called before every secure request
func (s *Subsystem) AuthenticationHandler(c *gin.Context) {
	// The AuthenticationHandler is a middleware function that is called before every secure request
	// that is used to get the user_id from the JWT token and store it in the request context to be
	// used by the Authorization handler
	user, err := s.authorizeUser(c.Request)
	if err != nil {
		c.Header("X-DTAC-AUTHENTICATION", "INCOMPLETE")
		s.Controller.Formatter.WriteUnauthorizedError(c, err)
		c.Abort()
		return
	}
	//s.Logger.Info("user granted access",
	//	zap.Uint64("userid", user.ID),
	//	zap.String("username", user.Username),
	//	zap.Any("groups", user.Groups))
	c.Header("X-DTAC-AUTHENTICATION", user.Username)

	c.Set("user_id", user.ID)
	c.Set("username", user.Username)
	c.Set("groups", user.Groups)
	c.Next()
}

// TODO: Need to make sure this function can access the context used for logging in
func (s *Subsystem) loginHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	start := time.Now()
	var u authndb.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	// check the admin user
	if (s.admin.Username != u.Username || s.admin.Password != u.Password) &&
		(s.guest.Username != u.Username || s.guest.Password != u.Password) {
		c.JSON(http.StatusUnauthorized, "invalid login credentials")
		return
	}

	// TODO: Fake lookup in database
	if u.Username == "admin" {
		u = s.admin // set the user to admins
	} else {
		u = s.guest // guest is the only other valid user at this point
	}

	token, err := s.createToken(u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := s.createAuth(u.ID, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("CreateAuth: %s", saveErr.Error()))
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	s.Controller.Formatter.WriteResponse(c, time.Since(start), tokens)
}

func (s *Subsystem) extractToken(r *http.Request) string {
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

func (s *Subsystem) extractTokenMetadata(token *jwt.Token) (*authndb.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("unable to extract access id from token")
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &authndb.AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}

	return nil, errors.New("failed to get claims from token")
}

func (s *Subsystem) createToken(userid uint64) (token *authndb.TokenDetails, err error) {

	td := &authndb.TokenDetails{
		AtExpires:   time.Now().Add(time.Minute * 15).Unix(),
		AccessUUID:  uuid.NewV4().String(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshUUID: uuid.NewV4().String(),
	}

	if os.Getenv("ACCESS_SECRET") == "" {
		err := os.Setenv("ACCESS_SECRET", "NEED_TO_GET_A_SECURE_SECRET_FROM_SOMEWHERE_IF_ENV_IS_EMPTY")
		if err != nil {
			s.Logger.Error("failed to set ACCESS_SECRET env variable", zap.Error(err))
		}
	}
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessSecret := base64.URLEncoding.EncodeToString([]byte(os.Getenv("ACCESS_SECRET")))
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	if os.Getenv("REFRESH_SECRET") == "" {
		err := os.Setenv("REFRESH_SECRET", "NEED_TO_GET_A_REFRESH_SECRET_FROM_SOMEWHERE_IF_ENV_IS_EMPTY")
		if err != nil {
			s.Logger.Error("failed to set REFRESH_SECRET env variable", zap.Error(err))
		}
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
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

func (s *Subsystem) createAuth(userid uint64, td *authndb.TokenDetails) (err error) {
	errAccess := s.Controller.AuthDB.UpdateDB(td.AccessUUID, strconv.Itoa(int(userid))) //todo: need to look into how to time-expire these entries
	if errAccess != nil {
		return errAccess
	}

	errRefresh := s.Controller.AuthDB.UpdateDB(td.RefreshUUID, strconv.Itoa(int(userid)))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (s *Subsystem) verifyToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(base64.URLEncoding.EncodeToString([]byte(os.Getenv("ACCESS_SECRET")))), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return token, nil
}

func (s *Subsystem) fetchAuth(authD *authndb.AccessDetails) (userID uint64, err error) {
	userIDStr, err := s.Controller.AuthDB.ViewDB(authD.AccessUUID)
	if userIDStr == "" || err != nil {
		return 0, fmt.Errorf("unable to find %s authn details in database", authD.AccessUUID)
	}
	userID, _ = strconv.ParseUint(userIDStr, 10, 64)
	return userID, nil
}

func (s *Subsystem) authorizeUser(r *http.Request) (user *authndb.User, err error) {
	tokenStr := s.extractToken(r)
	if tokenStr == "" {
		return nil, errors.New("invalid authorization header")
	}

	token, err := s.verifyToken(tokenStr)
	if err != nil {
		s.Logger.Error("failed to verify token", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	tokenAuth, err := s.extractTokenMetadata(token)
	if err != nil {
		s.Logger.Error("failed to get token metadata", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	userID, err := s.fetchAuth(tokenAuth)
	if err != nil {
		s.Logger.Error("failed to fetch authn", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	if userID == s.admin.ID {
		return &s.admin, nil
	} else if userID == s.guest.ID {
		return &s.guest, nil
	} else {
		return nil, fmt.Errorf("unable to find %v", userID)
	}
}
