package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

func main() {
	data, err := os.ReadFile("./medik.yaml")

	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	content := string(data)

	//fmt.Println(content)

	medik, err := config.Parse(content)

	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	fmt.Printf("%+v\n", medik)

	for _, v := range medik.Vitals {
		ty := v.Type
		fmt.Printf("type: %s\n", ty)
		parse, ok := exams.Parser(ty)

		if !ok {
			fmt.Printf("failed to find parser for type: %s", ty)
			continue
		}

		exam, err := parse(v)

		if err != nil {
			fmt.Printf("failed to parse exam: %v", err)
			continue
		}

		//fmt.Printf("%+v\n", exam)

		enforced, err := exam.Examinate()

		if err != nil {
			fmt.Printf("failed to examinate: %v", err)
			continue
		}

		fmt.Printf("enforced: %v\n", enforced)
	}
}
