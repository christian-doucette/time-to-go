package gtfs

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed stop-data/*.txt
var folder embed.FS

// reads through a CSV looking for a target value in at a specific column
func scanCsv(filePath string, targetFieldIndex int, targetValue string, returnFieldIndex int) (string, error) {
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
			return record[returnFieldIndex], nil
		}
	}

	return "", fmt.Errorf("Unable to find " + header[targetFieldIndex] + " " + targetValue + " in " + filePath)
}

// extracts the stop name from stops.txt file
func getStopName(stopId string) string {
	stopDataFilepaths := [7]string{
		"stop-data/bus-bronx.txt",
		"stop-data/bus-brooklyn.txt",
		"stop-data/bus-manhattan.txt",
		"stop-data/bus-mta-company.txt",
		"stop-data/bus-queens.txt",
		"stop-data/bus-staten-island.txt",
		"stop-data/subway.txt",
	}

	for _, stopDataFilepath := range stopDataFilepaths {
		stopName, err := scanCsv(stopDataFilepath, 0, stopId, 1)
		if err == nil {
			caser := cases.Title(language.English)
			return caser.String(stopName)
		}
	}

	panic("Could not find stopId " + stopId + " in any of the included files")
}

// formats time difference between two unix timestamps as string
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

// formats an arrivalEvent struct as a string
func (ae arrivalEvent) toString(currentTime int64) string {
	return fmt.Sprintf("%s | (%s) %s", formattedMinutesFromNow(currentTime, ae.expectedTime), ae.routeId, getStopName(ae.destinationStopId))
}

// formats a stopArrivalSnapshot as a list of strings
// numLines argument defines how many arrival times are included
func (sas stopArrivalSnapshot) toFormattedArrivalsList(numLines int) []string {
	currentTime := time.Now().Unix()
	formattedList := []string{}

	for _, arrivalEvent := range sas.arrivalEvents {
		if currentTime < arrivalEvent.expectedTime {
			formattedList = append(formattedList, arrivalEvent.toString(currentTime))
		}

		if len(formattedList) >= numLines {
			break
		}
	}

	return formattedList
}

// gets stop direction name
func getSubwayStopDirection(stopId string) string {
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

// formats subway title by combining stop direction and stop name into string
func subwayFormattedTitle(stopId string) string {
	return fmt.Sprintf("%s %s", getSubwayStopDirection(stopId), getStopName(stopId))
}

// formats bus title by using stop ID
func busFormattedTitle(stopId string) string {
	return getStopName(stopId)
}

func (sas stopArrivalSnapshot) busFormattedList(numLines int) []string {
	return append([]string{busFormattedTitle(sas.stopId)}, sas.toFormattedArrivalsList(numLines)...)
}

func (sas stopArrivalSnapshot) subwayFormattedList(numLines int) []string {
	return append([]string{subwayFormattedTitle(sas.stopId)}, sas.toFormattedArrivalsList(numLines)...)
}
