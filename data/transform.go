package data

import (
	"log"
	"os"
)

type DataJSON struct {
	Country   string
	Confirmed int32
	Deaths    int32
	Recovered int32
}

type GeneralJSON struct {
	Country   string
	Confirmed int32
	Deaths    int32
	Recovered int32
}

func Transform(json *[]byte, date string) error {
	jsonFile, err := os.Create(date + ".json")
	if err != nil {
		log.Println("Error to create json file")
		return err
	}

	defer jsonFile.Close()
	jsonFile.Write(*json)
	jsonFile.Close()

	return nil
}

func TransformGeneral(json *[]byte, date string) error {
	jsonFile, err := os.Create("general.json")
	if err != nil {
		log.Println("Error to create json file")
		return err
	}

	defer jsonFile.Close()
	jsonFile.Write(*json)
	jsonFile.Close()

	return nil
}
