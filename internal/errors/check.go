package errors

import (
	"log"
	"log/syslog"
	"os"
	"sync"
)

var (
	// WaitForwards is a wait group which wait send all announcements before panic
	WaitForwards = new(sync.WaitGroup)

	sysLogger *syslog.Writer
)

// Check helps debug critical errors without warnings from 'gocyclo' linter
func Check(err error) {
	if err != nil {
		log.Println(err.Error())

		// Wait what all users get announcement message first
		WaitForwards.Wait()

		err = sysLogger.Crit(err.Error())
		if err != nil {
			log.Panicln(err.Error())
		}

		os.Exit(1)
	}
}
