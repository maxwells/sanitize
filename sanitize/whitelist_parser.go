package sanitize

import(
	"os"
	"encoding/json"
)

// Load a new whitelist from a JSON file
func NewWhitelist(filepath string) (*Whitelist, error) {
	m, err := mapFromJson(filepath)
	if err != nil {
		return nil, err
	}

	// initialize whitelist
	whitelist := &Whitelist{
		Elements: make([]*Element, len(m)),
	}

	// read each element (key) and attribute list (value)
	// from map m as Element and Attribute instances into the whitelist
	elementCounter := 0
	for k, v := range m {
		attributes := v.([]interface{})

   		element := &Element{
			Tag: k,
			Attributes: make([]*Attribute, len(attributes)),
		}

		for i, attributeName := range(attributes) {
			element.Attributes[i] = &Attribute{
				Name: attributeName.(string),
			}
		}

   		whitelist.Elements[elementCounter] = element
   		elementCounter += 1;
	}

	return whitelist, nil
}

// helper function to read entirety of JSON file
// into an arbitrary map
func mapFromJson(filepath string) (map[string]interface{}, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	// prepare byte array to read json file into
	fileInfo, err := f.Stat()
	bytes := make([]byte, fileInfo.Size())

	_, err = f.Read(bytes)
	if err != nil {
		return nil, err
	}

	// unmarshal json file into contract-free interface
	var arbitraryJson interface{}
	err = json.Unmarshal(bytes, &arbitraryJson)

	return arbitraryJson.(map[string]interface{}), err
}