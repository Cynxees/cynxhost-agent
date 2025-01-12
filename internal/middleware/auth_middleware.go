package middleware

import (
	"context"
	"cynxhostagent/internal/constant/types"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/helper"
	contextmodel "cynxhostagent/internal/model/context"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(JWTManager *dependencies.JWTManager, next http.HandlerFunc, debug bool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Check for the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			apiResponse := response.APIResponse{
				Code:  responsecode.CodeAuthenticationError,
				Error: "Authorization header missing",
			}
			helper.WriteJSONResponse(w, http.StatusUnauthorized, apiResponse)
			return
		}

		// Check if the token starts with "Bearer"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			apiResponse := response.APIResponse{
				Code:  responsecode.CodeAuthenticationError,
				Error: "Invalid authorization token format",
			}
			helper.WriteJSONResponse(w, http.StatusUnauthorized, apiResponse)
			return
		}

		// You could verify the token here if needed (e.g., check JWT signature)
		token := parts[1]
		claims, err := JWTManager.VerifyToken(token)

		if err != nil { // Replace with your token verification logic
			apiResponse := response.APIResponse{
				Code:  responsecode.CodeAuthenticationError,
				Error: "Invalid or expired access token",
			}
			helper.WriteJSONResponse(w, http.StatusUnauthorized, apiResponse)
			return
		}

		// Extract user information from claims
		userId := claims.UserId // Adjust according to your claims structure

		// Convert userId to int
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			apiResponse := response.APIResponse{
				Code:  responsecode.CodeAuthenticationError,
				Error: "Invalid user ID format: " + err.Error(),
			}
			helper.WriteJSONResponse(w, http.StatusUnauthorized, apiResponse)
			return
		}

		// Inject user data into the request context
		ctx := context.WithValue(r.Context(), types.ContextKeyUser, contextmodel.User{
			Id: userIdInt,
		})

		next(w, r.WithContext(ctx))
	}

}
