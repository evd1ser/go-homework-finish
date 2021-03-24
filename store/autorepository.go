package store

import (
	"errors"
	"github.com/evd1ser/go-homework-finish/internal/app/models"
	"gorm.io/gorm"
)

type AutoRepository struct {
	store *Store
}

func NewAutoRepository(s *Store) *AutoRepository {
	s.db.AutoMigrate(&models.Auto{})

	return &AutoRepository{
		store: s,
	}
}

//Create auto in database
func (ar *AutoRepository) Create(a *models.Auto) (*models.Auto, bool, error) {
	//пытаемся найти пользователя в зарегистрированных
	_, found, err := ar.GetByMark(a.Mark)
	var exist bool
	// если произошла ошибка в базе
	if err != nil {
		return nil, exist, err
	}

	// если пользователь уже найден возвращаем ошибку
	if found {
		return nil, exist, nil
	}

	res := ar.store.db.Create(a)

	if res.Error != nil {
		return nil, exist, res.Error
	}

	exist = true

	return a, exist, nil
}

//Update auto in database
func (ar *AutoRepository) Update(a *models.Auto) (*models.Auto, error) {
	//пытаемся найти auto в зарегистрированных
	_, found, err := ar.GetByMark(a.Mark)

	// если произошла ошибка в базе
	if err != nil {
		return nil, err
	}

	// если пользователь уже найден возвращаем ошибку
	if !found {
		return nil, errors.New("auto not exists")
	}

	res := ar.store.db.Model(a).Updates(a)

	if res.Error != nil {
		return nil, res.Error
	}

	return a, nil
}

//Delete auto in database
func (ar *AutoRepository) Delete(a *models.Auto) (*models.Auto, error) {
	//пытаемся найти auto в зарегистрированных
	_, found, err := ar.GetByMark(a.Mark)

	// если произошла ошибка в базе
	if err != nil {
		return nil, err
	}

	// если пользователь уже найден возвращаем ошибку
	if !found {
		return nil, errors.New("auto not exists")
	}

	res := ar.store.db.Delete(&a)

	if res.Error != nil {
		return nil, res.Error
	}

	return a, nil
}

func (ar *AutoRepository) GetByMark(mark string) (*models.Auto, bool, error) {
	var auto models.Auto

	res := ar.store.db.First(&auto, "mark = ?", mark)

	var founded bool

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, founded, nil
	}

	if res.Error != nil {
		return nil, founded, res.Error
	}

	founded = true

	return &auto, true, nil
}

func (ar *AutoRepository) GetAll() ([]models.AutoApi, error) {

	var modelsRes = make([]models.AutoApi, 0)
	res := ar.store.db.Model(&models.Auto{}).Scan(&modelsRes)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return modelsRes, nil
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return modelsRes, nil
}
