package configutils

import (
	"os"
	"testing"
)

type testConfig struct {
	ParamInt int
	ParamStr string
}

var config = testConfig{}

const (
	incorrectJsonFileName string = "testdata/incorrect.json"
	correctJsonFileName   string = "testdata/correct.json"

	envVarName
	envVarIncorrectJsonValue string = `{"" 	"ParamInt" : 10, 	"ParamStr" : "Value" }`
	envVarCorrectJsonValue   string = `{   "ParamInt" : 10,   "ParamStr" : "Value" }`
)

func TestLoadConfigFromNoneexistingFile(t *testing.T) {
	err := loadConfigFromFile("NoneExistingFileName", &config)
	if err == nil {
		t.Error("Should fail when file does not exist")
	}
}

func TestLoadConfigFromFileWoName(t *testing.T) {
	if err := loadConfigFromFile("", &config); err == nil {
		t.Errorf("Should fail when file name is an empty string")
	}

}

func TestloadConfigFromFileWithIncorrectJson(t *testing.T) {
	if err := loadConfigFromFile(incorrectJsonFileName, &config); err == nil {
		t.Errorf("Should fail when file contains incorrect json")
	}
}

func TestLoadConfigFromFileWithCorrectJson(t *testing.T) {
	if err := loadConfigFromFile(correctJsonFileName, &config); err != nil {
		t.Errorf("Config should be read successfully: %v", err)
	}
	if config.ParamInt != 10 {
		t.Errorf("paramInt should be equal to 10")
	}
	if config.ParamStr != "Value" {
		t.Errorf("paramStr should be equal to 'Value'")
	}
}

func TestLoadFromEmptyEnvVar(t *testing.T) {
	if err := loadConfigFromEnv("", &config); err == nil {
		t.Errorf("Should fail when env varialbe is an empty string")
	}
}

func TestLoadFromIncorrectJsonInEnvVar(t *testing.T) {
	os.Setenv(envVarName, envVarIncorrectJsonValue)
	if err := loadConfigFromEnv(envVarName, &config); err == nil {
		t.Errorf("Should fail when environment variable contains incorrect json")
	}
}

func TestLoadFromCorrectJsonInEnvVar(t *testing.T) {
	os.Setenv(envVarName, envVarCorrectJsonValue)
	if err := loadConfigFromEnv(envVarName, &config); err != nil {
		t.Errorf("Should NOT fail when env variable contains correct json: %v", err)
	}
	if config.ParamInt != 10 {
		t.Errorf("paramInt should be equal to 10")
	}
	if config.ParamStr != "Value" {
		t.Errorf("paramStr should be equal to 'Value'")
	}
}
