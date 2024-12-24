package main

import (
	"fmt"

	"github.com/OJarrisonn/medik/pkg/rules"
)

func main() {
	vars := []string{
		"HOME",
		"USER",
		"PATH",
	}

	for _, v := range vars {
		if ok, _ := (&rules.EnvIsSet{EnvVar: v}).Validate(); ok {
			fmt.Printf("%s is set\n", v)
		} else {
			fmt.Printf("%s is not set\n", v)
		}
	}
}
