
# Weather Subscription API

A simple Go-based web service that provides weather data and email subscription functionality. Users can subscribe to weather updates for specific cities and receive notifications via email.

---

## 🚀 Features

- Subscribe to weather updates by email
- Confirm/unsubscribe with email links
- Daily/Hourly email frequency
- Fetch current weather for a city
- HTML form for subscribing and weather lookup
- REST API built with Gin
- PostgreSQL database (using GORM)
- Scheduled jobs with `robfig/cron`
- SMTP email delivery

---

## 🛠 Tech Stack

- **Go** + Gin
- **PostgreSQL**
- **GORM** for ORM
- **Render** for deployment
- **OpenWeatherMap API** (for weather data)
- **Gomail** for sending emails
- **robfig/cron** for job scheduling

---

## 🌐 Live Demo

Deployed on: [https://weather-api-mt26.onrender.com](https://weather-api-mt26.onrender.com)

---

## 📦 API Endpoints


### Get weather for a city


```http
  GET /api/weather?city={city}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `city` | `string` | **Required**.   |


### Subscription

| Method | Endpoint                 | Description                   |
|--------|--------------------------|-------------------------------|
| POST   | `/api/subscribe`          | Subscribe to weather updates   |
| GET    | `/api/confirm/:token`     | Confirm a subscription         |
| GET    | `/api/unsubscribe/:token` | Unsubscribe from updates       |

---

## 📄 Environment Variables

Change a `.env` file if you want to run it locally (or set them in Render):

```env
APP_PORT=8000
DB_HOST=db_host
DB_PORT=5432
DB_USERNAME=your_db_user
DB_PASSWORD=your_db_pass
DB_NAME=your_db_name

MAIL_EMAIL=your_email@example.com
MAIL_PASSWORD=your_email_password

OPENWEATHER_API_KEY=your_openweathermap_key
```
