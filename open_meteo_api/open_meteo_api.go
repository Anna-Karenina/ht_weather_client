package open_meteo_api

import (
	"encoding/json"
	"fmt"
	"go_education/weather/http_client"
	"io"
	"net/http"
	"time"
)

type Units struct {
	Time          string `json:"time"`
	Temperature   string `json:"temperature_2m"`
	Precipitation string `json:"precipitation"`
	WindSpeed     string `json:"wind_speed_10m"`
	WindDirection string `json:"wind_direction_10m"`
}

type HourlyTime struct {
	HourlyTimeTime []string  `json:"time"`
	Temperature    []float64 `json:"temperature_2m"`
	Precipitation  []float64 `json:"precipitation"`
	WindSpeed      []float64 `json:"wind_speed_10m"`
	WindDirection  []int     `json:"wind_direction_10m"`
}

type ForecastResponse struct {
	Elevation float32    `json:"elevation"`
	Units     Units      `json:"hourly_units"`
	Hourly    HourlyTime `json:"hourly"`
	TimeZone  string     `json:"timezone"`
	IpAdd     string
}

var openMetioClient *http_client.Client

func init() {
	baseUrl := "https://api.open-meteo.com/v1"
	openMetioClient = http_client.NewClient(

		&http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Println(req.Response.Status)
				fmt.Println("Redirect")
				return nil
			},
			// Transport: &LoggingTripper{
			// 	logger: os.Stdout,
			// 	next:   http.DefaultTransport,
			// },
			Timeout: time.Duration(time.Second * 3),
		},
		baseUrl,
	)
}

func GetData(lat string, long string, ip string) (ForecastResponse, error) {
	path := openMetioClient.
		Path(fmt.Sprintf("forecast?latitude=%s&longitude=%s&hourly=temperature_2m,precipitation,rain,showers,snowfall,wind_speed_10m,wind_direction_10m&start_date=2024-02-16&end_date=2024-02-16&timezone=auto", lat, long))
	res, err := openMetioClient.Get(path)
	if err != nil {
		return ForecastResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ForecastResponse{}, err
	}

	var r ForecastResponse
	if err = json.Unmarshal(body, &r); err != nil {
		return ForecastResponse{}, err
	}
	r.IpAdd = ip
	return r, nil

}
