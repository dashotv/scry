package search

import (
	"fmt"
	"time"
)

func timeToDateBucket(t time.Time) string {
	t = t.UTC()
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}
