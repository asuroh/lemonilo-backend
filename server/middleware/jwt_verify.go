package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"lemonilo-backend/model"
	apiHandler "lemonilo-backend/server/handler"
	"lemonilo-backend/usecase"
)

type jwtClaims struct {
	jwt.StandardClaims
}

// VerifyMiddlewareInit ...
type VerifyMiddlewareInit struct {
	*usecase.ContractUC
}

// VerifyPermissionInit ...
type VerifyPermissionInit struct {
	*usecase.ContractUC
	Menu string
}

func userContextInterface(ctx context.Context, req *http.Request, subject string, body map[string]interface{}) context.Context {
	return context.WithValue(ctx, subject, body)
}

func (m VerifyMiddlewareInit) verifyJWT(r *http.Request) (res map[string]interface{}, err error) {
	claims := &jwtClaims{}

	tokenAuthHeader := r.Header.Get("Authorization")
	if !strings.Contains(tokenAuthHeader, "Bearer") {
		return res, errors.New("Invalid token")
	}
	tokenAuth := strings.Replace(tokenAuthHeader, "Bearer ", "", -1)

	_, err = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := m.ContractUC.EnvConfig["TOKEN_SECRET"]
		return []byte(secret), nil
	})
	if err != nil {
		return res, errors.New("Invalid Token!")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return res, errors.New("Expired Token!")
	}

	// Decrypt payload
	res, err = m.ContractUC.Jwe.Rollback(claims.Id)
	if err != nil {
		return res, errors.New("Error when load the payload!")
	}

	return res, nil
}

func (m VerifyMiddlewareInit) verifyRefreshJWT(r *http.Request) (res map[string]interface{}, err error) {
	claims := &jwtClaims{}

	tokenAuthHeader := r.Header.Get("Authorization")
	if !strings.Contains(tokenAuthHeader, "Bearer") {
		return res, errors.New("Invalid token")
	}
	tokenAuth := strings.Replace(tokenAuthHeader, "Bearer ", "", -1)

	_, err = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := m.ContractUC.EnvConfig["TOKEN_REFRESH_SECRET"]
		return []byte(secret), nil
	})
	if err != nil {
		return res, errors.New("Invalid Token!")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return res, errors.New("Expired Token!")
	}

	// Decrypt payload
	res, err = m.ContractUC.Jwe.Rollback(claims.Id)
	if err != nil {
		return res, errors.New("Error when load the payload!")
	}

	return res, nil
}

// VerifyAdminTokenCredential ...
func (m VerifyMiddlewareInit) VerifyAdminTokenCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jweRes, err := m.verifyJWT(r)
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, err.Error(), []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Check id in table
		userUC := usecase.UserUC{ContractUC: m.ContractUC}
		user, err := userUC.FindByID(jweRes["id"].(string), false)
		if user.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		if user.Role != model.RoleCodeAdmin {
			apiHandler.RespondWithJSON(w, 401, 401, "Admin only access!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["username"] = user.UserName
		jweRes["email"] = user.Email
		jweRes["role"] = user.Role

		ctx := userContextInterface(r.Context(), r, "admin", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifyUserTokenCredential ...
func (m VerifyMiddlewareInit) VerifyUserTokenCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jweRes, err := m.verifyJWT(r)
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, err.Error(), []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Check id in table
		userUC := usecase.UserUC{ContractUC: m.ContractUC}
		user, err := userUC.FindByID(jweRes["id"].(string), false)
		if user.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["username"] = user.UserName
		jweRes["email"] = user.Email
		jweRes["role"] = user.Role

		ctx := userContextInterface(r.Context(), r, "user", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
