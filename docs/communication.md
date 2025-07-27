# Microservice Decomposition and Communication

The current **Weather API** is implemented as a monolith. It provides multiple distinct features. For better **scalability**, **maintainability**, and **faster deployment cycles**, we plan to decompose the application into microservices.

---

## Miscroservice architecture 

### 1) API Gateway

- Acts as the main entry point to the Weather API.
- Routes incoming requests to the appropriate backend services.
- Provides the following public endpoints:

- `GET /api/weather?city={city}` - Returns weather for the specified city  
- `POST /api/subscribe` - Subscribes a user to weather updates
- `GET /api/confirm/{token}` - Confirms user subscription 
- `GET /api/unsubscribe/{token}` - Unsubscribes the user

--- 

### 2) Weather Service 

- Fetches current weather data from multiple providers using the **Chain of Responsibility** pattern.
- Caches weather responses in **Redis** for 10 minutes to reduce API load and improve response time.

---

### 3) Subscription & Scheduler service

- Handles user subscription logic: subscribe, confirm, unsubscribe.
- Persists data to PostgreSQL
- Schedules and triggers weather update emails at defined intervals (daily/hourly).
- Publishes email events to a **message queue** (e.g., "WeatherEmailToSend", "EmailConfirmed").

---

### 4) Email service 

- Sends confirmation and weather update emails using **SMTP**.
- Listens for and processes email-related events from the message queue.
- Can be extended to support templates, retries, and bounce handling.

---

## Communication between services

### Synchronous Communication (HTTP / gRPC)

Used when a direct response is required, such as:

- `API Gateway → Subscription Service` (on `POST /subscribe`)
- `API Gateway → Weather Service` (on `GET /weather`)

Note: For internal service-to-service communication, gRPC can be considered for better performance and stricter typing.

### Asynchronous Communication (Message Broker)

Used for non-blocking, event-driven interactions:

- `Subscription Service → Email Service`  
  Events:
  - `EmailConfirmation`
  - `EmailSubscribed`
  - `EmailUnsubscribed`
  - `WeatherEmailToSend`

#### Recommended Message Brokers

| Broker     | Use Case Strengths                                                                 |
|------------|--------------------------------------------------------------------------------------|
| **RabbitMQ** | Reliable delivery, flexible routing, retry/delay queues. Ideal for early production. |
| **NATS**      | Lightweight pub/sub with high throughput  |

**Best initial choice:** Start with **RabbitMQ** for reliability and built-in retry/delay queues.

---

## Reliability

Reliability is very important in a microservice system, espesially when we aim for 99.9% uptime. Here are some simple ways we make sure the system keeps working even when something goes wrong:

- **Retry with Backoff**: service communication will implement retry logic with delay.
- **Dead Letter Queues (DLQs)**: Messages that repeatedly fail in processing will be moved to DLQs for later inspection.
- **Health Checks** : Add /health endpoints to verify if services are running and ready.
- **Monitoring** : Add logs and alerts to quicly detect and resolve problems


---

## Implementation plan 

Microservice decomposition will proceed in multiple phases:

1. **Extract Email Service**
   - Isolate SMTP logic.
   - Add RabbitMQ as a message broker.
   - Handle all email operations via asynchronous events.

2. **Improve Reliability**
   - Add health checks for all services.
   - Implement retry logic

3. **Extract Weather Service**
   - Move external provider logic and Redis caching to a standalone service.
   - Expose endpoint: `GET /weather?city={city}`.

4. **Extract Subscription & Scheduler Service**
   - Migrate all subscription-related logic from monolith.
   - Implement scheduled jobs using **cron** or **robfig/cron** to trigger events.

5. **Add Centralized Observability**
   - Logging
   - Metrics with **Prometheus**  
   - Dashboards with **Grafana**

This decomposition will enable the weather platform to grow efficiently, respond to traffic surges gracefully, and evolve features in isolation without high coupling or risk.