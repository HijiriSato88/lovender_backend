package service

import (
	"errors"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
	"lovender_backend/pkg/crypto"
	"lovender_backend/pkg/jwtutil"
)

type UserService interface {
	GetUser(id int64) (*models.User, error)
	Register(req *models.RegisterRequest) (*models.RegisterResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(id int64) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *userService) Register(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// メールアドレスの重複チェック
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// パスワードをハッシュ化
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// ユーザー作成
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// JWTトークン生成
	token, err := jwtutil.GenerateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	// 再度DBから取得してcreated_atを取得
	createdUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.RegisterResponse{
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		Password:  req.Password, // 元のパスワードを返す
		CreatedAt: createdUser.CreatedAt,
		Token:     token,
	}, nil
}

func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// ユーザーを取得
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// パスワード確認
	if !crypto.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// JWTトークン生成
	token, err := jwtutil.GenerateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
	}, nil
}
