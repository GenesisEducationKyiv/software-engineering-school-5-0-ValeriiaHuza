package mailer

type Frequency string

const (
	FrequencyHourly Frequency = "hourly"
	FrequencyDaily  Frequency = "daily"
)

type SubscriptionDTO struct {
	Email     string
	City      string
	Frequency Frequency
	Token     string
	Confirmed bool
}

type EmailType string

const (
	EmailTypeCreateSubscription EmailType = "CreateSubscription"
	EmailTypeConfirmSuccess     EmailType = "ConfirmSuccess"
)

type EmailJob struct {
	To           string
	EmailType    EmailType
	Subscription SubscriptionDTO
}

type WeatherUpdateJob struct {
	To           string
	EmailType    string
	Weather      WeatherDTO
	Subscription SubscriptionDTO
}

type WeatherDTO struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}
