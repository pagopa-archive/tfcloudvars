package tfcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	TF_CLOUD_URL = "https://app.terraform.io/api/v2/workspaces/%s/vars"
	BEARER_TOKEN = "Bearer %s"
	CONTENT_TYPE = "application/vnd.api+json"
)

type Attributes struct {
	Category    string `json:"category"`
	Created_at  string `json:"created-at,omitempty"`
	Description string `json:"description"`
	Hcl         bool   `json:"hcl"`
	Key         string `json:"key"`
	Sensitive   bool   `json:"sensitive"`
	Value       string `json:"value"`
}

type Links struct {
	Self string `json:"self"`
}

type ConfigData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ConfigLinks struct {
	Related string `json:"related"`
}

type Configurable struct {
	Data  ConfigData  `json:"data"`
	Links ConfigLinks `json:"links"`
}

type Relationships struct {
	Configurable `json:"configurable"`
}

type Data struct {
	Attributes    `json:"attributes"`
	ID            string `json:"id"`
	Links         `json:"links"`
	Relationships `json:"relationships"`
	Type          string `json:"type"`
}

type TerraformVars struct {
	Data []Data `json:"data"`
}

type Workspace struct {
	Data ConfigData `json:"data"`
}

type PayloadRelationships struct {
	Workspace Workspace `json:"workspace"`
}

type PayloadData struct {
	Attributes    `json:"attributes"`
	Relationships PayloadRelationships `json:"relationships"`
	Type          string               `json:"type"`
}

type Payload struct {
	Data PayloadData `json:"data"`
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func tojson(data interface{}, indent bool) (ret string, err error) {

	var byteArray []byte

	if indent == true {
		byteArray, err = json.MarshalIndent(data, "  ", "  ")
	} else {
		byteArray, err = json.Marshal(data)
	}

	ret = string(byteArray)

	return ret, err
}

func (v *TerraformVars) Json(indent bool) (data string, err error) {
	return tojson(v, indent)
}

func (v *Payload) Json(indent bool) (data string, err error) {
	return tojson(v, indent)
}

func (v *TerraformVars) TfVarsJson(indent bool) (data string, err error) {

	t := make(map[string]interface{})

	for _, v := range v.Data {

		if v.Sensitive == true {
			t[v.Key] = "sensitive"
		} else {
			t[v.Key] = v.Value
		}
	}

	var byteArray []byte

	if indent == true {
		byteArray, err = json.MarshalIndent(t, "  ", "  ")
	} else {
		byteArray, err = json.Marshal(t)
	}

	if err != nil {
		return "", fmt.Errorf("[ERROR] Marshalng map to json %t", err)
	}

	data = string(byteArray)

	return

}

func (v *TerraformVars) Get(w string, t string) (err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf(TF_CLOUD_URL, w), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf(BEARER_TOKEN, t))
	req.Header.Set("Content-Type", CONTENT_TYPE)

	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Http request status code %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err != nil {
		return err
	}

	if err := json.Unmarshal(respByte, &v); err != nil {
		return err
	}

	return nil
}

// Load TerraformVars from a json file
func (v *TerraformVars) Load(fileName string) (err error) {

	jsonFile, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	// fmt.Println(string(jsonFile))

	json.Unmarshal([]byte(jsonFile), &v)

	return nil
}

// Post the payload to the terraform cloud api that creates the variable.
// w is the workspace
// t is the bearer token
func (p *Payload) Post(w string, t string) (err error) {
	var resp *http.Response
	jsonPayload, err := p.Json(false)

	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(TF_CLOUD_URL, w),
		bytes.NewBuffer([]byte(jsonPayload)))

	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf(BEARER_TOKEN, t))
	req.Header.Set("Content-Type", CONTENT_TYPE)

	resp, err = Client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Http request status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	return nil

}

func (v *TerraformVars) Post(w string, t string) (err error) {

	for _, d := range v.Data {
		p := Payload{}
		p.Data.Type = "vars"
		p.Data.Attributes.Key = d.Attributes.Key
		p.Data.Attributes.Value = d.Attributes.Value
		p.Data.Attributes.Description = d.Attributes.Description
		p.Data.Attributes.Category = d.Attributes.Category
		p.Data.Attributes.Hcl = d.Attributes.Hcl
		p.Data.Attributes.Sensitive = d.Attributes.Sensitive
		p.Data.Relationships.Workspace.Data.ID = w
		p.Data.Relationships.Workspace.Data.Type = "workspaces"

		err := p.Post(w, t)

		if err != nil {
			log.Println(fmt.Sprintf("%s %s", d.Attributes.Key, err))
		} else {
			log.Println(fmt.Sprintf("%s created.", d.Attributes.Key))
		}

	}

	return nil
}
