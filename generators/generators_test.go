package generators

import (
	"testing"
	"time"
)

func TestSeconds(t *testing.T) {
	t.SkipNow()
	i := 0
	for now := range Seconds {
		t.Logf("Seconds at %s", now.Format(time.RFC3339))
		if i > 5 {
			break
		}
		i++
	}
}
