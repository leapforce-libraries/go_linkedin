package linkedin

import (
	"time"

	"cloud.google.com/go/civil"
)

type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (tr *TimeRange) StartDateGMT(locationName string) civil.Date {
	location, _ := time.LoadLocation(locationName)
	time_ := time.Unix(0, tr.Start*1000000).In(location)
	return civil.DateOf(time_)
}

func (tr *TimeRange) EndDateGMT(locationName string) civil.Date {
	location, _ := time.LoadLocation(locationName)
	time_ := time.Unix(0, tr.End*1000000).In(location)
	return civil.DateOf(time_)
}
