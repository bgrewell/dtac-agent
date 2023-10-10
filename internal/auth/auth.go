package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/auth_db"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/register"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewAuthSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "auth"
	as := AuthSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
		admin: auth_db.User{
			ID:       1,
			Username: c.Config.Auth.User,
			Password: c.Config.Auth.Pass,
		},
	}
	return &as
}

type AuthSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	admin      auth_db.User
}

func (as *AuthSubsystem) Register() error {
	if !as.Enabled() {
		as.Logger.Info("subsystem is disabled", zap.String("subsystem", as.Name()))
		return nil
	}
	// Create a group for this subsystem
	base := as.Controller.Router.Group(as.name)

	// Routes
	routes := []types.RouteInfo{
		{Group: base, HttpMethod: http.MethodPost, Path: "/login", Handler: as.loginHandler, Protected: false},
	}

	// Register routes
	register.RegisterRoutes(routes, as.Controller.SecureMiddleware)
	as.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

func (as *AuthSubsystem) Enabled() bool {
	return as.enabled
}

func (as *AuthSubsystem) AuthHandler(c *gin.Context) {
	user, err := as.authorizeUser(c.Request)
	if err != nil {
		c.Header("DTAC-AUTHORIZATION", "DENIED")
		helpers.WriteUnauthorizedResponseJSON(c, err)
		return
	}

	as.Logger.Info("user granted access", zap.Uint64("userid", user.ID), zap.String("username", user.Username))
	c.Header("DTAC-AUTHORIZATION", "GRANTED")
	c.Next()
}

func (as *AuthSubsystem) Name() string {
	return as.name
}

func (as *AuthSubsystem) loginHandler(c *gin.Context) {
	var u auth_db.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	// check the admin user
	if as.admin.Username != u.Username || as.admin.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "invalid login credentials")
		return
	}

	token, err := as.createToken(as.admin.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := as.createAuth(as.admin.ID, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("CreateAuth: %s", saveErr.Error()))
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func (as *AuthSubsystem) extractToken(r *http.Request) string {
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

func (as *AuthSubsystem) extractTokenMetadata(token *jwt.Token) (*auth_db.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("unable to extract access id from token")
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &auth_db.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}

	return nil, errors.New("failed to get claims from token")
}

func (as *AuthSubsystem) createToken(userid uint64) (token *auth_db.TokenDetails, err error) {

	td := &auth_db.TokenDetails{
		AtExpires:   time.Now().Add(time.Minute * 15).Unix(),
		AccessUuid:  uuid.NewV4().String(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshUuid: uuid.NewV4().String(),
	}

	if os.Getenv("ACCESS_SECRET") == "" {
		err := os.Setenv("ACCESS_SECRET", "NEED_TO_GET_A_SECURE_SECRET_FROM_SOMEWHERE_IF_ENV_IS_EMPTY")
		as.Logger.Error("failed to set ACCESS_SECRET env variable", zap.Error(err))
	}
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

	if os.Getenv("REFRESH_SECRET") == "" {
		err := os.Setenv("REFRESH_SECRET", "NEED_TO_GET_A_REFRESH_SECRET_FROM_SOMEWHERE_IF_ENV_IS_EMPTY")
		as.Logger.Error("failed to set REFRESH_SECRET env variable", zap.Error(err))
	}
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

func (as *AuthSubsystem) createAuth(userid uint64, td *auth_db.TokenDetails) (err error) {
	errAccess := as.Controller.AuthDB.UpdateDB(td.AccessUuid, strconv.Itoa(int(userid))) //todo: need to look into how to time-expire these entries
	if errAccess != nil {
		return errAccess
	}

	errRefresh := as.Controller.AuthDB.UpdateDB(td.RefreshUuid, strconv.Itoa(int(userid)))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (as *AuthSubsystem) verifyToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(base64.URLEncoding.EncodeToString([]byte(os.Getenv("ACCESS_SECRET")))), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return token, nil
}

func (as *AuthSubsystem) fetchAuth(authD *auth_db.AccessDetails) (userId uint64, err error) {
	userIdStr, err := as.Controller.AuthDB.ViewDB(authD.AccessUuid)
	if userIdStr == "" || err != nil {
		return 0, fmt.Errorf("unable to find %s auth details in database", authD.AccessUuid)
	}
	userId, _ = strconv.ParseUint(userIdStr, 10, 64)
	return userId, nil
}

func (as *AuthSubsystem) authorizeUser(r *http.Request) (user *auth_db.User, err error) {
	tokenStr := as.extractToken(r)
	if tokenStr == "" {
		return nil, errors.New("invalid authorization header")
	}

	token, err := as.verifyToken(tokenStr)
	if err != nil {
		as.Logger.Error("failed to verify token", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	tokenAuth, err := as.extractTokenMetadata(token)
	if err != nil {
		as.Logger.Error("failed to get token metadata", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	userId, err := as.fetchAuth(tokenAuth)
	if err != nil {
		as.Logger.Error("failed to fetch auth", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	if userId == as.admin.ID {
		return &as.admin, nil
	} else {
		return nil, fmt.Errorf("unable to find %v", userId)
	}
}
