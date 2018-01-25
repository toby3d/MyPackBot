package main

import "fmt"

// errCheck helps debug critical errors without warnings from 'gocyclo' linter
func errCheck(err error) {
	if err != nil {
		fmt.Sprintln(err.Error())
		waitForwards.Wait() // Wait what all users get announcement message
		panic(err.Error())
	}
}
