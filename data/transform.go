package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	current_time        = time.Now().Local().AddDate(0, 0, -1)
	five_days           = time.Now().Local().AddDate(0, 0, -5)
	fiveFormat   string = five_days.Format("01-02-2006")
	currentDay   string = current_time.Format("01-02-2006")
)

type DataJSON struct {
	Country   string
	Confirmed int32
	Deaths    int32
	Recovered int32
}

func restartCurrentDate() {
	currentDay = current_time.Format("01-02-2006")
}

func fileExist(general bool) bool {
	var errInfo error
	if general {
		_, errInfo = os.Stat(currentDay + "general.json")
	} else {
		_, errInfo = os.Stat(currentDay + ".json")
	}

	if os.IsNotExist(errInfo) {
		return false
	}

	return true
}

func Transform(json *[]byte, general bool) error {
	if general {
		jsonFile, err := os.Create(currentDay + "general.json")
		if err != nil {
			log.Println("Error to create json file")
			return err
		}

		defer jsonFile.Close()
		jsonFile.Write(*json)
		jsonFile.Close()
		return nil
	}

	jsonFile, err := os.Create(currentDay + ".json")
	if err != nil {
		log.Println("Error to create json file")
		return err
	}

	defer jsonFile.Close()
	jsonFile.Write(*json)
	jsonFile.Close()

	return nil
}

func validateExist(country string, allData *[]DataJSON) int {
	for key, val := range *allData {
		if val.Country == country {
			return key
		}
	}
	return -1
}

func loadJSON(name string, dataToJson *[]DataJSON) error { //pass pointer for load correct data
	jsonFile, err := os.Open(currentDay + name + ".json")

	if err != nil {
		log.Println("Error to load json file")
		return err
	}

	defer jsonFile.Close()

	bytesValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytesValue, &dataToJson)

	return nil
}

func filterByCountry(country string, jsonFilter *[]DataJSON) {
	var tmp []DataJSON
	for _, val := range *jsonFilter {
		if country == strings.ToLower(val.Country) {
			tmp := append(tmp, val)
			(*jsonFilter) = tmp
			break
		}
	}
}

func delete(name string) {
	_ = os.Remove(name)
}
