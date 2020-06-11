package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	register = "httpbinservice"
)

// PostOK mock successful Post request to httpbin
func (s *Server) PostOK(c *gin.Context) {
	url := fmt.Sprintf("http://localhost:8000/post")
	_, httpErr := s.CBHTTPPost(register, url, "", []byte(""))
	if httpErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": httpErr.Error(),
		})
	} else {
		c.String(http.StatusOK, "OK")
	}
}

// PostStatus mock response of specific status code from httpbin
func (s *Server) PostStatus(c *gin.Context) {
	code, _ := strconv.Atoi(c.Param("code"))
	url := fmt.Sprintf("http://localhost:8000/status/%d", code)
	_, httpErr := s.CBHTTPPost(register, url, "", []byte(""))
	if httpErr != nil {
		c.JSON(code, gin.H{
			"message": httpErr.Error(),
		})
	} else {
		c.String(http.StatusOK, "OK")
	}
}

// PostDelay mock delay response of specific seconds from httpbin
func (s *Server) PostDelay(c *gin.Context) {
	url := fmt.Sprintf("http://localhost:8000/delay/%s", c.Param("second"))
	_, httpErr := s.CBHTTPPost(register, url, "", []byte(""))
	if httpErr != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"message": httpErr.Error(),
		})
	} else {
		c.String(http.StatusOK, "OK")
	}
}
