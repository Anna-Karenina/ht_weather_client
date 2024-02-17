package app

import (
	geo_checker "go_education/weather/internal/geo_checker_api"
	open_meteo_api "go_education/weather/internal/open_meteo_api"
	tui "go_education/weather/internal/tui"
)

func RunApp() {
	app := tui.CreateApp()
	go app.Run()

	apiMessageChan := make(chan open_meteo_api.ForecastResponse)
	apiMessageErr := make(chan string)

	go func() {
		ipClientData, err := geo_checker.GetMyLocation()
		if err != nil {
			apiMessageErr <- "error"
			return
		}
		response, err := open_meteo_api.GetData(ipClientData.Lat, ipClientData.Long, ipClientData.IP)
		if err != nil {
			apiMessageErr <- "error"
			return
		}
		apiMessageChan <- response
	}()

	for {
		select {
		case <-apiMessageErr:
			{
				app.Send(tui.DataFetchedMessage{
					Action:  "error",
					Payload: open_meteo_api.ForecastResponse{}})
			}
		case msg := <-apiMessageChan:

			app.Send(
				tui.DataFetchedMessage{
					Action:  "dataFetched",
					Payload: msg,
				},
			)

		}
	}
}
