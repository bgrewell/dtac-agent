package authn

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authndb"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/middleware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

// AuthOutput is a struct to assist with describing the output format
type AuthOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewSubsystem creates a new authn subsystem
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "auth"
	as := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
		admin: authndb.User{ // This is all stubbed in until we get the authentication database up and running
			ID:             1,
			Username:       c.Config.Auth.User,
			UsernameHashed: fmt.Sprintf("%x", sha256.Sum256([]byte(c.Config.Auth.User))),
			Password:       fmt.Sprintf("%x", sha256.Sum256([]byte(c.Config.Auth.Pass))), //TODO: Store as a sha256 hash in the configuration file
			Groups:         []string{"admin"},
		},
		guest: authndb.User{
			ID:             2,
			Username:       "guest",
			UsernameHashed: fmt.Sprintf("%x", sha256.Sum256([]byte("guest"))),
			Password:       fmt.Sprintf("%x", sha256.Sum256([]byte("guest"))),
			Groups:         []string{"guest"},
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
	admin      authndb.User //TODO: Replace with actual authentication database
	guest      authndb.User
	endpoints  []*endpoint.Endpoint
}

// register registers the authn subsystem
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	// Create a group for this subsystem
	base := s.name

	// Endpoints
	authz := endpoint.AuthGroupGuest.String()
	s.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/login", base), endpoint.ActionCreate, s.loginHandler, false, authz, endpoint.WithBody(authndb.UserArgs{}), endpoint.WithOutput(AuthOutput{})),
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
func (s *Subsystem) Endpoints() []*endpoint.Endpoint {
	return s.endpoints
}

// Handler handles the authentication middleware
func (s *Subsystem) Handler(ep endpoint.Endpoint) endpoint.Func {
	// Bypass authentication for endpoints that don't use auth
	if !ep.Secure {
		return ep.Function
	}
	return s.AuthenticationHandler(ep.Function)
}

// TODO: Need to make sure this function can access the context used for logging in
func (s *Subsystem) loginHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		var u authndb.User

		// Transform the body into a RouteTableRow
		if err := json.Unmarshal(in.Body, &u); err != nil {
			return nil, nil, err
		}

		// Convert the users credentials into sha256 hashes
		userHash := fmt.Sprintf("%x", sha256.Sum256([]byte(u.Username)))
		passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(u.Password)))

		// Operations performed here are done this way to ensure constant time comparison where authentication checks will
		// always take the same approximate amount of time to avoid timing attacks
		check := func(a, b string) bool {
			return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
		}
		isAdminUsername := check(userHash, s.admin.UsernameHashed)
		isAdminPassword := check(passHash, s.admin.Password)
		isGuestUsername := check(userHash, s.guest.UsernameHashed)
		isGuestPassword := check(passHash, s.guest.Password)

		// check the users credentials //TODO: Replace with actual authentication database
		if !((isAdminUsername && isAdminPassword) || (isGuestUsername && isGuestPassword)) {
			return nil, nil, errors.New("invalid username or password")
		}

		// TODO: Fake lookup in database
		if u.Username == "admin" {
			u = s.admin // set the user to admins
		} else {
			u = s.guest // guest is the only other valid user at this point
		}

		token, err := s.createToken(u.ID)
		if err != nil {
			return nil, nil, err
		}

		saveErr := s.createAuth(u.ID, token)
		if saveErr != nil {
			return nil, nil, err
		}

		tokens := AuthOutput{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}

		// TODO: Should also package and transfer the refresh_token as a cookie here? (probably better to handle in the REST API)
		headers := map[string][]string{
			"Authorization": {fmt.Sprintf("Bearer %s", token.AccessToken)},
		}
		tokensJson, err := json.Marshal(tokens)
		return headers, tokensJson, nil
	}, "authentication tokens")
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

// START OF AuthenticationMiddleware portion of code

// Priority returns the priority of the middleware
func (s *Subsystem) Priority() middleware.Priority {
	return middleware.PriorityAuthentication
}

// AuthenticationHandler is the middleware function that is called before every secure request
func (s *Subsystem) AuthenticationHandler(next endpoint.Func) endpoint.Func {
	return func(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
		// The AuthenticationHandler is a middleware function that is called before every secure request
		// that is used to get the user_id from the JWT token and store it in the request context to be
		// used by the Authorization handler
		s.Logger.Debug("authentication middleware called")
		if auth, ok := in.Metadata[types.ContextAuthHeader.String()]; !ok {
			// Return error, API adapter should do a check to provide user with a more specific error
			return nil, errors.New("unable to authenticate user")
		} else {
			user, err := s.authorizeUser(auth)
			if err != nil {
				return nil, err
			}

			userJson, err := json.Marshal(user)
			if err != nil {
				return nil, err
			}
			in.Metadata[types.ContextAuthUser.String()] = string(userJson)
			return next(in)
		}
	}
}

func (s *Subsystem) authorizeUser(bearerToken string) (user *authndb.User, err error) {
	tokenStr := s.extractToken(bearerToken)
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

func (s *Subsystem) extractToken(bearerToken string) string {
	if bearerToken == "" {
		return ""
	}
	tokenArr := strings.Split(bearerToken, " ")
	if len(tokenArr) == 2 {
		return tokenArr[1]
	}
	return ""
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
