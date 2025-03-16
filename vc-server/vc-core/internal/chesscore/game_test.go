package chesscore

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)


func TestCustomVariantWithCheckmate(t *testing.T){
	runGameConfigTest(t,0,GameConfig{})
}

func runGameConfigTest(t *testing.T, configId int, expected GameConfig) {
	jsonFile, err := os.Open("test-configs.json")
	if err != nil {
		t.Errorf("Error opening json: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Errorf("Error reading json file: %v", err)
	}


	var gameConfigs []interface{}
	err = json.Unmarshal(byteValue, &gameConfigs)
	if err != nil {
		t.Errorf("Error unmarshalling json: %v", err)
	}
	if configId >len(gameConfigs){
		t.Errorf("Invalid configId")
	}
}

func compareGameConfigs(a, b GameConfig) bool {
	return reflect.DeepEqual(a, b)
}