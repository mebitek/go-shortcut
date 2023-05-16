package utils

import (
	"encoding/json"
	"os"
)

func InitConfingDirectory() string {
	configDir := getUserDir() + "/.config/go-shortcuts"
	if !exists(configDir) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			return err.Error()
		}
	}
	dbFile := configDir + "/db.json"
	if !exists(dbFile) {
		create, err := os.Create(dbFile)
		if err != nil {
			panic(err.Error())
		}
		output, err := json.Marshal(make([]string, 0))
		err = os.WriteFile(dbFile, output, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
		return create.Name()
	}
	return dbFile

}

func getUserDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	return dirname
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
