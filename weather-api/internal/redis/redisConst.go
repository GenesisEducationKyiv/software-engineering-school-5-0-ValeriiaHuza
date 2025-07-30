package redis

import "time"

const Delimeter = ":"
const WeatherKey = "weather" + Delimeter
const WeatherTTL = time.Minute * 15
