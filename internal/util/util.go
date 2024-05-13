package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatalf("error marshaling into JSON: %v", err)
	}
	fmt.Println(string(b))
}

func PrettyFormat(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatalf("error marshaling into JSON: %v", err)
	}
	return string(b)
}
