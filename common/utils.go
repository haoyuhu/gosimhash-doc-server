package common

import (
	"time"
	"log"
	"io/ioutil"
	"encoding/json"
)

const UINT_BASE = 16

func CountBinary(number int) int {
	count := 0
	for ; number != 0; {
		number &= number - 1
		count ++
	}
	return count
}

func Check2Power(number int) bool {
	return number != 1 && CountBinary(number) == 1
}

func Now() int64 {
	// current time in seconds
	return time.Now().Unix()
}

func ReadJsonFile(filePath string) (map[string]interface{}, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Cannot read json file, details: %s\n", err.Error())
		return nil, err
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(raw, &ret)
	if err != nil {
		log.Printf("Cannot convert json file's content to map, details: %s\n", err.Error())
		return nil, err
	}
	return ret, nil
}
