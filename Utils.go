package main

import "io/ioutil"

// ReadFile reads in a file as a string
func ReadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}

// ToStringMap converts map[interface{}]interface{} to map[string]string.
func ToStringMap(old map[interface{}]interface{}) map[string]string {
	newMap := make(map[string]string)

	for k, v := range old {
		newMap[k.(string)] = v.(string)
	}

	return newMap
}
