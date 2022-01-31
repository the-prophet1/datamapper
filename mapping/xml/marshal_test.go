package xml

import (
	"encoding/json"
	"fmt"
	"testing"
)

var xmlText = `<?xml version="1.0" encoding="utf-8"?>
				<note> 
				<to>10</to>  
				<from>10</from>  
				<heading>Reminder</heading>  
				<body>Don't forget me this weekend!</body> 
				</note>`

func TestXMLUnmarshal(t *testing.T) {
	xmlMap := make(map[string]interface{})

	if err := Unmarshal([]byte(xmlText), &xmlMap); err != nil {
		t.Error(err)
	}

	if data, err := Marshal(xmlMap); err != nil {
		t.Error(err)
	} else {
		fmt.Println(data)
	}

	if data, err := json.Marshal(&xmlMap); err != nil {
		t.Error(err)
	} else {
		fmt.Println(data)
	}

}
