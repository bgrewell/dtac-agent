package authn

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/bgrewell/dtac-agent/internal/authndb"
	"github.com/bgrewell/dtac-agent/internal/controller"
	"github.com/bgrewell/dtac-agent/internal/helpers"
	"github.com/bgrewell/dtac-agent/internal/interfaces"
	"github.com/bgrewell/dtac-agent/internal/middleware"
	"github.com/bgrewell/dtac-agent/internal/types"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		enabled:    c.Config.Subsystems.Auth,
		name:       name,
	}

	adminID := -1

	if u, err := c.AuthDB.ViewUser(adminID); err != nil {
		as.Logger.Warn("failed to get admin user from database. creating admin user", zap.Error(err))
		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(c.Config.Auth.Pass), bcrypt.DefaultCost)
		if err != nil {
			as.Logger.Fatal("failed to update admin user", zap.Error(err))
		}
		user := &authndb.User{
			ID:       adminID,
			Username: c.Config.Auth.User,
			Password: string(hash),
			Groups:   []string{"admin"},
		}
		err = c.AuthDB.CreateUserWithID(user)
		if err != nil {
			as.Logger.Fatal("failed to create admin user", zap.Error(err))
		}
	} else {
		// Ensure user/pass hasn't changed in the config
		if u.Username != c.Config.Auth.User {
			as.Logger.Info("updating admin user", zap.String("username", c.Config.Auth.User))
			u.Username = c.Config.Auth.User
			err = c.AuthDB.UpdateUser(u)
			if err != nil {
				as.Logger.Fatal("failed to update admin user", zap.Error(err))
			}
		}
		if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Config.Auth.Pass)); err != nil {
			as.Logger.Info("updating admin user password")
			var hash []byte
			hash, err = bcrypt.GenerateFromPassword([]byte(c.Config.Auth.Pass), bcrypt.DefaultCost)
			if err != nil {
				as.Logger.Fatal("failed to update admin user", zap.Error(err))
			}
			u.Password = string(hash)
			err = c.AuthDB.UpdateUser(u)
			if err != nil {
				as.Logger.Fatal("failed to update admin user", zap.Error(err))
			}
		}
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
	authzGuest := endpoint.AuthGroupGuest.String()
	authzOperator := endpoint.AuthGroupOperator.String()
	authzAdmin := endpoint.AuthGroupAdmin.String()
	s.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/login", base), endpoint.ActionCreate, "login handler", s.loginHandler, false, authzGuest, endpoint.WithBody(authndb.UserArgs{}), endpoint.WithOutput(AuthOutput{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/users", base), endpoint.ActionRead, "list users", s.listUsers, true, authzOperator, endpoint.WithOutput([]authndb.User{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/user", base), endpoint.ActionRead, "get user by id", s.getUser, true, authzOperator, endpoint.WithOutput(authndb.User{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/users", base), endpoint.ActionCreate, "create user", s.createUser, true, authzAdmin, endpoint.WithBody(authndb.User{}), endpoint.WithOutput(authndb.User{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/user", base), endpoint.ActionWrite, "update user", s.updateUser, true, authzAdmin, endpoint.WithBody(authndb.User{}), endpoint.WithOutput(authndb.User{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/user", base), endpoint.ActionDelete, "delete user", s.deleteUser, true, authzAdmin),
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
	if !ep.Secure || !s.enabled {
		return ep.Function
	}
	return s.AuthenticationHandler(ep.Function)
}

func (s *Subsystem) loginHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		var inputUser = &authndb.User{ID: -9999}

		// Transform the body into a RouteTableRow
		if err := json.Unmarshal(in.Body, &inputUser); err != nil {
			return nil, nil, err
		}

		// Usernames are always worked with in lowercase
		inputUser.Username = strings.ToLower(inputUser.Username)

		// Convert the users credentials into sha256 hashes
		userHash := fmt.Sprintf("%x", sha256.Sum256([]byte(inputUser.Username)))

		// Operations performed here are done this way to ensure constant time comparison where authentication checks will
		// always take the same approximate amount of time to avoid timing attacks
		check := func(a, b string) bool {
			return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
		}

		// Get Users from database
		users, err := s.Controller.AuthDB.ViewUsers()
		if err != nil {
			return nil, nil, err
		}

		// Check if the user is in the database (always check all the users to avoid timing attacks)
		var userExists bool
		var matchUser *authndb.User
		for _, user := range users {
			uh := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ToLower(user.Username))))
			if check(userHash, uh) {
				userExists = true
				matchUser = user
			}
		}

		// Ensure we don't have a nil user. We don't want to simply skip the password check due to timing differences that
		// would create which could be used to determine if a user exists or not. Technically we may want to revisit all
		// of this code and simply ensure that the function takes N time to complete regardless of any internal operations.
		if matchUser == nil {
			matchUser = &authndb.User{ID: -9999}
			matchUser.Password = "" // Even if password matches we'll get a credential failure below due to user not existing
		}

		// check password
		passwordMatch := false
		if err := bcrypt.CompareHashAndPassword([]byte(matchUser.Password), []byte(inputUser.Password)); err == nil {
			passwordMatch = true
		}

		// check the users credentials
		if !(userExists) || !(passwordMatch) {
			return nil, nil, errors.New("invalid username or password")
		}

		token, err := s.createToken(matchUser.ID)
		if err != nil {
			return nil, nil, err
		}

		saveErr := s.createAuth(matchUser.ID, token)
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
		tokensJSON, err := json.Marshal(tokens)
		return headers, tokensJSON, err
	}, "authentication tokens")
}

func (s *Subsystem) createToken(userid int) (token *authndb.TokenDetails, err error) {

	// Parse token expiration times from config
	accessExpiration, err := s.parseTokenExpiration(s.Controller.Config.Auth.AccessTokenExpiration)
	if err != nil {
		s.Logger.Error("failed to parse access token expiration, using default 15m", zap.Error(err))
		accessExpiration = time.Minute * 15
	}

	refreshExpiration, err := s.parseTokenExpiration(s.Controller.Config.Auth.RefreshTokenExpiration)
	if err != nil {
		s.Logger.Error("failed to parse refresh token expiration, using default 168h", zap.Error(err))
		refreshExpiration = time.Hour * 24 * 7
	}

	// Calculate expiration times
	var atExpires int64
	var rtExpires int64
	
	if accessExpiration == 0 {
		// "never" case - set to a very far future date (100 years from now)
		atExpires = time.Now().Add(time.Hour * 24 * 365 * 100).Unix()
	} else {
		atExpires = time.Now().Add(accessExpiration).Unix()
	}
	
	if refreshExpiration == 0 {
		// "never" case - set to a very far future date (100 years from now)
		rtExpires = time.Now().Add(time.Hour * 24 * 365 * 100).Unix()
	} else {
		rtExpires = time.Now().Add(refreshExpiration).Unix()
	}

	td := &authndb.TokenDetails{
		AtExpires:   atExpires,
		AccessUUID:  uuid.NewV4().String(),
		RtExpires:   rtExpires,
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

// parseTokenExpiration parses a duration string and returns the corresponding time.Duration.
// It supports standard Go duration format (e.g., "15m", "1h", "24h") and a special "never" value.
// If "never" is specified, it returns 0 to indicate no expiration.
func (s *Subsystem) parseTokenExpiration(expirationStr string) (time.Duration, error) {
	// Handle the special "never" value
	if strings.ToLower(strings.TrimSpace(expirationStr)) == "never" {
		s.Logger.Info("token expiration set to 'never'")
		return 0, nil
	}

	// Parse standard Go duration format
	duration, err := time.ParseDuration(expirationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format '%s': %w", expirationStr, err)
	}

	// Validate that the duration is positive
	if duration < 0 {
		return 0, fmt.Errorf("duration must be positive, got: %s", expirationStr)
	}

	return duration, nil
}

func (s *Subsystem) createAuth(userid int, td *authndb.TokenDetails) (err error) {
	errAccess := s.Controller.AuthDB.UpdateToken(td.AccessUUID, strconv.Itoa(userid)) //todo: need to look into how to time-expire these entries
	if errAccess != nil {
		return errAccess
	}

	errRefresh := s.Controller.AuthDB.UpdateToken(td.RefreshUUID, strconv.Itoa(userid))
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
	return func(in *endpoint.Request) (out *endpoint.Response, err error) {
		// The AuthenticationHandler is a middleware function that is called before every secure request
		// that is used to get the user_id from the JWT token and store it in the request context to be
		// used by the Authorization handler
		s.Logger.Debug("authentication middleware called")
		var ok bool
		var auth string
		if auth, ok = in.Metadata[types.ContextAuthHeader.String()]; !ok {
			// Return error, API adapter should do a check to provide user with a more specific error
			return nil, errors.New("unable to authenticate user")
		}

		user, err := s.authorizeUser(auth)
		if err != nil {
			return nil, err
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		in.Metadata[types.ContextAuthUser.String()] = string(userJSON)
		return next(in)
	}
}

func (s *Subsystem) authorizeUser(bearerToken string) (user *authndb.User, err error) {
	tokenStr := s.extractToken(bearerToken)
	if tokenStr == "" {
		return nil, errors.New("invalid authorization header")
	}

	// Check if static testing token is configured and matches
	if s.Controller.Config.Auth.StaticTestingToken != "" && tokenStr == s.Controller.Config.Auth.StaticTestingToken {
		s.Logger.Info("static testing token used for authentication")
		// Return the admin user when static testing token is used
		adminID := -1
		user, err := s.Controller.AuthDB.ViewUser(adminID)
		if err != nil {
			s.Logger.Error("failed to get admin user for static testing token", zap.Error(err))
			return nil, errors.New("unable to authorize token")
		}
		return user, nil
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

	return s.Controller.AuthDB.ViewUser(userID)
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

func (s *Subsystem) fetchAuth(authD *authndb.AccessDetails) (userID int, err error) {
	userIDStr, err := s.Controller.AuthDB.ViewToken(authD.AccessUUID)
	if userIDStr == "" || err != nil {
		return 0, fmt.Errorf("unable to find %s authn details in database", authD.AccessUUID)
	}
	userID, err = strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %s to int", userIDStr)
	}
	return userID, nil
}

func (s *Subsystem) extractTokenMetadata(token *jwt.Token) (*authndb.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("unable to extract access id from token")
		}
		if userID, ok := claims["user_id"].(float64); ok {
			return &authndb.AccessDetails{
				AccessUUID: accessUUID,
				UserID:     int(userID),
			}, nil
		}

		return nil, errors.New("unable to extract user id from token")
	}

	return nil, errors.New("failed to get claims from token")
}

func (s *Subsystem) listUsers(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		users, err := s.Controller.AuthDB.SafeViewUsers()
		if err != nil {
			return nil, err
		}
		return json.Marshal(users)
	}, "users configured for access to the system")
}

func (s *Subsystem) getUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		if m, ok := in.Parameters["id"]; ok && len(m) > 0 && m[0] != "" {
			id := m[0]
			uid, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			user, err := s.Controller.AuthDB.SafeViewUser(uid)
			if err != nil {
				return nil, err
			}
			return json.Marshal(user)
		}
		return nil, errors.New("missing parameter 'id'")

	}, "information for the specified user")
}

func (s *Subsystem) createUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		var user authndb.User
		if in.Body == nil || len(in.Body) == 0 {
			return nil, errors.New("missing body")
		}

		// Transform the body into a user
		if err := json.Unmarshal(in.Body, &user); err != nil {
			return nil, err
		}

		// Check to see if the user already exists (this is not secure against timing attacks because you need to be an admin already to do it)
		if s.Controller.AuthDB.UserExistsByUsername(user.Username) {
			return nil, errors.New("user already exists")
		}

		// Hash the password
		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hash)

		// Add to database
		err = s.Controller.AuthDB.CreateUser(&user)
		if err != nil {
			return nil, err
		}

		// Return safe view of updated user
		safeUser, err := s.Controller.AuthDB.SafeViewUser(user.ID)
		if err != nil {
			return nil, err
		}

		return json.Marshal(safeUser)
	}, "information for the user that has been created")
}

func (s *Subsystem) updateUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		var user authndb.User

		if m, ok := in.Parameters["id"]; ok && len(m) > 0 && m[0] != "" {
			if in.Body == nil || len(in.Body) == 0 {
				return nil, errors.New("missing body")
			}

			// Transform the body into a user
			if err := json.Unmarshal(in.Body, &user); err != nil {
				return nil, err
			}

			// Get DB user with id
			id := m[0]
			uid, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			dbUser, err := s.Controller.AuthDB.ViewUser(uid)
			if err != nil {
				return nil, err
			}

			// Update ID if not specified
			if user.ID == 0 {
				user.ID = dbUser.ID
			}

			// Verify
			if dbUser.ID != user.ID {
				return nil, errors.New("user id mismatch, you cannot change the user id")
			}
			if dbUser.Username != user.Username {
				return nil, errors.New("user username mismatch, you cannot change the username")
			}

			// If password changed then rehash it
			if (user.Password != "") && (fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))) != dbUser.Password) {
				s.Logger.Info("password changed, rehashing", zap.String("username", user.Username), zap.Int("id", user.ID))
				var hash []byte
				hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {
					return nil, err
				}
				user.Password = string(hash)
			}

			// Add to database
			err = s.Controller.AuthDB.UpdateUser(&user)
			if err != nil {
				return nil, err
			}

			// Return safe view of updated user
			safeUser, err := s.Controller.AuthDB.SafeViewUser(user.ID)
			if err != nil {
				return nil, err
			}

			return json.Marshal(safeUser)
		}
		return nil, errors.New("missing parameter 'id'")

	}, "information for the user that has been updated")
}

func (s *Subsystem) deleteUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		if m, ok := in.Parameters["id"]; ok && len(m) > 0 && m[0] != "" {
			id := m[0]
			uid, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			err = s.Controller.AuthDB.DeleteUser(uid)
			if err != nil {
				return nil, err
			}
			return json.Marshal(map[string]int{"deleted_uid": uid})
		}
		return nil, errors.New("missing parameter 'id'")

	}, "no information is returned by this endpoint")
}
