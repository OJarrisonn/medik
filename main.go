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
		if rules.CheckEnv(v) {
			fmt.Printf("%s is set\n", v)
		} else {
			fmt.Printf("%s is not set\n", v)
		}
	}
}
