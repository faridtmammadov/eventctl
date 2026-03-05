package output

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(data []byte) {
	var obj map[string]interface{}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println(string(data))
		return
	}

	pretty, _ := json.MarshalIndent(obj, "", "  ")

	fmt.Println(string(pretty))
}
