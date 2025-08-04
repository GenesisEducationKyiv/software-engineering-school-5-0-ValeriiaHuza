package weatherapi

type WeatherAPIErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  float64 `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}
