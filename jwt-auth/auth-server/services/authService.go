package services

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/domain"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/dto"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/logger"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
	Verify(urlParams map[string]string) (bool, error)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (a DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
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

func (a DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {
	// Parse token from params to jwt.Token
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return false, err
	} else {
		if jwtToken.Valid {
			mapClaims := jwtToken.Claims.(jwt.MapClaims)
			// Build domain.Claims from jwt.MapClaims
			if claims, err := domain.BuildClaimsFromJwtMapClaims(mapClaims); err != nil {
				return false, err
			} else {
				if claims.IsUserRole() {
					// Check if token is belonged to its user (same customer_id)
					if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
						return false, nil
					}
				}
				// Check if user have permissions to the resource
				isAuthorized := a.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
				return isAuthorized, nil
			}
		} else {
			return false, errors.New("invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthService(repo domain.AuthRepository, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo, permissions}
}
