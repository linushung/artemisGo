package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/linushung/artemis/internal/app/database/postgres"
)

func (s *Server) createArticle(ctx *gin.Context) {
	req := &postgres.ArticleReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	art := postgres.Article{
		Slug:        req.Title,
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		Tags:        req.Tags,
	}

	id := uuid.New()
	a, err := s.RDB.CreateArticle(id, art)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if err := s.RDB.TagArticle(a.TagId, req.Tags); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	a.Tags = req.Tags
	ctx.JSON(http.StatusCreated, gin.H{"article": a})
}

func (s *Server) feedArticle(ctx *gin.Context) {

}

func (s *Server) fetchArticle(ctx *gin.Context) {

}

func (s *Server) updateArticle(ctx *gin.Context) {

}

func (s *Server) deleteArticle(ctx *gin.Context) {

}
