package main

import (
	"fmt"

	"github.com/OJarrisonn/medik/pkg/check"
)

func main() {
	vars := []string{
		"HOME",
		"USER",
		"PATH",
	}

	for _, v := range vars {
		if check.CheckEnv(v) {
			fmt.Printf("%s is set\n", v)
		} else {
			fmt.Printf("%s is not set\n", v)
		}
	}
}
