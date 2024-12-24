package main

import (
	"fmt"

	"github.com/OJarrisonn/medik/pkg/exams"
)

func main() {
	vars := []string{
		"HOME",
		"USER",
		"PATH",
	}

	if ok, err := (&exams.EnvIsSet{Vars: vars}).Examinate(); ok {
		fmt.Printf("%s is set\n", vars)
	} else {
		fmt.Printf("%s is not set\n%v", vars, err)
	}
}
