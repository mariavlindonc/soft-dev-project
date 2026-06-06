package db

import (
	"backend/domain"

	"gorm.io/gorm"
)

type UserDAOImpl struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAOImpl {
	return &UserDAOImpl{db: db}
}

func (d *UserDAOImpl) Create(user *domain.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAOImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := d.db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAOImpl) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := d.db.Take(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

var _ UserDAO = (*UserDAOImpl)(nil)
