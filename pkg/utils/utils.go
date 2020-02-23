package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/vishwanathj/protovnfdparser/pkg/config"

	"github.com/ghodss/yaml"

	log "github.com/sirupsen/logrus"
	json_schema_val "github.com/vishwanathj/JSON-Parameterized-Data-Validator/pkg/jsondatavalidator"
)

const schemaStrPrefix = `{"$ref": "`
const schemaStrSuffix = `"}`

// schemaFileInputParam is name of schema file for input param files
const schemaFileInputParam = "inputParam.json"

// schemaFileDefineNonParam is name of schema file for non-parameterized templates.
// This is needed by the GenerateJSONSchemaFromParameterizedTemplate function
const schemaFileDefineNonParam = "vnfdDefineNonParam.json"

// schemaDir points to the relative path of where the schema files are located
var schemaDir string

// schemaStrVnfdInput is path to schema file for
// Parameterized templates
var schemaStrVnfdInput string

// schemaStrParameterizedInstance is path to schema file for
// instantiated Parameterized templates
var schemaStrParameterizedInstance string

// schemaStrPaginatedInstances is path to schema file
// for paginated output structure
var schemaStrPaginatedInstances string

func init() {
	schemaDir = config.GetConfigInstance().JsonSchemaConfig.SchemaDir
	schemaStrVnfdInput = schemaStrPrefix + schemaDir + "vnfdInputSchema.json#/vnfdInput" + schemaStrSuffix
	schemaStrParameterizedInstance = schemaStrPrefix + schemaDir + "vnfdInstanceSchema.json#/vnfdInstance" + schemaStrSuffix
	schemaStrPaginatedInstances = schemaStrPrefix + schemaDir + "vnfdPaginatedInstanceSchema.json#/vnfdsPaginatedInstances" + schemaStrSuffix
}

// ValidateVnfdPostBody validates the given JSON body against the parameterized
// VNFD Input JSON schema "parameterizedVnfdInputSchema.json" for compliance
func ValidateVnfdPostBody(body []byte) error {
	log.Debug()
	ioReaderObj := strings.NewReader(schemaStrVnfdInput)
	return json_schema_val.ValidateJSONBufAgainstSchema(body, ioReaderObj, "vnfdPostBody.json")
}

// ValidateVnfdInstanceBody validates the given JSON body against the parameterized
// VNFD Instance JSON schema "parameterizedVnfdInstanceSchema.json" for compliance
func ValidateVnfdInstanceBody(jsonval []byte) error {
	log.Debug()
	ioReaderObj := strings.NewReader(schemaStrParameterizedInstance)
	return json_schema_val.ValidateJSONBufAgainstSchema(jsonval, ioReaderObj, "vnfdInstanceBody.json")
}

// ValidatePaginatedVnfdsInstancesBody validates that JSON body returning the
// Vnfds adhere to the pagination format
func ValidatePaginatedVnfdsInstancesBody(jsonval []byte) error {
	log.Debug()
	ioReaderObj := strings.NewReader(schemaStrPaginatedInstances)
	return json_schema_val.ValidateJSONBufAgainstSchema(jsonval, ioReaderObj, "vnfdsPaginatedInstancesBody.json")
}

// ValidateInputParamAgainstParameterizedVnfd validates the given "input_param"
// JSON file against the dynamically generated JSON Schema
func ValidateInputParamAgainstParameterizedVnfd(inputParamJSON []byte,
	parameterizedVnfdJSON []byte) error {
	log.Debug()
	inputParamDynSchema, e := GenerateJSONSchemaFromParameterizedTemplate(parameterizedVnfdJSON)
	if e != nil {
		return e
	}
	data, e := yaml.YAMLToJSON(inputParamDynSchema)
	if e != nil {
		return e
	}
	fmt.Println(string(data))

	return json_schema_val.ValidateJSONBufAgainstSchema(inputParamJSON, strings.NewReader(string(data)), "inputParam.json")
}

// GenerateJSONSchemaFromParameterizedTemplate generated a dynamic schema
// by parsing the template for parameterized variables and looking up
// allowable values for those parameterized variables.
func GenerateJSONSchemaFromParameterizedTemplate(parameterizedJSON []byte) ([]byte, error) {
	abspath := schemaDir + schemaFileDefineNonParam
	nonParamDefineJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(abspath)
	if err != nil {
		return nil, err
	}

	abspath = schemaDir + schemaFileInputParam
	inputParamSchemaJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(abspath)
	if err != nil {
		return nil, err
	}

	return json_schema_val.GenerateJSONSchemaFromParameterizedTemplate(parameterizedJSON, nonParamDefineJSONBuf, inputParamSchemaJSONBuf, []string{"vnfd_id", "name"}, `\${1}(.*)`)
}

// GetSchemaDefinitionFileAsJSONBuf reads a Schema file and returns JSON buf
func GetSchemaDefinitionFileAsJSONBuf(schemaFileName string) ([]byte, error) {
	log.Debug()
	yamlText, err := ioutil.ReadFile(filepath.Clean(schemaFileName))

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(yamlText, &m)
	log.Debug(string(yamlText), err)

	return yamlText, err
}
