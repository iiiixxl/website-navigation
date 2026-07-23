package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"sitenavigation/model"
	"sitenavigation/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	redisHelper *utils.RedisHelper
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:          db,
		redisHelper: utils.NewRedisHelper(),
	}
}

func (s *UserService) getUserCacheKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}

func (s *UserService) getUserByAccountCacheKey(account string) string {
	return fmt.Sprintf("user:account:%s", account)
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func (s *UserService) CreateUser(req *model.CreateUserRequest, createUserID *int) (*model.User, error) {
	var count int64
	if err := s.db.Model(&model.User{}).Where("account = ?", req.Account).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("账号已存在")
	}

	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Account:      req.Account,
		Password:     string(hashedPassword),
		Salt:         salt,
		Name:         req.Name,
		EnName:       req.EnName,
		CreateTime:   time.Now(),
		CreateUserID: createUserID,
		IsDeleted:    false,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id int) (*model.User, error) {
	cacheKey := s.getUserCacheKey(id)

	var user model.User
	if err := s.redisHelper.GetJSON(cacheKey, &user); err == nil {
		return &user, nil
	}

	user = model.User{}
	if err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	_ = s.redisHelper.SetJSON(cacheKey, &user, time.Hour)
	return &user, nil
}

func (s *UserService) GetUserByAccount(account string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("account = ? AND is_deleted = ?", account, false).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) VerifyPassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.Salt))
	return err == nil
}

func (s *UserService) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	if err := s.db.Where("is_deleted = ?", false).Order("id").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(id int, req *model.UpdateUserRequest) (*model.User, error) {
	currentUser, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":    req.Name,
		"en_name": req.EnName,
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+currentUser.Salt), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	if err := s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.clearUserCache(id, currentUser.Account)
	return s.GetUserByID(id)
}

func (s *UserService) DeleteUser(id int) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Account == "admin" {
		return errors.New("不允许删除 admin 用户")
	}

	if err := s.db.Model(&model.User{}).Where("id = ?", id).Update("is_deleted", true).Error; err != nil {
		return err
	}

	s.clearUserCache(id, user.Account)
	return nil
}

func (s *UserService) clearUserCache(id int, account string) {
	_ = s.redisHelper.Del(s.getUserCacheKey(id), s.getUserByAccountCacheKey(account))
}
