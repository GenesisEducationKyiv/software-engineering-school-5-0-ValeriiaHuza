//go:build unit
// +build unit

package service

import (
	"testing"

	"github.com/ValeriiaHuza/weather_api/dto"
	appErr "github.com/ValeriiaHuza/weather_api/error"
	"github.com/ValeriiaHuza/weather_api/models"
	"github.com/ValeriiaHuza/weather_api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateSubscriptionInput_ValidInput(t *testing.T) {
	mockWeather := new(MockWeatherService)
	mockRepo := new(MockRepo)

	mockWeather.On("GetWeather", "Kyiv").Return(&dto.WeatherDTO{}, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)

	svc := service.NewSubscribeService(mockWeather, nil, mockRepo)

	freq, err := svc.ValidateSubscriptionInput("test@example.com", "Kyiv", "daily")

	assert.Nil(t, err)
	assert.Equal(t, models.Frequency("daily"), freq)
}

func TestValidateSubscriptionInput_InvalidEmail(t *testing.T) {
	mockWeather := new(MockWeatherService)
	mockRepo := new(MockRepo)
	svc := service.NewSubscribeService(mockWeather, nil, mockRepo)

	_, err := svc.ValidateSubscriptionInput("invalid-email", "Kyiv", "daily")

	assert.Equal(t, appErr.ErrInvalidInput, err)
}

func TestValidateSubscriptionInput_EmptyFields(t *testing.T) {
	svc := service.NewSubscribeService(nil, nil, nil)

	_, err := svc.ValidateSubscriptionInput("", "Kyiv", "daily")
	assert.Equal(t, appErr.ErrInvalidInput, err)

	_, err = svc.ValidateSubscriptionInput("test@example.com", "", "daily")
	assert.Equal(t, appErr.ErrInvalidInput, err)

	_, err = svc.ValidateSubscriptionInput("test@example.com", "Kyiv", "")
	assert.Equal(t, appErr.ErrInvalidInput, err)
}

func TestValidateSubscriptionInput_WeatherServiceFails(t *testing.T) {
	mockWeather := new(MockWeatherService)
	mockRepo := new(MockRepo)

	mockWeather.On("GetWeather", "InvalidCity").Return(nil, appErr.ErrInvalidInput)

	svc := service.NewSubscribeService(mockWeather, nil, mockRepo)

	_, err := svc.ValidateSubscriptionInput("test@example.com", "InvalidCity", "daily")

	assert.Equal(t, appErr.ErrInvalidInput, err)
}

func TestValidateSubscriptionInput_InvalidFrequency(t *testing.T) {
	mockWeather := new(MockWeatherService)
	mockRepo := new(MockRepo)

	mockWeather.On("GetWeather", "Kyiv").Return(&dto.WeatherDTO{}, nil)

	svc := service.NewSubscribeService(mockWeather, nil, mockRepo)

	_, err := svc.ValidateSubscriptionInput("test@example.com", "Kyiv", "yearly")

	assert.Equal(t, appErr.ErrInvalidInput, err)
}

func TestValidateSubscriptionInput_EmailAlreadySubscribed(t *testing.T) {
	mockWeather := new(MockWeatherService)
	mockRepo := new(MockRepo)

	mockWeather.On("GetWeather", "Kyiv").Return(&dto.WeatherDTO{}, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(&models.Subscription{}, nil)

	svc := service.NewSubscribeService(mockWeather, nil, mockRepo)

	_, err := svc.ValidateSubscriptionInput("test@example.com", "Kyiv", "daily")

	assert.Equal(t, appErr.ErrEmailSubscribed, err)
}

func TestValidateSubscriptionInput_RepoError(t *testing.T) {
	mockWeather := new(MockWeatherService)
	mockRepo := new(MockRepo)

	mockWeather.On("GetWeather", "Kyiv").Return(&dto.WeatherDTO{}, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, appErr.New(404, "db error"))

	svc := service.NewSubscribeService(mockWeather, nil, mockRepo)

	_, err := svc.ValidateSubscriptionInput("test@example.com", "Kyiv", "daily")

	assert.Equal(t, appErr.ErrInvalidInput, err)
}

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeather(city string) (*dto.WeatherDTO, *appErr.AppError) {
	args := m.Called(city)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*appErr.AppError)
	}
	return args.Get(0).(*dto.WeatherDTO), nil
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindByEmail(email string) (*models.Subscription, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), nil
}

// Unused methods
func (m *MockRepo) Create(models.Subscription) error                 { return nil }
func (m *MockRepo) Update(models.Subscription) error                 { return nil }
func (m *MockRepo) Delete(models.Subscription) error                 { return nil }
func (m *MockRepo) FindByToken(string) (*models.Subscription, error) { return nil, nil }
func (m *MockRepo) FindByFrequencyAndConfirmation(models.Frequency) ([]models.Subscription, error) {
	return nil, nil
}
