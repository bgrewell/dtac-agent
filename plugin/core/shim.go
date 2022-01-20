package core

import "github.com/gin-gonic/gin"

type Shim struct {
	Name string
}

func (s *Shim) CustomFileHandler(c *gin.Context) {

}
