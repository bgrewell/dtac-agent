package handlers

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	Routes gin.RoutesInfo
	Info   BasicInfo
)

func init() {
	Info = BasicInfo{}
	Info.Update()
}

func SecretTestHandler(c *gin.Context) {
	user, err := AuthorizeUser(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
	}
	c.JSON(http.StatusOK, gin.H{"user": user.ID, "secret": "somesupersecretvalue"})
}

func HomeHandler(c *gin.Context) {
	// Update Routes
	start := time.Now()
	Info.UpdateRoutes(Routes)
	WriteResponseJSON(c, time.Since(start), Info)
}
