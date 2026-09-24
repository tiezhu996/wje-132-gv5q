package service

import (
	"errors"
	"log/slog"
	"strings"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户业务逻辑。
type UserService struct {
	repo   *repository.UserRepository
	logger *slog.Logger
}

// NewUserService 构造用户服务。
func NewUserService(repo *repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

// Register 注册用户。
func (s *UserService) Register(phone, password, name, role string) (*model.User, error) {
	phone = strings.TrimSpace(phone)
	if len(password) < 6 {
		return nil, util.NewAppError(constants.CodeValidationFailed, "User[phone="+phone+"] register: password too short")
	}
	if _, err := s.repo.FindByPhone(phone); err == nil {
		return nil, util.NewAppError(constants.CodeUserExists, constants.MsgPhoneExists)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, util.Wrap(err, "User[phone=%s] register check failed", phone)
	}
	if role == "" {
		role = constants.RoleWorker
	}
	if !contains(constants.UserRoleValues, role) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "User[role="+role+"] register: invalid role")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.Wrap(err, "User[phone=%s] register hash failed", phone)
	}
	u := &model.User{Phone: phone, PasswordHash: string(hash), Name: name, Role: role}
	if err := s.repo.Create(u); err != nil {
		return nil, util.Wrap(err, "User[phone=%s] register create failed", phone)
	}
	s.logger.Info(constants.LogUserRegisterSuccess, "user_id", u.ID)
	return u, nil
}

// Login 登录并返回 JWT。
func (s *UserService) Login(secret string, expireHours int, phone, password string) (string, *model.User, error) {
	u, err := s.repo.FindByPhone(phone)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", nil, util.NewAppError(constants.CodeInvalidCredentials, constants.MsgInvalidCredentials)
		}
		return "", nil, util.Wrap(err, "User[phone=%s] login find failed", phone)
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		s.logger.Warn(constants.LogUserLoginFailed, "phone", phone)
		return "", nil, util.NewAppError(constants.CodeInvalidCredentials, constants.MsgInvalidCredentials)
	}
	token, err := util.GenerateToken(secret, util.DurationHours(expireHours), u.ID, u.Phone, u.Role)
	if err != nil {
		return "", nil, util.Wrap(err, "User[phone=%s] login token failed", phone)
	}
	s.logger.Info(constants.LogUserLoginSuccess, "user_id", u.ID)
	return token, u, nil
}

// GetByID 查询用户。
func (s *UserService) GetByID(id uint64) (*model.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "User[id=%d] get failed", id)
	}
	return u, nil
}

// UpdateProfile 修改资料。
func (s *UserService) UpdateProfile(id uint64, name, avatar string) (*model.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "User[id=%d] update profile find failed", id)
	}
	if name != "" {
		u.Name = name
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	if err := s.repo.Update(u); err != nil {
		return nil, util.Wrap(err, "User[id=%d] update profile save failed", id)
	}
	s.logger.Info(constants.LogUserProfileUpdate, "user_id", u.ID)
	return u, nil
}

// List 分页查询用户。
func (s *UserService) List(page, pageSize int) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize)
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
