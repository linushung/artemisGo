package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/linushung/artemis/internal/app/authorization"
	"github.com/linushung/artemis/internal/app/database/postgres"
)

func (s *Server) createChiUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req := &postgres.RegisterReq{}
	if err := validateReq(r.Body, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := &postgres.Poster{
		Email: req.Email,
		Username: req.Username,
		Password: string(hash),
		Role: string(postgres.User),
	}

	if err := s.RDB.CreatePoster(*p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (s *Server) loginChiUser(w http.ResponseWriter, r *http.Request) {
	req := &postgres.LoginReq{}
	if err := validateReq(r.Body, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := s.RDB.SelectPosterByEmail(req.Email)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sql: no rows in result set" {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	c := authorization.Claims{
		Username: p.Username,
		Role: p.Role,
		Jti: p.Email,
		Subject: p.Email,
	}
	token, err := s.JWTMgr.GenerateJWT(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jwt := struct { Token string } {token }
	json.NewEncoder(w).Encode(jwt)
}

func (s *Server) updateChiUser(w http.ResponseWriter, r *http.Request) {
	token := strings.Fields(r.Header.Get("Authorization"))[1]
	c, err := s.JWTMgr.VerifyJWT(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	req := &postgres.UpdateReq{}
	if err := validateReq(r.Body, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := s.RDB.UpdatePoster(c.Subject, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) fetchCurrentChiUser(w http.ResponseWriter, r *http.Request) {
	token := strings.Fields(r.Header.Get("Authorization"))[1]
	c, err := s.JWTMgr.VerifyJWT(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	p, err := s.RDB.SelectPosterByEmail(c.Subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) fetchChiUserProfile(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) followChiUser(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) unFollowChiUser(w http.ResponseWriter, r *http.Request) {

}
