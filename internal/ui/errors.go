package ui

import (
	"errors"
)

var (
	errInteractiveModeNotApplicable = errors.New("interactive mode is not applicable")
	errCouldntAddDataToTable        = errors.New("couldn't add data to table")
	errCouldntRenderTable           = errors.New("couldn't render table")
)
