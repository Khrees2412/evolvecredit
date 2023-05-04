package services

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/khrees2412/evolvecredit/models"
	"github.com/khrees2412/evolvecredit/repositories"
	"github.com/khrees2412/evolvecredit/types"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

type IAuthService interface {
	Register(body types.RegisterRequest) (*types.RegisterResponse, error)
	Login(body types.LoginRequest) (*types.LoginResponse, error)
	IssueToken(u *models.User) (*types.TokenResponse, error)
	ParseToken(token string) (*types.Claims, error)
}

type authService struct {
	jwtSecret string
	userRepo  repositories.IUserRepository
}

// NewAuthService will instantiate AuthService
func NewAuthService() IAuthService {
	return &authService{
		jwtSecret: os.Getenv("JWT_SECRET"),
		userRepo:  repositories.NewUserRepo(),
	}
}

var (
	errCouldNotSetPassword = errors.New("could not set password")
	errInvalidPassword     = errors.New("invalid password")
)

func (as *authService) Register(body types.RegisterRequest) (*types.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errCouldNotSetPassword
	}
	password := string(hashedPassword)
	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  password,
	}
	err = as.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}, nil
}

func (as *authService) Login(body types.LoginRequest) (*types.LoginResponse, error) {
	user, err := as.userRepo.FindByEmail(body.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return nil, errInvalidPassword
	}

	issueResponse, err := as.IssueToken(user)

	return &types.LoginResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     issueResponse,
	}, nil
}

func (as *authService) IssueToken(u *models.User) (*types.TokenResponse, error) {
	nowTime := time.Now()
	expConf, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))

	if err != nil {
		return nil, err
	}
	expireTime := nowTime.Add(time.Duration(int64(expConf)) * time.Minute)

	claims := types.Claims{
		Email:  u.Email,
		UserId: u.UserId,
		StandardClaims: jwt.StandardClaims{
			Subject:   u.UserId,
			IssuedAt:  nowTime.Unix(),
			ExpiresAt: expireTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := jwtToken.SignedString([]byte(as.jwtSecret))

	if err != nil {
		return nil, err
	}

	return &types.TokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   claims.ExpiresAt,
		Issuer:      claims.Issuer,
	}, nil
}

func (as *authService) ParseToken(token string) (*types.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&types.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(as.jwtSecret), nil
		},
	)

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*types.Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
