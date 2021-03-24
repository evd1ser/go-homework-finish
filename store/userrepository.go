package store

import (
	"errors"
	"fmt"
	"github.com/evd1ser/go-homework-finish/internal/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	store *Store
}

func NewUserRepository(s *Store) *UserRepository {
	s.db.AutoMigrate(&models.User{})

	return &UserRepository{
		store: s,
	}
}

//Create user in database
func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	//пытаемся найти пользователя в зарегистрированных
	_, found, err := ur.FindByLogin(u.Username)

	// если произошла ошибка в базе
	if err != nil {
		return nil, err
	}

	// если пользователь уже найден возвращаем ошибку
	if found {
		return nil, errors.New("user already exists")
	}

	// все ок - создаем пользователя
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("something wrong")
	}

	u.Password = string(hashedPassword)
	res := ur.store.db.Create(u)

	if res.Error != nil {
		return nil, res.Error
	}

	return u, nil
}

//Find by login
func (ur *UserRepository) FindByLogin(username string) (*models.User, bool, error) {
	var user models.User

	fmt.Println(username)

	res := ur.store.db.First(&user, "username = ?", username)
	var founded bool

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, founded, nil
	}

	if res.Error != nil {
		return nil, founded, res.Error
	}
	founded = true

	return &user, founded, nil
}
