package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/parse"
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
		parse, ok := parse.GetExamParser(ty)

		if !ok {
			fmt.Printf("failed to find parser for type: %s\n", ty)
			continue
		}

		exam, err := parse(v)

		if err != nil {
			fmt.Printf("failed to parse exam: %v\n", err)
			continue
		}

		//fmt.Printf("%+v\n", exam)

		enforced, errs := exam.Examinate()

		if errs != nil {
			for _, e := range errs {
				fmt.Printf("error: %v\n", e)
			}
			continue
		}

		fmt.Printf("enforced: %v\n", enforced)
	}
}
