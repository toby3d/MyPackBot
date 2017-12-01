package main

// errCheck helps debug critical errors without warnings from 'gocyclo' linter
func errCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
