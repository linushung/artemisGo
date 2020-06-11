package authorization

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	once       sync.Once
	instance   JWTMgr
	privateKey *rsa.PrivateKey
)

const (
	ValidPeriodMinute = 3
	JwtClaimsIssuer   = "artemis-MockIdentityManager"
	JwtClaimsAudience = "angular-realworld"
)

type JWTMgr interface {
	GenerateJWT(claims Claims) (string, error)
	VerifyJWT(token string) (Claims, error)
}

type Claims struct {
	Username string
	Role     string
	Jti      string
	Subject  string
}

type jwtClaims struct {
	Username string
	Role     string
	jwt.StandardClaims
}

type MockIdentityManager struct {
	Type string
}

func InitJWTService() {
	once.Do(func() {
		instance = MockIdentityManager{"artemisJWT"}

		// Use a single instance of Validate, it caches struct info
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Fatalf("***** [JWT][FAIL] ***** Failed to create RSA key pair:: %s", err)
			os.Exit(1)
		}

		privateKey = key
	})
}

func GetJWTMgr() JWTMgr {
	return instance
}

/* Ref: https://github.com/dgrijalva/jwt-go */
func (mgr MockIdentityManager) GenerateJWT(c Claims) (string, error) {
	issueTime := time.Now()

	claims := &jwtClaims{
		Username: c.Username,
		Role:     c.Role,
		StandardClaims: jwt.StandardClaims{
			Audience:  JwtClaimsAudience,
			ExpiresAt: issueTime.Add(ValidPeriodMinute * time.Minute).Unix(),
			Id:        c.Jti,
			IssuedAt:  issueTime.Unix(),
			Issuer:    JwtClaimsIssuer,
			NotBefore: issueTime.Unix(),
			Subject:   c.Subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	jtwStr, err := token.SignedString(privateKey)
	if err != nil {
		log.Errorf("***** [JWT][FAIL] ***** Failed to create JWT token:: %v", err)
		return "", err
	}

	return jtwStr, nil
}

func (mgr MockIdentityManager) VerifyJWT(token string) (Claims, error) {
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return &privateKey.PublicKey, nil
	})
	if err != nil {
		log.Errorf("***** [JWT][FAIL] ***** Failed to verify JWT:: %v", err)
		return Claims{}, err
	}

	return Claims{
		Username: claims.Username,
		Role: claims.Role,
		Jti: claims.Id,
		Subject: claims.Subject,
	}, nil
}

func VerifyJWTHandler(mgr JWTMgr) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Fields(ctx.Request.Header["Authorization"][0])[1]
		c, err := mgr.VerifyJWT(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "message": err.Error() })
			return
		}

		ctx.Set("token", c)
		ctx.Next()
	}
}

/* https://tools.ietf.org/html/rfc7519#section-4.1 */
/**
 * Reserved claims:
 * iss (issuer): Issuer of the JWT
 * sub (subject): Subject of the JWT (the user)
 * aud (audience): Recipient for which the JWT is intended
 * exp (expiration time): Time after which the JWT expires
 * nbf (not before time): Time before which the JWT must not be accepted for processing
 * iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
 * jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
 * */
