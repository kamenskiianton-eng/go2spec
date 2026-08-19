package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Professor struct {
	Name       string     `json:"name"`
	ScienceID  int        `json:"science_id"`
	IsWorking  bool       `json:"is_working"`
	University University `json:"university"`
}

type University struct {
	Name string `json:"name"`
	City string `json:"city"`
}

func main() {
	prof := Professor{
		Name:      "Ivan",
		ScienceID: 7333,
		IsWorking: true,
		University: University{
			Name: "MFUA",
			City: "Moscow",
		},
	}

	byteArr, err := json.MarshalIndent(prof, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(byteArr))
	err = os.WriteFile("output.json", byteArr, 0666) //-rw-rw-rw-
	if err != nil {
		log.Fatal(err)
	}
}
