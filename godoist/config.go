package godoist

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
)

type TodoistConfig struct {
	ApiKey string `json:"apiKey"`
}

type Config struct {
	Todoist TodoistConfig `json:"todoist"`
}

var CONFIG Config

func LoadConfig() {
	usr, _ := user.Current()
	configFilePath := fmt.Sprintf("%s/.godoist", usr.HomeDir)
	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	configContents, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configContents, &CONFIG)
	if err != nil {
		panic(err)
	}
}
