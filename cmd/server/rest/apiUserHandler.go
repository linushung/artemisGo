package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/linushung/artemis/internal/app/authorization"
	"github.com/linushung/artemis/internal/app/database/postgres"
)

func statusCode(err error) int {
	switch err.Error() {
	case "sql: no rows in result set", "row(s) affected: 0":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	req := &postgres.RegisterReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	p := &postgres.Poster{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hash),
		Role:     string(postgres.User),
	}

	if err := s.RDB.CreatePoster(*p); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": p})
}

func (s *Server) loginUser(ctx *gin.Context) {
	req := &postgres.LoginReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	p, err := s.RDB.SelectPosterByEmail(req.Email)
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c := authorization.Claims{
		Username: p.Username,
		Role:     p.Role,
		Jti:      p.Email,
		Subject:  p.Email,
	}
	token, err := s.JWTMgr.GenerateJWT(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) updateUser(ctx *gin.Context) {
	req := &postgres.UpdateReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c := ctx.MustGet("token").(authorization.Claims)
	p, err := s.RDB.UpdatePoster(c.Subject, req)
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": p})
}

func (s *Server) fetchCurrentUser(ctx *gin.Context) {
	c := ctx.MustGet("token").(authorization.Claims)
	p, err := s.RDB.SelectPosterByEmail(c.Subject)
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": p})
}

func (s *Server) fetchUserProfile(ctx *gin.Context) {
	p, err := s.RDB.SelectPosterByUsername(ctx.Param("username"))
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	followers, err := s.RDB.FetchFollowersByEmail(p.Email)
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	isFollowing := false
	c := ctx.MustGet("token").(authorization.Claims)
	for _, f := range followers {
		if c.Username == f {
			isFollowing = true
			break
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"profile": &postgres.Profile{
		Username:  p.Username,
		Image:     p.Image,
		Bio:       p.Bio,
		Following: isFollowing,
	}})
}

func (s *Server) followUser(ctx *gin.Context) {
	p, err := s.RDB.SelectPosterByUsername(ctx.Param("username"))
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	c := ctx.MustGet("token").(authorization.Claims)
	if err := s.RDB.FollowPoster(p.Email, c.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"profile": &postgres.Profile{
		Username:  p.Username,
		Image:     p.Image,
		Bio:       p.Bio,
		Following: true,
	}})
}

func (s *Server) unFollowUser(ctx *gin.Context) {
	p, err := s.RDB.SelectPosterByUsername(ctx.Param("username"))
	if err != nil {
		ctx.JSON(statusCode(err), gin.H{"message": err.Error()})
		return
	}

	c := ctx.MustGet("token").(authorization.Claims)
	if err := s.RDB.UnFollowPoster(p.Email, c.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"profile": &postgres.Profile{
		Username:  p.Username,
		Image:     p.Image,
		Bio:       p.Bio,
		Following: false,
	}})
}
