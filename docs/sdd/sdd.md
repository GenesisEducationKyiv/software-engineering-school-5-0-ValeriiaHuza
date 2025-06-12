# System Design  
## Weather Subscription API  

### Context 

The Weather Subscription API is a backend service that lets users check the current weather and receive weather updates by email. It uses data from WeatherApi.com and allows users to choose a city and how often they want updates (daily or hourly).

Users subscribe by entering their email, city, and preferred frequency. They receive a confirmation email with a link to activate the subscription. The service sends weather update emails at the selected intervals and includes an unsubscribe link in each email.

### 1) System requirements  

#### Functional requirements:  

- Allow users to query current weather conditions for a given city.  
- Retrieve real-time weather data from WeatherApi.com.  
- Let users subscribe to weather updates by providing their email, city, and frequency (daily/hourly).  
- Send a confirmation email with a link after a subscription request.  
- Allow only one active subscription per email-city pair (no duplicates).  
- Send weather update emails according to the selected frequency.  
- Include an unsubscribe link in every update email.  

#### Nonfunctional requirements:  

- System should be available 99.9% of the time.  
- It must scale to support 100,000+ active subscriptions.  
- Ensure low-latency (<150ms) API responses and minimal delay in sending emails.  
- Achieve 99.99% reliable email delivery with retry logic.  

### 2) Constraints:  

- Minimal budget for third-party services, including:  
  - WeatherAPI.com subscription plan  
  - Cloud/Server hosting  
- WeatherAPI.com free plan limits:  
  - 1 million API calls per month  

### 3) Load estimation:  

- Active users: 100,000  
- User subscription: 2-3 (total ~250,000 subscriptions)  
- Daily API Requests: ~150,000–300,000 (weather data fetch + email confirmation links)  
- Daily Email Sends: Up to 250,000 emails  

### 4) Components design  

#### 4.1) API Service  

- `GET /api/weather?city={city}` - Returns current temperature, humidity, and weather description for the specified city  
- `POST /api/subscribe` - Subscribes a user to weather updates by city and frequency; sends confirmation email.  
- `GET /api/confirm/{token}` - Confirms the subscription using the provided token.  
- `GET /api/unsubscribe/{token}` - Unsubscribes the user using a link from the update email.  

#### 4.2) SubscriptionService  

The SubscriptionService is responsible for managing user subscriptions to weather updates. It enables users to subscribe to weather updates for a specific city and frequency (daily/hourly), confirm their subscription with a link, and unsubscribe when desired.  

Responsibilities:  

- Create, confirm, and unsubscribe user subscription  
- Validate user email, city, frequency and ability of user to create subscription  
- Generate unique token via UUID  

#### 4.3) WeatherService  

The WeatherService is responsible for interacting with the external Weather API to retrieve current weather data for a given city.  

Responsibilities:  

- Fetch current weather data from an external weather API.  
- Handle and parse external API errors.  
- Parse and transform JSON responses into structured DTOs.  

#### 4.4) Mailer  

Mailer composes and sends emails for users.  

- Generate HTML emails for:  
  - Subscription confirmation  
  - Successful confirmation  
  - Weather updates  
- Emails sent using SMTP (Gmail) via the gomail.v2 library.  

### 5) Database Design  

**Subscription entity**  

| Field      | Type       | Constraints             |
|------------|------------|-------------------------|
| id         | serial     | Primary Key             |
| email      | varchar    | NOT NULL                |
| city       | varchar    | NOT NULL                |
| frequency  | varchar(10)| NOT NULL                |
| token      | varchar    | NOT NULL, UNIQUE        |
| confirmed  | bool       | NOT NULL DEFAULT false  |
| created_at | timestamp  | NOT NULL                |
| updated_at | timestamp  | NOT NULL                |
| deleted_at | timestamp  |                         |

### 6) Deployment  

The service can be deployed easily using `docker-compose.yml`.  

We choose docker because of : 

- Simple deployment – Start the entire system (DB, backend, any other services) with one command anywhere
- Lightweight and fast – Containers start quickly and use fewer resources than VMs

For starting system use this command 

```bash
docker-compose up --build -d
```

### 7) Future Enhancements  

- Custom Frequency: Support for advanced scheduling (e.g. every 3 hours, weekends only).  
- Conditional Alerts: Notify users based on conditions like rain, heatwave, or high humidity.  
- Add other types of weather update sending, for example SMS.  
- Move beyond email-token confirmation to JWT or OAuth-based login.  
- Add WebSockets for instant weather updates.  
- Store recent weather queries in Redis or DB to reduce API calls and latency.  
- Deploy app on AWS or other cloud service. 