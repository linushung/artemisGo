package rest

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/linushung/artemis/cmd/server"
	"github.com/linushung/artemis/internal/app/authorization"
	"github.com/linushung/artemis/internal/pkg/configs"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
	defaultIdleTimeout  = 120 * time.Second
)

// Server represents a restful server
type Server struct {
	server.BaseServer
	Server http.Server
}

// InitRestServer run a HTTP server
func InitRestServer(base server.BaseServer) {
	/* Ref:
	1. https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	2. https://blog.cloudflare.com/exposing-go-on-the-internet/
	3. https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b
	*/
	srv := http.Server{
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}
	router := createRouter(&Server{base, srv})

	httpPort := configs.GetConfigStr("service.rest.port")
	log.Infof("***** [SERVER:REST] ***** Start a HTTP Server on port %s ......", httpPort)
	if err := http.ListenAndServe(httpPort, router); err != nil {
		log.Fatalf("***** [SERVER:REST][FAIL] ***** Failed to start HTTP Server: %v", err)
	}
}

func createRouter(s *Server) *gin.Engine {
	/* Ref: https://github.com/gin-gonic/gin */
	router := gin.Default()

	/* Health Check */
	router.GET("/ping", s.HTTPPing)
	/* pprof */
	pprof.Register(router, "/debug/pprof")

	/* Hystrix */
	hystrixGroup := router.Group("/hystrix")
	{
		hystrixGroup.POST("/ok", s.PostOK)
		hystrixGroup.POST("/status/{code}", s.PostStatus)
		hystrixGroup.POST("/delay/{second}", s.PostDelay)
	}

	/* Artemis */
	basicGroup := router.Group("/api/users")
	{
		basicGroup.POST("/", s.createUser)
		basicGroup.POST("/login", s.loginUser)
	}

	jwtAuth := router.Group("/api")
	jwtAuth.Use(authorization.VerifyJWTHandler(s.JWTMgr))
	{
		userGroup := jwtAuth.Group("/users")
		{
			userGroup.PUT("/", s.updateUser)
			userGroup.GET("/", s.fetchCurrentUser)
		}
		profileGroup := jwtAuth.Group("/profiles")
		{
			profileGroup.GET("/:username", s.fetchUserProfile)
			profileGroup.POST("/:username/follow", s.followUser)
			profileGroup.DELETE("/:username/follow", s.unFollowUser)
		}
		articleGroup := jwtAuth.Group("/articles")
		{
			articleGroup.POST("/", s.createArticle)
			articleGroup.GET("/feed", s.feedArticle)
			articleGroup.GET("/:slug", s.fetchArticle)
			articleGroup.PUT("/:slug", s.updateArticle)
			articleGroup.DELETE("/:slug", s.deleteArticle)
		}
	}

	return router
}

func createChiRouter(s *Server) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.Recoverer)

	/* Health Check */
	// router.Get("/ping", s.HTTPPing)
	/* pprof */
	router.Mount("/debug", middleware.Profiler())

	/* Hystrix */
	// router.Post("/ok", s.PostOK)
	// router.Post("/status/{code}", s.PostStatus)
	// router.Post("/delay/{second}", s.PostDelay)

	/* Artemis */
	router.Route("/api", func(r chi.Router) {
		r.Post("/users", s.createChiUser)
		r.Post("/users/login", s.loginChiUser)
		r.Put("/users", s.updateChiUser)
		r.Get("/users", s.fetchCurrentChiUser)
		r.Get("/profiles/{username}", s.fetchChiUserProfile)
		r.Post("/profiles/{username}/follow", s.followChiUser)
		r.Delete("/profiles/{username}/follow", s.unFollowChiUser)
	})

	return router
}
