package entity

import (
	"bufio"
	"github.com/asaskevich/govalidator"
	"log"
	"os"
)

var Params = make([]string, 0, 0)

func init() {
	inputFile, err := os.Open("/app/params")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatalf(scanner.Err().Error())
		}
		if govalidator.IsAlpha(scanner.Text()) && scanner.Text() != "" {
			Params = append(Params, scanner.Text())
		}
	}
}
