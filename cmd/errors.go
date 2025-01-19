package cmd

import (
	"errors"
)

var (
	errCouldntUnmarshalToJSON  = errors.New("couldn't unmarshal data to JSON")
	errCouldntFetchDataFromDB  = errors.New("couldn't fetch data from hours' DB")
	errCouldntUpdateDataInDB   = errors.New("couldn't update data in hours' DB")
	errCouldntParseTaskID      = errors.New("couldn't parse the argument for task ID as an integer")
	errTaskAlreadyBeingTracked = errors.New("task is already being tracked")
	errCouldntParseBeginTS     = errors.New("couldn't parse begin timestamp")
	errCouldntParseEndTS       = errors.New("couldn't parse end timestamp")
)
