package tfcloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/uolter/cptfcvars/tfcloud/mocks"
)

func TestGetValidResponse(t *testing.T) {

	Client = &mocks.MockClient{}

	j := `{
		"data": [
		  {
			"id": "var-yYXXXXXXXXX",
			"type": "vars",
			"attributes": {
			  "key": "cidr_subnet",
			  "value": "[\"10.0.5.0/24\"]",
			  "sensitive": false,
			  "category": "terraform",
			  "hcl": true,
			  "created-at": "2021-04-14T15:08:53.569Z",
			  "description": null
			},
			"relationships": {
			  "configurable": {
				"data": {
				  "id": "ws-gxxxxxxxxxx",
				  "type": "workspaces"
				},
				"links": {
				  "related": "/api/v2/organizations/test/workspaces/test"
				}
			  }
			},
			"links": {
			  "self": "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxxx"
			}
		  },
		  {
			"id": "var-xxxxxxxxxxxxxxx",
			"type": "vars",
			"attributes": {
			  "key": "database_name",
			  "value": "db",
			  "sensitive": false,
			  "category": "terraform",
			  "hcl": false,
			  "created-at": "2021-03-30T15:07:10.784Z",
			  "description": null
			},
			"relationships": {
			  "configurable": {
				"data": {
				  "id": "ws-xxxxxxxxxx",
				  "type": "workspaces"
				},
				"links": {
				  "related": "/api/v2/organizations/test/workspaces/test"
				}
			  }
			},
			"links": {
			  "self": "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxx"
			}
		  }
		]
	  }`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(j)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	v := TerraformVars{}

	err := v.Get("", "")

	if err != nil {
		t.Log(fmt.Printf("error %s ", err))
		t.Fail()
	}

	// Test num records
	expected := 2
	actual := len(v.Data)

	if expected != actual {
		t.Log(fmt.Printf("xerror expected %d actual %d", expected, actual))
		t.Fail()
	}

	// Test key
	expected_key := "cidr_subnet"
	actual_key := v.Data[0].Attributes.Key

	if expected_key != actual_key {
		t.Log(fmt.Printf("error expected %s actual %s", expected_key, actual_key))
		t.Fail()
	}

	// Test Value
	expected_val := "[\"10.0.5.0/24\"]"
	actual_val := v.Data[0].Attributes.Value

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected_val, actual_val))
		t.Fail()
	}
}

func TestGetEmptyResponse(t *testing.T) {

	Client = &mocks.MockClient{}

	json := `{
		"data": []
	  }`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	v := TerraformVars{}

	err := v.Get("", "")

	if err != nil {
		t.Log(fmt.Printf("error %s ", err))
		t.Fail()
	}

	// Test num records
	expected := string(rune(0))
	actual := string(rune(len(v.Data)))

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}

func TestGet404Response(t *testing.T) {

	Client = &mocks.MockClient{}

	json := `{
		"errors": [
		  {
			"status": "404",
			"title": "not found"
		  }
		]
	  }`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {

		return &http.Response{
			StatusCode: 404,
			Body:       r,
		}, nil
	}

	v := TerraformVars{}

	err := v.Get("", "")

	if err == nil {
		t.Log(fmt.Printf("error %s ", err))
		t.Fail()
	}

	// Test num records
	expected := string(rune(0))
	actual := string(rune(len(v.Data)))

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}

func TestTerraformVarsJsonOK(t *testing.T) {

	tf := TerraformVars{
		Data: []Data{
			{
				Attributes: Attributes{
					Category:    "category",
					Created_at:  "created-at",
					Description: "description",
					Hcl:         false,
					Key:         "key",
					Sensitive:   false,
					Value:       "value",
				},
				ID: "id",
				Links: Links{
					Self: "self",
				},
				Relationships: Relationships{
					Configurable: Configurable{
						Data: ConfigData{
							ID:   "id",
							Type: "type",
						},
						Links: ConfigLinks{
							Related: "related",
						},
					},
				},
				Type: "type",
			},
		},
	}

	actual, err := tf.Json(false)

	if err != nil {
		t.Log(fmt.Printf("error %s ", err))
		t.Fail()
	} else {
		expected := `{"data":[{"attributes":{"category":"category","created-at":"created-at","description":"description","hcl":false,"key":"key","sensitive":false,"value":"value"},"id":"id","links":{"self":"self"},"relationships":{"configurable":{"data":{"id":"id","type":"type"},"links":{"related":"related"}}},"type":"type"}]}`

		if expected != actual {
			t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
			t.Fail()
		}
	}

}

func TestPayloadJsonOK(t *testing.T) {

	p := Payload{
		Data: PayloadData{
			Attributes: Attributes{
				Category:    "category",
				Description: "description",
				Hcl:         false,
				Key:         "key",
				Sensitive:   false,
				Value:       "value",
			},
			Relationships: PayloadRelationships{
				Workspace: Workspace{
					Data: ConfigData{
						ID:   "id",
						Type: "type"},
				},
			},
			Type: "type"},
	}

	actual, err := p.Json(false)

	if err != nil {
		t.Log(fmt.Printf("error %s ", err))
		t.Fail()
	} else {
		expected := `{"data":{"attributes":{"category":"category","description":"description","hcl":false,"key":"key","sensitive":false,"value":"value"},"relationships":{"workspace":{"data":{"id":"id","type":"type"}}},"type":"type"}}`

		if expected != actual {
			t.Log(fmt.Printf("error expected %s \nactual %s", expected, actual))
			t.Fail()
		}
	}

}

func TestPostPayloadOK(t *testing.T) {

	j := `{
		"data": {
		  "id":"var-EavQ1LztoRTQHSNT",
		  "type":"vars",
		  "attributes": {
			"key":"some_key",
			"value":"some_value",
			"description":"some description",
			"sensitive":false,
			"category":"terraform",
			"hcl":false
		  },
		  "relationships": {
			"configurable": {
			  "data": {
				"id":"ws-4j8p6jX1w33MiDC7",
				"type":"workspaces"
			  },
			  "links": {
				"related":"/api/v2/organizations/my-organization/workspaces/my-workspace"
			  }
			}
		  },
		  "links": {
			"self":"/api/v2/workspaces/ws-4j8p6jX1w33MiDC7/vars/var-EavQ1LztoRTQHSNT"
		  }
		}
	  }`

	Client = &mocks.MockClient{}
	r := ioutil.NopCloser(bytes.NewReader([]byte(j)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {

		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	p := Payload{
		Data: PayloadData{
			Attributes: Attributes{
				Category:    "category",
				Description: "description",
				Hcl:         false,
				Key:         "key",
				Sensitive:   false,
				Value:       "value",
			},
			Relationships: PayloadRelationships{
				Workspace: Workspace{
					Data: ConfigData{
						ID:   "id",
						Type: "type"},
				},
			},
			Type: "type"},
	}

	err := p.Post("", "")

	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestPostKeyAlreadyExist(t *testing.T) {

	j := `{
		"errors": [
		  {
			"status": "422",
			"title": "invalid attribute",
			"detail": "Key has already been taken",
			"source": {
			  "pointer": "/data/attributes/key"
			}
		  }
		]
	  }`

	Client = &mocks.MockClient{}
	r := ioutil.NopCloser(bytes.NewReader([]byte(j)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {

		return &http.Response{
			StatusCode: 422,
			Body:       r,
		}, nil
	}

	p := Payload{
		Data: PayloadData{
			Attributes: Attributes{
				Category:    "category",
				Description: "description",
				Hcl:         false,
				Key:         "key",
				Sensitive:   false,
				Value:       "value",
			},
			Relationships: PayloadRelationships{
				Workspace: Workspace{
					Data: ConfigData{
						ID:   "id",
						Type: "type"},
				},
			},
			Type: "type"},
	}

	err := p.Post("", "")

	if err == nil {
		t.Log(err)
		t.Fail()
	}

}

func TestPostWorkspaceNotFound(t *testing.T) {

	j := `{
		"errors": [
		  {
			"status": "404",
			"title": "not found"
		  }
		]
	  }`

	Client = &mocks.MockClient{}
	r := ioutil.NopCloser(bytes.NewReader([]byte(j)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {

		return &http.Response{
			StatusCode: 404,
			Body:       r,
		}, nil
	}

	p := Payload{
		Data: PayloadData{
			Attributes: Attributes{
				Category:    "category",
				Description: "description",
				Hcl:         false,
				Key:         "key",
				Sensitive:   false,
				Value:       "value",
			},
			Relationships: PayloadRelationships{
				Workspace: Workspace{
					Data: ConfigData{
						ID:   "id",
						Type: "type"},
				},
			},
			Type: "type"},
	}

	err := p.Post("", "")

	if err == nil {
		t.Log(err)
		t.Fail()
	} else if err.Error() != "Http request status code 404" {
		t.Log("Expected error 404")
		t.Fail()
	}

}

func TestLoadJson(t *testing.T) {
	v := TerraformVars{}

	content := `{
		"data": [
		  {
		  "id": "var-yYXXXXXXXXX",
		  "type": "vars",
		  "attributes": {
			"key": "cidr_subnet",
			"value": "[\"10.0.5.0/24\"]",
			"sensitive": false,
			"category": "terraform",
			"hcl": true,
			"created-at": "2021-04-14T15:08:53.569Z",
			"description": null
		  },
		  "relationships": {
			"configurable": {
			"data": {
			  "id": "ws-gxxxxxxxxxx",
			  "type": "workspaces"
			},
			"links": {
			  "related": "/api/v2/organizations/test/workspaces/test"
			}
			}
		  },
		  "links": {
			"self": "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxxx"
		  }
		  },
		  {
		  "id": "var-xxxxxxxxxxxxxxx",
		  "type": "vars",
		  "attributes": {
			"key": "database_name",
			"value": "db",
			"sensitive": false,
			"category": "terraform",
			"hcl": false,
			"created-at": "2021-03-30T15:07:10.784Z",
			"description": null
		  },
		  "relationships": {
			"configurable": {
			"data": {
			  "id": "ws-xxxxxxxxxx",
			  "type": "workspaces"
			},
			"links": {
			  "related": "/api/v2/organizations/test/workspaces/test"
			}
			}
		  },
		  "links": {
			"self": "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxx"
		  }
		  }
		]
	  }`

	err := v.Load(content)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	expected := 2
	actual := len(v.Data)

	if expected != actual {
		t.Log(fmt.Printf("error expected %d actual %d", expected, actual))
		t.Fail()
	}

	expected_key := "cidr_subnet"
	actual_key := v.Data[0].Key

	if expected_key != actual_key {
		t.Log(fmt.Printf("error expected %s actual %s", expected_key, actual_key))
		t.Fail()
	}

}

func TestInvalidJson(t *testing.T) {
	v := TerraformVars{}

	err := v.Load("aaa")

	expected := err.Error()
	actual := "Invalid json invalid character 'a' looking for beginning of value"

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}

func TestToTfVarsJsonHcl(t *testing.T) {
	v := TerraformVars{}

	v.Data = append(v.Data, Data{ID: "var-a12eqwq", Type: "vars", Attributes: Attributes{
		Category:    "terraform",
		Created_at:  "2021-04-14T15:08:53.569Z",
		Description: "null",
		Hcl:         true,
		Key:         "cidr_subnet",
		Sensitive:   false,
		Value:       "[\"10.0.5.0/24\"]",
	}, Links: Links{
		Self: "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxxx",
	}})

	expected := "{\"cidr_subnet\":[\"10.0.5.0/24\"]}"
	actual, _ := v.ToTfVars(false)

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}

func TestToTfVarsJsonStr(t *testing.T) {
	v := TerraformVars{}

	v.Data = append(v.Data, Data{ID: "var-a12eqwq", Type: "vars", Attributes: Attributes{
		Category:    "terraform",
		Created_at:  "2021-04-14T15:08:53.569Z",
		Description: "null",
		Hcl:         false,
		Key:         "name",
		Sensitive:   false,
		Value:       "api",
	}, Links: Links{
		Self: "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxxx",
	}})

	expected := "{\"name\":\"api\"}"
	actual, _ := v.ToTfVars(false)

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}

func TestToTfVarsJsonNumber(t *testing.T) {
	v := TerraformVars{}

	v.Data = append(v.Data, Data{ID: "var-a12eqwq", Type: "vars", Attributes: Attributes{
		Category:    "terraform",
		Created_at:  "2021-04-14T15:08:53.569Z",
		Description: "null",
		Hcl:         false,
		Key:         "count",
		Sensitive:   false,
		Value:       "10",
	}, Links: Links{
		Self: "/api/v2/workspaces/ws-xxxxxx/vars/var-xxxxxx",
	}})

	expected := "{\"count\":10}"
	actual, _ := v.ToTfVars(false)

	if expected != actual {
		t.Log(fmt.Printf("error expected %s actual %s", expected, actual))
		t.Fail()
	}
}
