package repository

import (
	"github.com/ValeriiaHuza/weather_api/db"
	"github.com/ValeriiaHuza/weather_api/models"
)

type SubscriptionRepository interface {
	Create(sub models.Subscription) error
	Update(sub models.Subscription) error
	FindByToken(token string) (*models.Subscription, error)
	Delete(sub models.Subscription) error
	FindByEmail(email string) (*models.Subscription, error)
	FindByFrequencyAndConfirmation(freq models.Frequency) ([]models.Subscription, error)
}

type subscriptionRepository struct{}

func NewSubscriptionRepository() SubscriptionRepository {
	return &subscriptionRepository{}
}

func (r *subscriptionRepository) Create(sub models.Subscription) error {
	return db.DB.Create(&sub).Error
}

func (r *subscriptionRepository) Update(sub models.Subscription) error {
	return db.DB.Save(&sub).Error
}

func (r *subscriptionRepository) FindByToken(token string) (*models.Subscription, error) {
	var sub models.Subscription
	err := db.DB.Where("token = ?", token).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) Delete(sub models.Subscription) error {
	return db.DB.Unscoped().Delete(&sub).Error
}

func (r *subscriptionRepository) FindByEmail(email string) (*models.Subscription, error) {
	var sub models.Subscription
	err := db.DB.Where("email = ?", email).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) FindByFrequencyAndConfirmation(freq models.Frequency) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := db.DB.Where("frequency = ? AND confirmed = true", freq).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}
