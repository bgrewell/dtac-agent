package middleware

import (
	"errors"
	"fmt"
	. "github.com/BGrewell/system-api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

var (
	lockout *LockoutStatus
	locktimeout int
)

const (
	LOCKOUT_UNLOCKED = "unlocked"
	LOCKOUT_LOCKED = "locked"
	LOCKOUT_PATH = "/lockout"
)

type LockoutStatus struct {
	Status string `json:"status"`
	Key string `json:"key,omitempty"`
	Host string `json:"host,omitempty"`
	Expiration *time.Time  `json:"expiration,omitempty"`
}

type LockoutError struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Lock LockoutStatus `json:"lock"`
}

func RegisterLockoutHandler(r *gin.Engine, lockoutTimeout int) {
	locktimeout = lockoutTimeout
	r.GET(LOCKOUT_PATH, GetLockoutHandler)
	r.POST(LOCKOUT_PATH, GetLockoutHandler)
	r.POST(LOCKOUT_PATH + "/:timeout", CreateLockoutHandler)
	r.PUT(LOCKOUT_PATH, RefreshLockoutHandler)
	r.PUT(LOCKOUT_PATH + "/:timeout", RefreshLockoutHandler)
	r.DELETE(LOCKOUT_PATH, DeleteLockoutHandler)
}

func GetLockoutHandler(c *gin.Context) {
	start := time.Now()
	WriteResponseJSON(c, time.Since(start), lockout)
}

func CreateLockoutHandler(c *gin.Context) {
	start := time.Now()

	if lockout.Status == LOCKOUT_LOCKED {
		WriteErrorResponseJSON(c, errors.New("unable to create lock as one already exists"))
	}

	timeoutStr := c.Param("timeout")
	timeout := locktimeout
	if timeoutStr != "" {
		t, err := strconv.Atoi(timeoutStr)
		if err != nil {
			WriteErrorResponseJSON(c, errors.New(fmt.Sprintf("failed to parse timeout value: %s", err)))
			return
		}
		timeout = t
	}

	expire := time.Now().Add(time.Duration(timeout) * time.Second)
	lockout = &LockoutStatus{
		Status:     LOCKOUT_LOCKED,
		Host: c.ClientIP(),
		Expiration: &expire,
		Key: uuid.New().String(),
	}
	WriteResponseJSON(c, time.Since(start), lockout)
}

func RefreshLockoutHandler(c *gin.Context) {
	start := time.Now()

	timeoutStr := c.Param("timeout")
	timeout := locktimeout
	if timeoutStr != "" {
		t, err := strconv.Atoi(timeoutStr)
		if err != nil {
			WriteErrorResponseJSON(c, errors.New(fmt.Sprintf("failed to parse timeout value: %s", err)))
			return
		}
		timeout = t
	}

	if lockout.Status == LOCKOUT_LOCKED && !(lockout.Expiration == nil || time.Now().After(*lockout.Expiration)) {
		*lockout.Expiration = time.Now().Add(time.Duration(timeout) * time.Second)
	} else {
		WriteErrorResponseJSON(c, errors.New("unable to refresh lock as no active lock exists"))
		return
	}
	WriteResponseJSON(c, time.Since(start), lockout)
}

func DeleteLockoutHandler(c *gin.Context) {
	start := time.Now()
	if lockout.Status == LOCKOUT_UNLOCKED || lockout.Expiration == nil || time.Now().After(*lockout.Expiration) {
		WriteErrorResponseJSON(c, errors.New("unable to delete lock as no active lock exists"))
	}
	lockout = &LockoutStatus{
		Status:     LOCKOUT_UNLOCKED,
		Expiration: nil,
	}
	WriteResponseJSON(c, time.Since(start), lockout)
}

func LockoutMiddleware() gin.HandlerFunc {
	t := time.Now()
	lockout = &LockoutStatus{
		Status:     LOCKOUT_UNLOCKED,
		Expiration: &t,
	}

	return func(c *gin.Context) {

		if lockout.Expiration != nil && time.Now().After(*lockout.Expiration) {
			lockout = &LockoutStatus{
				Status:     LOCKOUT_UNLOCKED,
				Expiration: nil,
			}
		}

		if lockout.Status == LOCKOUT_UNLOCKED || strings.HasPrefix(c.Request.URL.String(), LOCKOUT_PATH) {
			c.Next()
			return
		} else if lockout.Status == LOCKOUT_LOCKED && c.GetHeader("LOCKOUT_KEY") == lockout.Key {
			c.Next()
			return
		} else if c.GetHeader("LOCKOUT_KEY") != "" && c.GetHeader("LOCKOUT_KEY") != lockout.Key {
			le := LockoutError{
				Status:  "error",
				Message: fmt.Sprintf("your lockout key %s does not match the currently active key", c.GetHeader("LOCKOUT_KEY")),
				Lock:    *lockout,
			}
			c.AbortWithStatusJSON(401, le)
			return
		} else {
			le := LockoutError{
				Status: "error",
				Message: "system is currently locked out",
				Lock:   *lockout,
			}
			c.AbortWithStatusJSON(401, le)
			return
		}
	}
}