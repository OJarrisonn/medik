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

	for _, v := range vars {
		if ok, _ := (&exams.EnvIsSet{EnvVar: v}).Examinate(); ok {
			fmt.Printf("%s is set\n", v)
		} else {
			fmt.Printf("%s is not set\n", v)
		}
	}
}
