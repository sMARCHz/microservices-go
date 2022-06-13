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
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(map[string]string) *errs.AppError
	Refresh(dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (a DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	if validationErr := req.Validate(); validationErr != nil {
		return nil, validationErr
	}

	// Get the user detail in the database
	var login *domain.Login
	var appErr *errs.AppError
	if login, appErr = a.repo.ValidateUser(req.Username, req.Password); appErr != nil {
		return nil, appErr
	}

	// Build claims and authToken for generate token
	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	// Generate access and refresh token
	var accessToken, refreshToken string
	if accessToken, appErr = authToken.GenerateAccessToken(); appErr != nil {
		return nil, appErr
	}
	if refreshToken, appErr = a.repo.GenerateAndSaveRefreshToken(authToken); appErr != nil {
		return nil, appErr
	}
	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	if validationErr := validateUrlParams(urlParams); validationErr != nil {
		return validationErr
	}

	// Parse string in the urlParams to token
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthorizationError(err.Error())
	} else {
		if jwtToken.Valid {
			claims := jwtToken.Claims.(*domain.AccessTokenClaims)
			if claims.IsUserRole() {
				// Check if token is belonged to its user or not (same customer_id)
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

func (a DefaultAuthService) Refresh(req dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if validationErr := req.Validate(); validationErr != nil {
		return nil, validationErr
	}

	// Check if access token is valid
	if validationErr := req.IsAccessTokenValid(); validationErr != nil {
		if validationErr.Errors == jwt.ValidationErrorExpired {
			// If access token expired, check refresh token is existed in the store or not
			var appErr *errs.AppError
			if appErr = a.repo.IsRefreshTokenExisted(req.RefreshToken); appErr != nil {
				return nil, appErr
			}

			var accessToken string
			if accessToken, appErr = domain.GenerateAccessTokenFromRefreshToken(req.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errs.NewAuthenticationError("invalid token")
	}
	return nil, errs.NewAuthenticationError("cannot generate new access token until the current one expires")
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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
