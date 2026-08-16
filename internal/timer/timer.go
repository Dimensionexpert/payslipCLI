package timer

import (
	"log"
	"time"
)

func Track(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %s", name, time.Since(start))
	}
}
