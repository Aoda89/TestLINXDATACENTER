package main

import (
	"TestLINXDATACENTER/CSV"
	"TestLINXDATACENTER/JSON"
	"log"
	"os"
	"strings"
)

func checkType(flag string) (result, name string) {

	split := strings.Split(flag, ".")

	if split[1] == "json" {
		result := "json"
		name := flag
		return result, name

	} else if split[1] == "csv" {
		result := "csv"
		name := flag
		return result, name
	} else {
		panic("Ошибка формата")
	}
}

func main() {

	args := os.Args
	if len(args) == 1 {
		log.Fatal("Отсутствуют аргумены")
	}
	flag, name := checkType(args[1])

	switch flag {

	case "csv":
		CSV.ReadCSV(name)

	case "json":
		JSON.ReadJSON(name)

	}

}
