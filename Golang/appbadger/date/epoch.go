package date

import (
	"encoding/json"
	"os"
	"time"
)

func Epoch(value int64, utc bool) error {

	location := time.Local
	offset := OffSetUTC()

	if utc {
		location = time.UTC
		offset = "00:00"
	}

	if value < 0 {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Error  string
			Status bool
		}{
			Error:  "epoch value must be a positive integer",
			Status: false,
		})
	}

	var strDate string
	if utc {
		strDate = time.Unix(value, 0).UTC().Format("2006/01/02 15:04:05")
	} else { // Local
		strDate = time.Unix(value, 0).Local().Format("2006/01/02 15:04:05")
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Epoch    int64
		Date     string
		Location string
		Offset   string
	}{
		Epoch:    value,
		Date:     strDate,
		Location: location.String(),
		Offset:   offset,
	})
}
