package server

import (
	"github.com/linushung/artemis/internal/app/authorization"
	"github.com/linushung/artemis/internal/app/database/postgres"
)

// BaseServer represents a generic server
type BaseServer struct {
	RDB    postgres.RDB
	CircuitBreakerManager
	authorization.JWTMgr
}

// NewBaseServer return an instance of BaseServer struct.
func NewBaseServer() *BaseServer {
	return &BaseServer{
		postgres.InitPostgreSQL(),
		GetCircuitBreakerMgr(),
		authorization.GetJWTMgr(),
	}
}
