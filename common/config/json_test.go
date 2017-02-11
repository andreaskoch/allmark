// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"strings"
	"testing"
)

func Test_SerializeConfig_NoErrorIsReturned(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)

	config := &Config{
		Server: Server{
			ThemeFolderName: "/some/folder",
			HTTP: HTTP{
				Enabled: true,
			},
		},
	}

	serializer := JSONSerializer{}

	// act
	err := serializer.SerializeConfig(writeBuffer, config)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("The serialization of the config object return an error. %s", err)
	}
}

func Test_SerializeConfig_JSONContainsConfigValues(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)

	config := &Config{
		Server: Server{
			ThemeFolderName: "/some/folder",
			HTTP: HTTP{
				Enabled: true,
			},
		},
	}

	serializer := JSONSerializer{}

	// act
	serializer.SerializeConfig(writeBuffer, config)

	// assert
	json := writeBuffer.String()

	// assert: json contains theme folder
	if !strings.Contains(json, config.Server.ThemeFolderName) {
		t.Fail()
		t.Logf("The produced json does not contain the 'ThemeFolderName' value %q. The produced JSON is this: %s", config.Server.ThemeFolderName, json)
	}

	// assert: json contains http enabled
	if !strings.Contains(json, "true") {
		t.Fail()
		t.Logf("The produced json does not contain the 'Http Enabled' value %q. The produced JSON is this: %s", config.Server.HTTP.Enabled, json)
	}
}

func Test_SerializeConfig_JSONIsFormatted(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)

	config := &Config{
		Server: Server{
			ThemeFolderName: "/some/folder",
			HTTP: HTTP{
				Enabled: true,
			},
		},
	}

	serializer := JSONSerializer{}

	// act
	serializer.SerializeConfig(writeBuffer, config)

	// assert
	json := writeBuffer.String()

	// assert: json contains theme folder
	if !strings.Contains(json, "\n") {
		t.Fail()
		t.Logf("The produced json does not seem to be formatted. The produced JSON is this: %s", json)
	}
}

func Test_DeserializeConfig_EmptyObjectString_NoErrorIsReturned(t *testing.T) {
	// arrange
	json := `{}`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	_, err := serializer.DeserializeConfig(jsonReader)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("The deserialization of %q should not produce an error. But it did produce this error: %s", json, err)
	}
}

func Test_DeserializeConfig_FullConfigString_AllFieldsAreSet(t *testing.T) {
	// arrange
	json := `{
		"Server": {
			"ThemeFolderName": "/some/folder",
			"HTTP": {
				"Enabled": true
			}
		}
	}`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	config, _ := serializer.DeserializeConfig(jsonReader)

	// assert: Theme folder
	if config.Server.ThemeFolderName == "" {
		t.Fail()
		t.Logf("The deserialized config object should have the %q field properly initialized. Deserialization result: %#v", "ThemeFolderName", config)
	}

	// assert: http enabled
	if config.Server.HTTP.Enabled != true {
		t.Fail()
		t.Logf("The deserialized config object should have the %q field properly initialized. Deserialization result: %#v", "Http.Enabled", config)
	}
}

func Test_DeserializeConfig_ObjectWithDifferentFields_ConfigWithDefaultValuesIsReturned(t *testing.T) {
	// arrange
	json := `{
		"Name": "Ladi da",
		"AnotherField": {
		},
		"SomeList": [ "1", "2", "3" ]
	}
	`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	config, _ := serializer.DeserializeConfig(jsonReader)

	// assert
	emptyConfig := Config{}
	if config.Server.ThemeFolderName != emptyConfig.Server.ThemeFolderName {
		t.Fail()
		t.Logf("When the JSON cannot be mapped to the Config type the deserializer should return an uninitialized config object.")
	}
}

func Test_DeserializeConfig_EmptyString_ErrorIsReturned(t *testing.T) {
	// arrange
	json := ""
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	_, err := serializer.DeserializeConfig(jsonReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("DeserializeConfig should return an error if supplied JSON is invalid")
	}
}

func Test_DeserializeConfig_InvalidJSON_ErrorIsReturned(t *testing.T) {
	// arrange
	json := `dsajdklasdj/(/)(=7897402
		38748902
		;;;
	`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	_, err := serializer.DeserializeConfig(jsonReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("DeserializeConfig should return an error if supplied JSON is invalid")
	}
}
