package mapping

import (
	"encoding/json"
	"testing"
)

var xmlText = `<?xml version="1.0" encoding="utf-8"?>
				<note> 
				<to>Tove</to>  
				<from>Jani</from>  
				<heading>Reminder</heading>  
				<body>Don't forget me this weekend!</body> 
				</note>`

func TestXMLUnmarshal(t *testing.T) {
	xmlMap := make(map[string]interface{})

	if err := XMLUnmarshal([]byte(xmlText), &xmlMap); err != nil {
		t.Error(err)
	}

	res, err := json.Marshal(&xmlMap)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}
