package repository

import (
	"github.com/ValeriiaHuza/weather_api/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub models.Subscription) error
	Update(sub models.Subscription) error
	FindByToken(token string) (*models.Subscription, error)
	Delete(sub models.Subscription) error
	FindByEmail(email string) (*models.Subscription, error)
	FindByFrequencyAndConfirmation(freq models.Frequency) ([]models.Subscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(database *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: database}
}

func (r *subscriptionRepository) Create(sub models.Subscription) error {
	return r.db.Create(&sub).Error
}

func (r *subscriptionRepository) Update(sub models.Subscription) error {
	return r.db.Save(&sub).Error
}

func (r *subscriptionRepository) FindByToken(token string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("token = ?", token).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) Delete(sub models.Subscription) error {
	return r.db.Unscoped().Delete(&sub).Error
}

func (r *subscriptionRepository) FindByEmail(email string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("email = ?", email).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) FindByFrequencyAndConfirmation(freq models.Frequency) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := r.db.Where("frequency = ? AND confirmed = true", freq).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}
