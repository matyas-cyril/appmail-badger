package date

import (
	pkgDefine "appbadger/define"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func Conv(day, hour string, utc bool) error {

	pkgDefine.TrimSpace(&day, &hour)

	// Vérifier la syntaxe de day
	if _, err := time.Parse("2006/01/02", day); err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Error  string
			Status bool
		}{
			Error:  err.Error(),
			Status: false,
		})
	}

	// Vérifier la syntaxe de heure
	if _, err := time.Parse("15:04:05", hour); err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Error  string
			Status bool
		}{
			Error:  err.Error(),
			Status: false,
		})
	}

	location := time.Local
	offset := OffSetUTC()

	if utc {
		location = time.UTC
		offset = "00:00"
	}

	date := fmt.Sprintf("%s %s", strings.TrimSpace(day), strings.TrimSpace(hour))

	t, err := time.ParseInLocation("2006/01/02 15:04:05", date, location)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Error  string
			Status bool
		}{
			Error:  err.Error(),
			Status: false,
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Epoch    int64
		Date     string
		Location string
		Offset   string
		Status   bool
	}{
		Epoch:    t.Unix(),
		Date:     date,
		Location: location.String(),
		Offset:   offset,
		Status:   true,
	})

}
