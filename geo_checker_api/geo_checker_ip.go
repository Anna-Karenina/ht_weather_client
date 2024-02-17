package geo_checker

import (
	"encoding/json"
	"fmt"
	"go_education/weather/http_client"
	"io"
	"net/http"
	"strings"
	"time"
)

type IpClientData struct {
	IP      string
	Long    string
	Lat     string
	City    string
	Country string
}

type ipConfigResponse struct {
	Ip      string `json:"ip"`
	City    string `json:"city"`
	Country string `json:"country"`
	Loc     string `json:"loc"`
}

var ipClient *http_client.Client

func init() {
	baseUrl := "https://ipinfo.io"
	ipClient = http_client.NewClient(

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

func GetMyLocation() (IpClientData, error) {
	res, err := ipClient.Get(ipClient.Path(""))
	if err != nil {
		return IpClientData{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return IpClientData{}, err
	}

	var r ipConfigResponse
	if err = json.Unmarshal(body, &r); err != nil {
		return IpClientData{}, err
	}

	loc := strings.Split(r.Loc, ",")
	return IpClientData{
		Lat:     loc[0],
		Long:    loc[1],
		City:    r.City,
		Country: r.Country,
		IP:      r.Ip,
	}, nil
}
