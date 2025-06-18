package repository

import (
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(database *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: database}
}

func (r *SubscriptionRepository) Create(sub subscription.Subscription) error {
	return r.db.Create(&sub).Error
}

func (r *SubscriptionRepository) Update(sub subscription.Subscription) error {
	return r.db.Save(&sub).Error
}

func (r *SubscriptionRepository) FindByToken(token string) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	err := r.db.Where("token = ?", token).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Delete(sub subscription.Subscription) error {
	return r.db.Unscoped().Delete(&sub).Error
}

func (r *SubscriptionRepository) FindByEmail(email string) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	err := r.db.Where("email = ?", email).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) FindByFrequencyAndConfirmation(
	freq subscription.Frequency) ([]subscription.Subscription, error) {
	var subs []subscription.Subscription
	err := r.db.Where("frequency = ? AND confirmed = true", freq).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}
