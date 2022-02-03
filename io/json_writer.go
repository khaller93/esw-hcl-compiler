package io

import (
    "encoding/json"
    "io/ioutil"
    "os"
)

// transforms the given data into JSON and writes it to the given file.
func WriteJSON(data interface{}, filepath string) error {
    if data != nil {
        json_data, err := json.Marshal(data)
        if err == nil {
            return ioutil.WriteFile(filepath, json_data, os.ModePerm)
        }
        return err
    }
    return ioutil.WriteFile(filepath, []byte(""), os.ModePerm)
}
