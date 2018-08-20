// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func WriteToFile(file, data string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}

func Clean(pathList ...string) error {
	for _, p := range pathList {
		err := os.RemoveAll(p)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
func cleanAll() error {
	return clean("Hosts.json", "Applications.json", "Hosts", "Applications")
}
*/
func DumpToFile(file string, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, json, 0644)
}