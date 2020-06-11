package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPPing generates PONG response to a Ping request for health checking
func (s *BaseServer) HTTPPing(c *gin.Context) {
	c.String(http.StatusOK, "PONG")
}
