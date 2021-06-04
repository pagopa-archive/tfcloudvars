package tfcloud

import (
	"encoding/json"
	"io/ioutil"
)

type TfConfig struct {
	Credentials struct {
		App_terraform_io struct {
			Token string `json:"token"`
		} `json:"app.terraform.io"`
	} `json:"credentials"`
}

func (c *TfConfig) Read(fileName string) (err error) {
	jsonFile, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	json.Unmarshal([]byte(jsonFile), &c)

	return nil
}
