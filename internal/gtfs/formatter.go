package gtfs

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

func scanCsv(filePath string, targetFieldIndex int, targetValue string, returnFieldIndex int) string {
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to read input file " + filePath + ", error: " + err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	header, err := csvReader.Read()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if record[targetFieldIndex] == targetValue {
			return record[returnFieldIndex]
		}
	}

	panic("Unable to find " + header[targetFieldIndex] + " " + targetValue + " in " + filePath)
}

func getStopName(stopId string) string {
	return scanCsv("internal/gtfs/stops.txt", 0, stopId, 1)
}

func getStopDirection(stopId string) string {
	direction := stopId[len(stopId)-1]

	switch direction {
	case 'N':
		return "Uptown"
	case 'S':
		return "Downtown"
	default:
		return "Unknown Stop Direction"

	}
}

func formattedTitle(stopId string) string {
	return fmt.Sprintf("%s %s", getStopName(stopId), getStopDirection(stopId))
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

func (sas stopArrivalSnapshot) toString(numLines int) string {
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

	return strings.Join(formattedList, "\n")
}
