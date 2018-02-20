package errors

import (
	"log"
	"log/syslog"
	"os"
	"sync"
)

var (
	WaitForwards = new(sync.WaitGroup)
	sysLogger    *syslog.Writer
)

// Check helps debug critical errors without warnings from 'gocyclo' linter
func Check(err error) {
	if err != nil {
		log.Println(err.Error())

		// Wait what all users get announcement message first
		WaitForwards.Wait()

		sysLogger.Crit(err.Error())
		os.Exit(1)
	}
}
