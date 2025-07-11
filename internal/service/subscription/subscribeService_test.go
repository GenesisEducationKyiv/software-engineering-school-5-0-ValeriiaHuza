//go:build unit
// +build unit

package subscription

import (
	"errors"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type mockWeatherService struct {
	mock.Mock
}

func (m *mockWeatherService) GetWeather(city string) (*client.WeatherDTO, error) {
	args := m.Called(city)
	dto, _ := args.Get(0).(*client.WeatherDTO)
	return dto, args.Error(1)
}

type mockMailService struct {
	mock.Mock
}

func (m *mockMailService) SendConfirmationEmail(sub Subscription) {
	m.Called(sub)
}
func (m *mockMailService) SendConfirmSuccessEmail(sub Subscription) {
	m.Called(sub)
}
func (m *mockMailService) SendWeatherUpdateEmail(sub Subscription, weather client.WeatherDTO) {
	m.Called(sub, weather)
}

type mockSubscriptionRepository struct {
	mock.Mock
}

func (m *mockSubscriptionRepository) Create(sub Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}
func (m *mockSubscriptionRepository) Update(sub Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}
func (m *mockSubscriptionRepository) FindByToken(token string) (*Subscription, error) {
	args := m.Called(token)
	sub, _ := args.Get(0).(*Subscription)
	return sub, args.Error(1)
}
func (m *mockSubscriptionRepository) Delete(sub Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}
func (m *mockSubscriptionRepository) FindByEmail(email string) (*Subscription, error) {
	args := m.Called(email)
	sub, _ := args.Get(0).(*Subscription)
	return sub, args.Error(1)
}
func (m *mockSubscriptionRepository) FindByFrequencyAndConfirmation(freq Frequency) ([]Subscription, error) {
	args := m.Called(freq)
	subs, _ := args.Get(0).([]Subscription)
	return subs, args.Error(1)
}

// --- Tests ---

func TestSubscribeForWeatherUpdates_Success(t *testing.T) {
	mockWeather := new(mockWeatherService)
	mockMail := new(mockMailService)
	mockRepo := new(mockSubscriptionRepository)

	mockWeather.On("GetWeather", "Kyiv").Return(&client.WeatherDTO{}, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("record not found"))
	mockRepo.On("Create", mock.AnythingOfType("Subscription")).Return(nil)
	mockMail.On("SendConfirmationEmail", mock.AnythingOfType("Subscription")).Return()

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
	mockWeather.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockMail.AssertExpectations(t)
}

func TestSubscribeForWeatherUpdates_WeatherServiceError(t *testing.T) {
	mockWeather := new(mockWeatherService)
	mockMail := new(mockMailService)
	mockRepo := new(mockSubscriptionRepository)

	mockWeather.On("GetWeather", "Kyiv").Return(nil, errors.New("weather error"))

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))

	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockMail.AssertNotCalled(t, "SendConfirmationEmail", mock.Anything)
	assert.EqualError(t, err, "weather error")
	mockWeather.AssertExpectations(t)
}

func TestSubscribeForWeatherUpdates_EmailAlreadySubscribed(t *testing.T) {
	mockWeather := new(mockWeatherService)
	mockMail := new(mockMailService)
	mockRepo := new(mockSubscriptionRepository)

	mockWeather.On("GetWeather", "Kyiv").Return(&client.WeatherDTO{}, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(&Subscription{Email: "test@example.com"}, nil)

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))
	assert.Equal(t, ErrEmailAlreadySubscribed, err)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockMail.AssertNotCalled(t, "SendConfirmationEmail", mock.Anything)
	mockWeather.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestSubscribeForWeatherUpdates_CreateError(t *testing.T) {
	mockWeather := new(mockWeatherService)
	mockMail := new(mockMailService)
	mockRepo := new(mockSubscriptionRepository)

	mockWeather.On("GetWeather", "Kyiv").Return(&client.WeatherDTO{}, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("record not found"))
	mockRepo.On("Create", mock.AnythingOfType("Subscription")).Return(errors.New("db error"))

	service := &SubscribeService{
		weatherService:         mockWeather,
		mailService:            mockMail,
		subscriptionRepository: mockRepo,
	}

	err := service.SubscribeForWeatherUpdates("test@example.com", "Kyiv", Frequency("daily"))
	assert.Equal(t, ErrFailedToSaveSubscription, err)
	mockMail.AssertNotCalled(t, "SendConfirmationEmail", mock.Anything)
	mockWeather.AssertExpectations(t)
	mockRepo.AssertExpectations(t)

}
func TestConfirmSubscription_Success(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockMail := new(mockMailService)
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: false,
	}

	mockRepo.On("FindByToken", "token123").Return(mockSub, nil)
	mockRepo.On("Update", mock.MatchedBy(func(sub Subscription) bool {
		return sub.Email == mockSub.Email && sub.Confirmed
	})).Return(nil)
	mockMail.On("SendConfirmSuccessEmail", mock.MatchedBy(func(sub Subscription) bool {
		return sub.Email == mockSub.Email && sub.Confirmed
	})).Return()

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		mailService:            mockMail,
	}

	err := service.ConfirmSubscription("token123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMail.AssertExpectations(t)
}

func TestConfirmSubscription_TokenNotFound(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockMail := new(mockMailService)

	mockRepo.On("FindByToken", "invalid-token").Return(nil, errors.New("not found"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		mailService:            mockMail,
	}

	err := service.ConfirmSubscription("invalid-token")

	mockMail.AssertNotCalled(t, "SendConfirmSuccessEmail", mock.Anything)
	assert.Equal(t, ErrTokenNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestConfirmSubscription_UpdateError(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockMail := new(mockMailService)
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: false,
	}

	mockRepo.On("FindByToken", "token123").Return(mockSub, nil)
	mockRepo.On("Update", mock.AnythingOfType("Subscription")).Return(errors.New("update error"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		mailService:            mockMail,
	}

	err := service.ConfirmSubscription("token123")

	mockMail.AssertNotCalled(t, "SendConfirmSuccessEmail", mock.Anything)
	assert.Equal(t, ErrFailedToSaveSubscription, err)
	mockRepo.AssertExpectations(t)
}
func TestUnsubscribe_Success(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: true,
	}

	mockRepo.On("FindByToken", "token123").Return(mockSub, nil)
	mockRepo.On("Delete", *mockSub).Return(nil)

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("token123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUnsubscribe_TokenNotFound(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockRepo.On("FindByToken", "invalid-token").Return(nil, errors.New("not found"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("invalid-token")
	assert.Equal(t, ErrTokenNotFound, err)

	mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestUnsubscribe_SubscriptionNil(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockRepo.On("FindByToken", "token123").Return(nil, nil)

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("token123")
	assert.NoError(t, err)

	mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestUnsubscribe_DeleteError(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockSub := &Subscription{
		Email:     "test@example.com",
		City:      "Kyiv",
		Frequency: Frequency("daily"),
		Token:     "token123",
		Confirmed: true,
	}

	mockRepo.On("FindByToken", "token123").Return(mockSub, nil)
	mockRepo.On("Delete", *mockSub).Return(errors.New("delete error"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	err := service.Unsubscribe("token123")
	assert.Equal(t, ErrInvalidInput, err)
	mockRepo.AssertExpectations(t)
}
func TestEmailSubscribed_ReturnsTrueWhenSubscribed(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockRepo.On("FindByEmail", "test@example.com").Return(&Subscription{Email: "test@example.com"}, nil)

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	subscribed := service.emailSubscribed("test@example.com")
	assert.True(t, subscribed)

	mockRepo.AssertExpectations(t)
}

func TestEmailSubscribed_ReturnsError(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("db error"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	subscribed := service.emailSubscribed("test@example.com")
	assert.False(t, subscribed)

	mockRepo.AssertExpectations(t)
}

func TestGetConfirmedSubscriptionsByFrequency_ReturnsSubscriptions(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	expectedSubs := []Subscription{
		{Email: "a@example.com", City: "Kyiv", Frequency: Frequency("daily"), Confirmed: true},
		{Email: "b@example.com", City: "Lviv", Frequency: Frequency("daily"), Confirmed: true},
	}
	mockRepo.On("FindByFrequencyAndConfirmation", Frequency("daily")).Return(expectedSubs, nil)

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	result := service.GetConfirmedSubscriptionsByFrequency(Frequency("daily"))
	assert.Equal(t, expectedSubs, result)
	mockRepo.AssertExpectations(t)
}

func TestGetConfirmedSubscriptionsByFrequency_RepoError_ReturnsEmptySlice(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockRepo.On("FindByFrequencyAndConfirmation", Frequency("daily")).Return(nil, errors.New("db error"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
	}

	result := service.GetConfirmedSubscriptionsByFrequency(Frequency("daily"))
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestSendSubscriptionEmails_SendsEmails(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockWeather := new(mockWeatherService)
	mockMail := new(mockMailService)

	subs := []Subscription{
		{Email: "test1@example.com", City: "Kyiv", Frequency: Frequency("daily"), Confirmed: true},
		{Email: "test2@example.com", City: "Lviv", Frequency: Frequency("daily"), Confirmed: true},
	}
	mockRepo.On("FindByFrequencyAndConfirmation", Frequency("daily")).Return(subs, nil)
	mockWeather.On("GetWeather", "Kyiv").Return(&client.WeatherDTO{Temperature: 10}, nil)
	mockWeather.On("GetWeather", "Lviv").Return(&client.WeatherDTO{Temperature: 20}, nil)
	mockMail.On("SendWeatherUpdateEmail", subs[0], client.WeatherDTO{Temperature: 10}).Return()
	mockMail.On("SendWeatherUpdateEmail", subs[1], client.WeatherDTO{Temperature: 20}).Return()

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		weatherService:         mockWeather,
		mailService:            mockMail,
	}

	service.SendSubscriptionEmails(Frequency("daily"))
	mockRepo.AssertExpectations(t)
	mockWeather.AssertExpectations(t)
	mockMail.AssertExpectations(t)
}

func TestSendSubscriptionEmails_WeatherError_SkipsEmail(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	mockWeather := new(mockWeatherService)
	mockMail := new(mockMailService)

	subs := []Subscription{
		{Email: "a@example.com", City: "Kyiv", Frequency: Frequency("daily"), Confirmed: true},
	}
	mockRepo.On("FindByFrequencyAndConfirmation", Frequency("daily")).Return(subs, nil)
	mockWeather.On("GetWeather", "Kyiv").Return(nil, errors.New("weather error"))

	service := &SubscribeService{
		subscriptionRepository: mockRepo,
		weatherService:         mockWeather,
		mailService:            mockMail,
	}

	service.SendSubscriptionEmails(Frequency("daily"))

	mockMail.AssertNotCalled(t, "SendWeatherUpdateEmail", mock.Anything, mock.Anything)

	mockRepo.AssertExpectations(t)
	mockWeather.AssertExpectations(t)
	mockMail.AssertExpectations(t)
}
