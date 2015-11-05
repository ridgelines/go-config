package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"strconv"
)

type JSONFile struct {
	*FileProvider
}

func NewJSONFile(path string) *JSONFile {
	return &JSONFile{
		FileProvider: NewFileProvider(path),
	}
}

func (this *JSONFile) Load() error {
	encodedJSON, err := ioutil.ReadFile(this.Path)
	if err != nil {
		return err
	}

	decodedJSON := map[string]interface{}{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		return err
	}

	tokens, err := this.flatten(decodedJSON, "")
	if err != nil {
		return err
	}

	this.tokens = tokens
	return nil

}

func (this *JSONFile) GetTokens() map[string]string {
	return this.tokens
}

func (this *JSONFile) flatten(inputJSON map[string]interface{}, namespace string) (map[string]string, error) {
	flattened := map[string]string{}

	for key, value := range inputJSON {
		var token string
		if namespace == "" {
			token = key
		} else {
			token = fmt.Sprintf("%s.%s", namespace, key)
		}

		if child, ok := value.(map[string]interface{}); ok {
			tokens, err := this.flatten(child, token)
			if err != nil {
				return nil, err
			}

			for k, v := range tokens {
				flattened[k] = v
			}
		} else {
			flattened[token] = fmt.Sprintf("%v", value)
		}
	}

	return flattened, nil
}

/*
func flatten(inputJSON map[string]interface{}, lkey string, flattened *map[string]interface{}) {
	for rkey, value := range inputJSON {
		key := lkey + rkey

		fmt.Println("key: ", key, value)
	}


			if _, ok := value.(string); ok {
				(*flattened)[key] = value.(string)
			} else if _, ok := value.(float64); ok {
				strconv.ParseInt((value.(float64)))
				(*flattened)[key] = strconv.ParseInt((value.(float64)))
			} else if _, ok := value.(bool); ok {
				(*flattened)[key] = strconv.ParseBool(value.(bool))
			} else if _, ok := value.([]interface{}); ok {
				for i := 0; i < len(value.([]interface{})); i++ {
					if _, ok := value.([]string); ok {
						stringI := string(i)
						(*flattened)[stringI] = value.(string)
						/// think this is wrong
					} else if _, ok := value.([]int); ok {
						stringI := string(i)
						(*flattened)[stringI] = value.(int)
					} else {
						flatten(value.([]interface{})[i].(map[string]interface{}), key+":"+strconv.Itoa(i)+":", flattened)
					}
				}
			} else {
				flatten(value.(map[string]interface{}), key+".", flattened)
			}
		}

}
/*

/*

func load(filename string) map[string]interface{} {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	mappedJSON := decodeJSON(fileContents)

	return mappedJSON
}

func decodeJSON(encodedJSON []byte) map[string]interface{} {
	decoded := map[string]interface{}{}
	err := json.Unmarshal(encodedJSON, &decoded)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

func main() {
	mappedJSON := load("global.json")

	flatten(mappedJSON, lkey, &flattened)
	for key, value := range flattened {
		fmt.Printf("%v:%v\n", key, value)
	}
}
*/
