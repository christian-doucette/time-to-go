package mta

import (
	"io"
	"net/http"
)

// formats authorization header for MTA API request
func authorizationHeader(mtaApiKey string) (string, string) {
	return "x-api-key", mtaApiKey
}

// checks if status code from response is OK
func ok(status int) bool {
	return 200 <= status && status < 300
}

// calls a specific MTA api endpoint
func callMtaUrl(mtaApiKey string, url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set(authorizationHeader(mtaApiKey))
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	if !ok(resp.StatusCode) {
		panic("MTA API call failed: " + resp.Status + " " + string(responseData))
	}

	return responseData
}

// calls all MTA API realtime feed endpoints and joins the results together
func CallAllRealtimeFeedApis(mtaApiKey string) []byte {
	realtimeFeedUrls := [8]string{
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

	for _, realtimeFeedUrl := range realtimeFeedUrls {
		combinedFeedResponses = append(combinedFeedResponses, callMtaUrl(mtaApiKey, realtimeFeedUrl)...)
	}

	return combinedFeedResponses
}
