package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"my-flat-login/internal/model" // Replace with your actual module path

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct (must match the Claims struct in your service)
type Claims struct {
	FirebaseID string `json:"firebase_id"`
	jwt.RegisteredClaims
}

// JWT middleware to authenticate requests
func JWTAuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get the token from the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Parse and validate the JWT
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 3. Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract claims"})
			return
		}

		// 4. Get the user from the database (you'll need a UserRepository)
		// userRepository := // ... get your UserRepository instance
		// user, err := userRepository.FindByFirebaseID(ctx, claims.FirebaseID)
		// if err != nil {
		// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		// 	return
		// }

		// 5. Set the user in the context
		ctx.Set("user", &model.User{FirebaseID: claims.FirebaseID}) // Or set the actual user object

		ctx.Next()
	}
}
