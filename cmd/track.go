package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	pers "github.com/dhth/hours/internal/persistence"
)

func renderActiveTLDetails(db *sql.DB, writer io.Writer) error {
	details, err := pers.FetchActiveTaskDetails(db)
	if errors.Is(err, pers.ErrNoTaskActive) {
		return err
	} else if err != nil {
		return fmt.Errorf("%w: %w", errCouldntFetchDataFromDB, err)
	}

	result, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntUnmarshalToJSON, err.Error())
	}

	fmt.Fprintln(writer, string(result))

	return nil
}

func startTracking(db *sql.DB, writer io.Writer, taskID int, comment *string) error {
	details, err := pers.FetchActiveTaskDetails(db)
	var noTaskActive bool
	if errors.Is(err, pers.ErrNoTaskActive) {
		noTaskActive = true
	} else if err != nil {
		return fmt.Errorf("%w: %w", errCouldntFetchDataFromDB, err)
	}

	now := time.Now()
	switch noTaskActive {
	case true:
		_, err := pers.InsertNewTL(db, taskID, now, comment)
		if err != nil {
			return fmt.Errorf("%w: %w", errCouldntUpdateDataInDB, err)
		}
	case false:
		if details.TaskID == taskID {
			return errTaskAlreadyBeingTracked
		}
		_, err := pers.QuickSwitchActiveTL(db, taskID, now, comment)
		if err != nil {
			return fmt.Errorf("%w: %w", errCouldntUpdateDataInDB, err)
		}
	}

	return renderActiveTLDetails(db, writer)
}

func updateTracking(db *sql.DB, writer io.Writer, beginTS time.Time, comment *string) error {
	_, err := pers.FetchActiveTaskDetails(db)
	if errors.Is(err, pers.ErrNoTaskActive) {
		return err
	} else if err != nil {
		return fmt.Errorf("%w: %w", errCouldntFetchDataFromDB, err)
	}

	err = pers.EditActiveTL(db, beginTS, comment)
	if err != nil {
		return fmt.Errorf("%w: %w", errCouldntUpdateDataInDB, err)
	}

	return renderActiveTLDetails(db, writer)
}

func stopTracking(db *sql.DB, writer io.Writer, beginTS, endTS *time.Time, comment *string) error {
	details, err := pers.FetchActiveTaskDetails(db)
	if errors.Is(err, pers.ErrNoTaskActive) {
		return err
	} else if err != nil {
		return fmt.Errorf("%w: %w", errCouldntFetchDataFromDB, err)
	}

	var bTS time.Time
	if beginTS == nil {
		bTS = details.CurrentLogBeginTS
	} else {
		bTS = *beginTS
	}

	var eTS time.Time
	if endTS == nil {
		eTS = time.Now()
	} else {
		eTS = *endTS
	}

	var commentToUse *string
	if comment == nil {
		commentToUse = details.CurrentLogComment
	} else {
		commentToUse = comment
	}

	err = pers.FinishActiveTL(db, details.TaskID, bTS, eTS, commentToUse)
	if err != nil {
		return fmt.Errorf("%w: %w", errCouldntUpdateDataInDB, err)
	}

	tlDetails, err := pers.FetchTLByID(db, details.TLID)
	if err != nil {
		return fmt.Errorf("%w: %w", errCouldntFetchDataFromDB, err)
	}

	result, err := json.MarshalIndent(tlDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntUnmarshalToJSON, err.Error())
	}

	fmt.Fprintln(writer, string(result))

	return nil
}
