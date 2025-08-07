package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
)

const (
	keycloakIssuer = "http://localhost:8080/realms/master"                               // Ваш issuer
	jwksURL        = "http://localhost:8080/realms/master/protocol/openid-connect/certs" // JWKS endpoint
	audience       = "golang-api"                                                        // Audience из Keycloak
)

var keySet jwk.Set

func FetchJWKS() error {
	ctx := context.Background()
	set, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	keySet = set
	return nil
}

func StartJWKSRefresh() {
	if err := FetchJWKS(); err != nil {
		fmt.Printf("Initial JWKS fetch failed: %v\n", err)
	}

	go func() {
		for range time.Tick(1 * time.Hour) {
			if err := FetchJWKS(); err != nil {
				fmt.Printf("JWKS refresh failed: %v\n", err)
			}
		}
	}()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" || strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("kid not found")
			}

			key, found := keySet.LookupKeyID(kid)
			if !found {
				return nil, fmt.Errorf("key not found")
			}

			var rawKey interface{}
			if err := key.Raw(&rawKey); err != nil {
				return nil, fmt.Errorf("failed to get raw key: %w", err)
			}
			return rawKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		iss, ok := claims["iss"].(string)
		if !ok || iss != keycloakIssuer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid issuer"})
			return
		}

		aud, ok := claims["aud"].([]interface{})
		if !ok {
			audStr, ok := claims["aud"].(string)
			if ok {
				aud = []interface{}{audStr}
			}
		}
		audienceFound := false
		for _, a := range aud {
			if a == audience {
				audienceFound = true
				break
			}
		}
		if !audienceFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
			return
		}

		if scopes, ok := claims["scope"].(string); ok {
			if !strings.Contains(scopes, "api:read") {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient scopes"})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No scopes in token"})
			return
		}

		if username, ok := claims["preferred_username"].(string); ok {
			c.Set("user", username)
		}

		c.Next()
	}
}
