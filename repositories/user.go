package repositories

import (
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) error
	FindByUserId(userId string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	DoesEmailExist(email string) (bool, error)
}

type userRepo struct {
	db *gorm.DB
}

// NewUserRepo will instantiate User Repository
func NewUserRepo() IUserRepository {
	return &userRepo{
		db: database.DB(),
	}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByUserId(userId string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) DoesEmailExist(email string) (bool, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}
