package mta

import (
	"io"
	"net/http"
)

func mtaUrl(line byte) string {
	switch line {
	case 'A', 'C', 'E':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace"
	case 'B', 'D', 'F', 'M':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm"
	case 'G':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g"
	case 'J', 'Z':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz"
	case 'N', 'Q', 'R', 'W':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw"
	case 'L':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l"
	case '1', '2', '3', '4', '5', '6', '7':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"
	case 'S':
		return "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-si"
	default:
		panic("Invalid StopID (first letter must be one of 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'J', 'L', 'M', 'N', 'Q', 'R', 'S', 'W', 'Z', '1', '2', '3', '4', '5', '6', '7')")
	}
}

func authorizationHeader(mtaApiKey string) (string, string) {
	return "x-api-key", mtaApiKey
}

func ok(status int) bool {
	return 200 <= status && status < 300
}

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
