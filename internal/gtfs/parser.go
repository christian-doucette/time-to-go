package gtfs

import (
	"sort"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

type arrivalEvent struct {
	expectedTime      int64
	routeId           string
	destinationStopId string
}

type stopArrivalSnapshot struct {
	stopId        string
	arrivalEvents []arrivalEvent
}

// parses raw GTFS file to extract arrival times for a specific stop
func parse(stopId string, gtfsRaw []byte) stopArrivalSnapshot {
	feed := gtfs.FeedMessage{}
	err := proto.Unmarshal(gtfsRaw, &feed)
	if err != nil {
		panic(err.Error())
	}

	var arrivalEvents []arrivalEvent

	for _, entity := range feed.Entity {
		tripUpdate := entity.GetTripUpdate()

		stopTimeUpdates := tripUpdate.GetStopTimeUpdate()

		for _, stopTimeUpdate := range stopTimeUpdates {
			if stopId == stopTimeUpdate.GetStopId() {
				destinationStopId := stopTimeUpdates[len(stopTimeUpdates)-1].GetStopId()
				routeId := tripUpdate.GetTrip().GetRouteId()
				arrivalTime := stopTimeUpdate.GetArrival().GetTime()

				arrivalEvent := arrivalEvent{destinationStopId: destinationStopId, routeId: routeId, expectedTime: arrivalTime}
				arrivalEvents = append(arrivalEvents, arrivalEvent)
			}
		}
	}

	sort.Slice(arrivalEvents, func(i, j int) bool {
		return arrivalEvents[i].expectedTime < arrivalEvents[j].expectedTime
	})

	return stopArrivalSnapshot{stopId: stopId, arrivalEvents: arrivalEvents}
}

// extracts the next numLines arrival times for stopId from raw GTFS file
func ExtractSubwayStopArrivalTimes(gtfsRaw []byte, stopId string, numLines int, leastMinutesAhead int64) []string {
	return parse(stopId, gtfsRaw).subwayFormattedList(numLines, leastMinutesAhead)
}

// extracts the next numLines arrival times for stopId from raw GTFS file
func ExtractBusStopArrivalTimes(gtfsRaw []byte, stopId string, numLines int, leastMinutesAhead int64) []string {
	return parse(stopId, gtfsRaw).busFormattedList(numLines, leastMinutesAhead)
}
