package mta

import (
	"io"
	"net/http"
)

// checks if status code from response is OK
func ok(status int) bool {
	return 200 <= status && status < 300
}

// makes HTTP GET request using the given url and headers
func getRequest(url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	if !ok(resp.StatusCode) {
		panic("Request failed: " + resp.Status + " " + string(responseData))
	}

	return responseData
}

// calls all MTA API subway realtime feed endpoints and joins the results together
func CallAllSubwayRealtimeFeedApis() []byte {
	subwayRealtimeFeedUrls := [8]string{
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs",
		"https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-si",
	}

	combinedFeedResponses := []byte{}

	for _, realtimeFeedUrl := range subwayRealtimeFeedUrls {
		combinedFeedResponses = append(combinedFeedResponses, getRequest(realtimeFeedUrl)...)
	}

	return combinedFeedResponses
}

func CallBusRealtimeFeedApi() []byte {
	busRealtimeFeedUrl := "https://gtfsrt.prod.obanyc.com/tripUpdates"
	return getRequest(busRealtimeFeedUrl)
}
