package gtfs

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"time"
)

//go:embed stops.txt
var folder embed.FS

func scanCsv(filePath string, targetFieldIndex int, targetValue string, returnFieldIndex int) string {
	f, err := folder.Open(filePath)
	if err != nil {
		panic("Unable to read input file " + filePath + ", error: " + err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	header, err := csvReader.Read()
	if err != nil {
		panic(err.Error())
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		if record[targetFieldIndex] == targetValue {
			return record[returnFieldIndex]
		}
	}

	panic("Unable to find " + header[targetFieldIndex] + " " + targetValue + " in " + filePath)
}

func getStopName(stopId string) string {
	return scanCsv("stops.txt", 0, stopId, 1)
}

func getStopDirection(stopId string) string {
	direction := stopId[len(stopId)-1]

	switch direction {
	case 'N':
		return "Uptown"
	case 'S':
		return "Downtown"
	default:
		panic("Invalid StopID (must end with 'N' or 'S')")

	}
}

func formattedTitle(stopId string) string {
	return fmt.Sprintf("%s %s", getStopDirection(stopId), getStopName(stopId))
}

func formattedMinutesFromNow(currentTime int64, futureTime int64) string {

	minutesFromNow := float64(futureTime-currentTime) / 60.0

	if minutesFromNow < 1.0 {
		return "<1m"
	}

	roundedMinutesFromNow := int64(math.Round(minutesFromNow))

	if roundedMinutesFromNow < 10 {
		return fmt.Sprintf(" %dm", roundedMinutesFromNow)
	} else {
		return fmt.Sprintf("%dm", roundedMinutesFromNow)
	}
}

func (ae arrivalEvent) toString(currentTime int64) string {
	return fmt.Sprintf("%s | (%s) %s", formattedMinutesFromNow(currentTime, ae.expectedTime), ae.routeId, getStopName(ae.destinationStopId))
}

func (sas stopArrivalSnapshot) toFormattedList(numLines int) []string {
	currentTime := time.Now().Unix()
	formattedList := []string{formattedTitle(sas.stopId)}

	for _, arrivalEvent := range sas.arrivalEvents {
		if currentTime < arrivalEvent.expectedTime {
			formattedList = append(formattedList, arrivalEvent.toString(currentTime))
		}

		if len(formattedList) > numLines {
			break
		}
	}

	return formattedList
}
