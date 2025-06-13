package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Frequency string

const (
	FrequencyHourly Frequency = "hourly"
	FrequencyDaily  Frequency = "daily"
)

type Subscription struct {
	gorm.Model
	ID        int       `gorm:"primary_key"`
	Email     string    `gorm:"unique;not null"`
	City      string    `gorm:"not null"`
	Frequency Frequency `gorm:"type:varchar(10);not null"`
	Token     string    `gorm:"unique;not null"`
	Confirmed bool      `gorm:"not null;default:false"`
}

func ParseFrequency(freq string) (Frequency, error) {
	switch freq {
	case string(FrequencyHourly):
		return FrequencyHourly, nil
	case string(FrequencyDaily):
		return FrequencyDaily, nil
	default:
		return "", fmt.Errorf("invalid frequency: %s", freq)
	}
}
