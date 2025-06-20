//go:build unit
// +build unit

package subscription

import (
	"errors"
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/stretchr/testify/assert"
)

// --- Mocks ---

type mockWeatherService struct {
	getWeatherFunc func(city string) (*client.WeatherDTO, error)
}

func (m *mockWeatherService) GetWeather(city string) (*client.WeatherDTO, error) {
	return m.getWeatherFunc(city)
}

type mockMailService struct {
	sentConfirmationEmail      *Subscription
	sentConfirmSuccessEmail    *Subscription
	sentWeatherUpdateEmailSub  *Subscription
	sentWeatherUpdateEmailData *client.WeatherDTO
}

func (m *mockMailService) SendConfirmationEmail(sub Subscription) {
	m.sentConfirmationEmail = &sub
}
func (m *mockMailService) SendConfirmSuccessEmail(sub Subscription) {
	m.sentConfirmSuccessEmail = &sub
}
func (m *mockMailService) SendWeatherUpdateEmail(sub Subscription, weather client.WeatherDTO) {
	m.sentWeatherUpdateEmailSub = &sub
	m.sentWeatherUpdateEmailData = &weather
}

type mockSubscriptionRepository struct {
	createFunc                      func(sub Subscription) error
	updateFunc                      func(sub Subscription) error
	findByTokenFunc                 func(token string) (*Subscription, error)
	deleteFunc                      func(sub Subscription) error
	findByEmailFunc                 func(email string) (*Subscription, error)
	findByFrequencyAndConfirmationF func(freq Frequency) ([]Subscription, error)
}

func (m *mockSubscriptionRepository) Create(sub Subscription) error {
	return m.createFunc(sub)
}
func (m *mockSubscriptionRepository) Update(sub Subscription) error {
	return m.updateFunc(sub)
}
func (m *mockSubscriptionRepository) FindByToken(token string) (*Subscription, error) {
	return m.findByTokenFunc(token)
}
func (m *mockSubscriptionRepository) Delete(sub Subscription) error {
	return m.deleteFunc(sub)
}
func (m *mockSubscriptionRepository) FindByEmail(email string) (*Subscription, error) {
	return m.findByEmailFunc(email)
}
func (m *mockSubscriptionRepository) FindByFrequencyAndConfirmation(freq Frequency) ([]Subscription, error) {
	return m.findByFrequencyAndConfirmationF(freq)
}

// --- Tests ---

func TestSubscribeForWeatherUpdates_Success(t *testing.T) {
	mockWeather := &mockWeatherService{
		getWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			return &client.WeatherDTO{}, nil
		},
	}
	mockMail := &mockMailService{}
	mockRepo := &mockSubscriptionRepository{
		findByEmailFunc: func(email string) (*Subscription, error) {
			return nil, nil // not subscribed
		},
		createFunc: func(sub Subscription) error {
			return nil
		},
	}

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	email := "test@example.com"
	city := "Kyiv"
	freq := Frequency("daily")

	err := service.SubscribeForWeatherUpdates(email, city, freq)
	assert.NoError(t, err)
	assert.NotNil(t, mockMail.sentConfirmationEmail)
	assert.Equal(t, email, mockMail.sentConfirmationEmail.Email)
	assert.Equal(t, city, mockMail.sentConfirmationEmail.City)
	assert.Equal(t, freq, mockMail.sentConfirmationEmail.Frequency)
	assert.False(t, mockMail.sentConfirmationEmail.Confirmed)
	assert.NotEmpty(t, mockMail.sentConfirmationEmail.Token)
}

func TestSubscribeForWeatherUpdates_WeatherServiceError(t *testing.T) {
	mockWeather := &mockWeatherService{
		getWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			return nil, errors.New("weather error")
		},
	}
	mockMail := &mockMailService{}
	mockRepo := &mockSubscriptionRepository{}

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))
	assert.EqualError(t, err, "weather error")
	assert.Nil(t, mockMail.sentConfirmationEmail)
}

func TestSubscribeForWeatherUpdates_EmailAlreadySubscribed(t *testing.T) {
	mockWeather := &mockWeatherService{
		getWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			return &client.WeatherDTO{}, nil
		},
	}
	mockMail := &mockMailService{}
	mockRepo := &mockSubscriptionRepository{
		findByEmailFunc: func(email string) (*Subscription, error) {
			return &Subscription{Email: email}, nil // already subscribed
		},
	}

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))
	assert.Equal(t, ErrEmailAlreadySubscribed, err)
	assert.Nil(t, mockMail.sentConfirmationEmail)
}

func TestSubscribeForWeatherUpdates_FindByEmailError(t *testing.T) {
	mockWeather := &mockWeatherService{
		getWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			return &client.WeatherDTO{}, nil
		},
	}
	mockMail := &mockMailService{}
	mockRepo := &mockSubscriptionRepository{
		findByEmailFunc: func(email string) (*Subscription, error) {
			return nil, errors.New("db error")
		},
	}

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))
	assert.Equal(t, ErrInvalidInput, err)
	assert.Nil(t, mockMail.sentConfirmationEmail)
}

func TestSubscribeForWeatherUpdates_CreateError(t *testing.T) {
	mockWeather := &mockWeatherService{
		getWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			return &client.WeatherDTO{}, nil
		},
	}
	mockMail := &mockMailService{}
	mockRepo := &mockSubscriptionRepository{
		findByEmailFunc: func(email string) (*Subscription, error) {
			return nil, nil
		},
		createFunc: func(sub Subscription) error {
			return errors.New("db error")
		},
	}

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))
	assert.Equal(t, ErrFailedToSaveSubscription, err)
	assert.Nil(t, mockMail.sentConfirmationEmail)
}
func TestConfirmSubscription_Success(t *testing.T) {
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: false,
	}
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			assert.Equal(t, "token123", token)
			return mockSub, nil
		},
		updateFunc: func(sub Subscription) error {
			assert.True(t, sub.Confirmed)
			assert.Equal(t, mockSub.Email, sub.Email)
			return nil
		},
	}
	mockMail := &mockMailService{}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		mailService:            mockMail,
	}

	err := service.ConfirmSubscription("token123")
	assert.NoError(t, err)
	assert.NotNil(t, mockMail.sentConfirmSuccessEmail)
	assert.Equal(t, "test@example.com", mockMail.sentConfirmSuccessEmail.Email)
	assert.True(t, mockMail.sentConfirmSuccessEmail.Confirmed)
}

func TestConfirmSubscription_TokenNotFound(t *testing.T) {
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			return nil, errors.New("not found")
		},
	}
	mockMail := &mockMailService{}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		mailService:            mockMail,
	}

	err := service.ConfirmSubscription("invalid-token")
	assert.Equal(t, ErrTokenNotFound, err)
	assert.Nil(t, mockMail.sentConfirmSuccessEmail)
}

func TestConfirmSubscription_UpdateError(t *testing.T) {
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: false,
	}
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			return mockSub, nil
		},
		updateFunc: func(sub Subscription) error {
			return errors.New("update error")
		},
	}
	mockMail := &mockMailService{}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		mailService:            mockMail,
	}

	err := service.ConfirmSubscription("token123")
	assert.Equal(t, ErrFailedToSaveSubscription, err)
	assert.Nil(t, mockMail.sentConfirmSuccessEmail)
}
func TestUnsubscribe_Success(t *testing.T) {
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: true,
	}
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			assert.Equal(t, "token123", token)
			return mockSub, nil
		},
		deleteFunc: func(sub Subscription) error {
			assert.Equal(t, mockSub.Email, sub.Email)
			return nil
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("token123")
	assert.NoError(t, err)
}

func TestUnsubscribe_TokenNotFound(t *testing.T) {
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			return nil, errors.New("not found")
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("invalid-token")
	assert.Equal(t, ErrTokenNotFound, err)
}

func TestUnsubscribe_SubscriptionNil(t *testing.T) {
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			return nil, nil
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("token123")
	assert.NoError(t, err)
}

func TestUnsubscribe_DeleteError(t *testing.T) {
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: true,
	}
	mockRepo := &mockSubscriptionRepository{
		findByTokenFunc: func(token string) (*Subscription, error) {
			return mockSub, nil
		},
		deleteFunc: func(sub Subscription) error {
			return errors.New("delete error")
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("token123")
	assert.Equal(t, ErrInvalidInput, err)
}
func TestEmailSubscribed_ReturnsTrueWhenSubscribed(t *testing.T) {
	mockRepo := &mockSubscriptionRepository{
		findByEmailFunc: func(email string) (*Subscription, error) {
			return &Subscription{Email: email}, nil
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	subscribed, err := service.emailSubscribed("test@example.com")
	assert.True(t, subscribed)
	assert.NoError(t, err)
}

func TestEmailSubscribed_ReturnsError(t *testing.T) {
	mockRepo := &mockSubscriptionRepository{
		findByEmailFunc: func(email string) (*Subscription, error) {
			return nil, errors.New("db error")
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	subscribed, err := service.emailSubscribed("test@example.com")
	assert.False(t, subscribed)
	assert.EqualError(t, err, "db error")
}

func TestGetConfirmedSubscriptionsByFrequency_ReturnsSubscriptions(t *testing.T) {
	expectedSubs := []Subscription{
		{Email: "a@example.com", City: "Kyiv", Frequency: Frequency("daily"), Confirmed: true},
		{Email: "b@example.com", City: "Lviv", Frequency: Frequency("daily"), Confirmed: true},
	}
	mockRepo := &mockSubscriptionRepository{
		findByFrequencyAndConfirmationF: func(freq Frequency) ([]Subscription, error) {
			assert.Equal(t, Frequency("daily"), freq)
			return expectedSubs, nil
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	result := service.GetConfirmedSubscriptionsByFrequency(Frequency("daily"))
	assert.Equal(t, expectedSubs, result)
}

func TestGetConfirmedSubscriptionsByFrequency_RepoError_ReturnsEmptySlice(t *testing.T) {
	mockRepo := &mockSubscriptionRepository{
		findByFrequencyAndConfirmationF: func(freq Frequency) ([]Subscription, error) {
			return nil, errors.New("db error")
		},
	}
	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	result := service.GetConfirmedSubscriptionsByFrequency(Frequency("daily"))
	assert.Empty(t, result)
}
