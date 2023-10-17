package authn

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authn_db"
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

func NewAuthnSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "auth"
	as := AuthnSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
		admin: authn_db.User{ // This is all stubbed in until we get the authentication database up and running
			ID:       1,
			Username: c.Config.Auth.User,
			Password: c.Config.Auth.Pass,
			Groups:   []string{"admin"},
		},
		guest: authn_db.User{
			ID:       2,
			Username: "guest",
			Password: "guest",
			Groups:   []string{"guest"},
		},
	}
	return &as
}

type AuthnSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	admin      authn_db.User
	guest      authn_db.User
}

func (as *AuthnSubsystem) Register() error {
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

func (as *AuthnSubsystem) Enabled() bool {
	return as.enabled
}

func (as *AuthnSubsystem) AuthenticationHandler(c *gin.Context) {
	// The AuthenticationHandler is a middleware function that is called before every secure request
	// that is used to get the user_id from the JWT token and store it in the request context to be
	// used by the Authorization handler
	user, err := as.authorizeUser(c.Request)
	if err != nil {
		c.Header("X-DTAC-AUTHENTICATION", "INCOMPLETE")
		helpers.WriteUnauthorizedResponseJSON(c, err)
		c.Abort()
		return
	}
	//as.Logger.Info("user granted access",
	//	zap.Uint64("userid", user.ID),
	//	zap.String("username", user.Username),
	//	zap.Any("groups", user.Groups))
	c.Header("X-DTAC-AUTHENTICATION", user.Username)

	c.Set("user_id", user.ID)
	c.Set("username", user.Username)
	c.Set("groups", user.Groups)
	c.Next()
}

func (as *AuthnSubsystem) Name() string {
	return as.name
}

func (as *AuthnSubsystem) loginHandler(c *gin.Context) {
	var u authn_db.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	// check the admin user
	if (as.admin.Username != u.Username || as.admin.Password != u.Password) &&
		(as.guest.Username != u.Username || as.guest.Password != u.Password) {
		c.JSON(http.StatusUnauthorized, "invalid login credentials")
		return
	}

	// TODO: Fake lookup in database
	if u.Username == "admin" {
		u = as.admin // set the user to admins
	} else {
		u = as.guest // guest is the only other valid user at this point
	}

	token, err := as.createToken(u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := as.createAuth(u.ID, token)
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

func (as *AuthnSubsystem) extractToken(r *http.Request) string {
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

func (as *AuthnSubsystem) extractTokenMetadata(token *jwt.Token) (*authn_db.AccessDetails, error) {
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
		return &authn_db.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}

	return nil, errors.New("failed to get claims from token")
}

func (as *AuthnSubsystem) createToken(userid uint64) (token *authn_db.TokenDetails, err error) {

	td := &authn_db.TokenDetails{
		AtExpires:   time.Now().Add(time.Minute * 15).Unix(),
		AccessUuid:  uuid.NewV4().String(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshUuid: uuid.NewV4().String(),
	}

	if os.Getenv("ACCESS_SECRET") == "" {
		err := os.Setenv("ACCESS_SECRET", "NEED_TO_GET_A_SECURE_SECRET_FROM_SOMEWHERE_IF_ENV_IS_EMPTY")
		if err != nil {
			as.Logger.Error("failed to set ACCESS_SECRET env variable", zap.Error(err))
		}
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
		if err != nil {
			as.Logger.Error("failed to set REFRESH_SECRET env variable", zap.Error(err))
		}
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

func (as *AuthnSubsystem) createAuth(userid uint64, td *authn_db.TokenDetails) (err error) {
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

func (as *AuthnSubsystem) verifyToken(tokenStr string) (*jwt.Token, error) {
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

func (as *AuthnSubsystem) fetchAuth(authD *authn_db.AccessDetails) (userId uint64, err error) {
	userIdStr, err := as.Controller.AuthDB.ViewDB(authD.AccessUuid)
	if userIdStr == "" || err != nil {
		return 0, fmt.Errorf("unable to find %s authn details in database", authD.AccessUuid)
	}
	userId, _ = strconv.ParseUint(userIdStr, 10, 64)
	return userId, nil
}

func (as *AuthnSubsystem) authorizeUser(r *http.Request) (user *authn_db.User, err error) {
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
		as.Logger.Error("failed to fetch authn", zap.Error(err))
		return nil, errors.New("unable to authorize token")
	}

	if userId == as.admin.ID {
		return &as.admin, nil
	} else if userId == as.guest.ID {
		return &as.guest, nil
	} else {
		return nil, fmt.Errorf("unable to find %v", userId)
	}
}
