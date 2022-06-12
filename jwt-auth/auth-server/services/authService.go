package services

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/domain"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/dto"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (a DefaultAuthService) Login(req dto.LoginRequest) (*string, *errs.AppError) {
	// Get the user detail if user is validated
	user, err := a.repo.ValidateUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := user.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	if err := validateUrlParams(urlParams); err != nil {
		return err
	}
	// Parse token from params to jwt.Token
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthorizationError(err.Error())
	} else {
		if jwtToken.Valid {
			claims := jwtToken.Claims.(*domain.Claims)
			if claims.IsUserRole() {
				// Check if token is belonged to its user (same customer_id)
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewAuthorizationError("request not verified with the token claims")
				}
			}
			// Check if user have permissions to the resource
			isAuthorized := a.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
			}
			return nil
		} else {
			return errs.NewAuthorizationError("invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func validateUrlParams(urlParams map[string]string) *errs.AppError {
	if urlParams["token"] == "" {
		return errs.NewAuthorizationError("missing token")
	} else if urlParams["routeName"] == "" {
		return errs.NewAuthorizationError("missing routeName")
	}
	return nil
}

func NewAuthService(repo domain.AuthRepository, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo, permissions}
}
