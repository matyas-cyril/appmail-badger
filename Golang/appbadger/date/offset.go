package date

import (
	"fmt"
	"time"
)

func OffSetUTC() string {

	_, offset := time.Now().Zone()

	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}

	return fmt.Sprintf("%s%02d:%02d", sign, offset/3600, (offset%3600)/60)
}
