package mta

import (
	"io/ioutil"
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
		panic("Invalid Stop ID")
	}
}

func authorizationHeader(mtaApiKey string) (string, string) {
	return "x-api-key", mtaApiKey
}

func CallRealtimeFeedApi(mtaApiKey string, line byte) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", mtaUrl(line), nil)
	req.Header.Set(authorizationHeader(mtaApiKey))
	resp, err := client.Do(req)

	if err != nil {
		panic(err.Error())
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return responseData
}
