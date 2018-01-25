package main

// errCheck helps debug critical errors without warnings from 'gocyclo' linter
func errCheck(err error) {
	if err != nil {
		waitForwards.Wait() // Wait what all users get announcement message
		panic(err.Error())
	}
}
