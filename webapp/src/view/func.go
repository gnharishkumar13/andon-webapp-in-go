package view

import (
	"fmt"
	"html/template"
	"math"
	"time"
)

// DurationToHHMMSS converts a duration to a string with the format "HH:MM:SS ".
func DurationToHHMMSS(d time.Duration) string {
	h := d.Truncate(time.Hour).Hours()
	m := d.Truncate(time.Minute).Minutes() - h*60
	s := d.Truncate(time.Second).Seconds() - h*3600 - m*60
	return fmt.Sprintf("%02.f:%02.f:%02.f", h, math.Abs(m), math.Abs(s))
}

func registerFunctions(t *template.Template) {
	fm := map[string]interface{}{
		"durationToHHMMSS": DurationToHHMMSS,
	}
	t.Funcs(fm)
}
