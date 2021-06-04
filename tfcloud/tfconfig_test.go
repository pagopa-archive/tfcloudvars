package tfcloud

import (
	"fmt"
	"testing"
)

func TestValidConfig(t *testing.T) {

	c := TfConfig{}
	err := c.Read("./mocks/credentials.tfrc.json")

	if err != nil {
		t.Log(fmt.Printf("erroro reading config file %t", err))
		t.Fail()
	}

	expected := "this.is.a.test"
	actual := c.Credentials.App_terraform_io.Token

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}

func TestConfigFileNotFound(t *testing.T) {

	c := TfConfig{}
	err := c.Read("./mocks/notfound.tfrc.json")

	if err == nil {
		t.Log(fmt.Printf("error reading config file %t", err))
		t.Fail()
	}

}
