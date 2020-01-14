package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"

	log "github.com/sirupsen/logrus"
	json_schema_val "github.com/vishwanathj/JSON-Parameterized-Data-Validator/pkg/jsondatavalidator"
)

// SchemaDir points to the relative path of where the schema files are located
var SchemaDir string

// SchemaInputPath is path to schema file for
// Parameterized templates
var SchemaInputPath string

// SchemaParameterizedInstanceRelPath is path to schema file for
// instantiated Parameterized templates
var SchemaParameterizedInstanceRelPath string

// SchemaPaginatedInstancesRelPath is path to schema file
// for paginated output structure
var SchemaPaginatedInstancesRelPath string

// SchemaFileInputParam is name of schema file for input param files
var SchemaFileInputParam string

// SchemaFileDefineNonParam is name of schema file for non-parameterized templates.
// This is needed by the GenerateJSONSchemaFromParameterizedTemplate function
var SchemaFileDefineNonParam string

func init() {
	log.Debug()
	localUnitTest := os.Getenv("TEST")
	log.Debug(localUnitTest)

	if localUnitTest == "true" {
		SchemaDir = "../schema/"
		SchemaInputPath = "../schema/vnfdInputSchema.json#/vnfdInput"
		SchemaParameterizedInstanceRelPath = "../schema/vnfdInstanceSchema.json#/vnfdInstance"
		SchemaPaginatedInstancesRelPath = "../schema/vnfdPaginatedInstanceSchema.json#/vnfdsPaginatedInstances"
		SchemaFileInputParam = "inputParam.json"
		SchemaFileDefineNonParam = "vnfdDefineNonParam.json"
	} else {
		SchemaDir = "/usr/share/vnfdservice/schema/"
		SchemaInputPath = "/usr/share/vnfdservice/schema/vnfdInputSchema.json#/vnfdInput"
		SchemaParameterizedInstanceRelPath = "/usr/share/vnfdservice/schema/vnfdInstanceSchema.json#/vnfdInstance"
		SchemaPaginatedInstancesRelPath = "/usr/share/vnfdservice/schema/vnfdPaginatedInstanceSchema.json#/vnfdsPaginatedInstances"
		SchemaFileInputParam = "inputParam.json"
		SchemaFileDefineNonParam = "vnfdDefineNonParam.json"
	}
}

// ValidateVnfdPostBody validates the given JSON body against the parameterized
// VNFD Input JSON schema "parameterizedVnfdInputSchema.json" for compliance
func ValidateVnfdPostBody(body []byte) error {
	log.Debug()
	var schemaText = GetSchemaStringWhenGivenFilePath(SchemaInputPath)
	ioReaderObj := strings.NewReader(schemaText)
	return json_schema_val.ValidateJSONBufAgainstSchema(body, ioReaderObj, "vnfdPostBody.json")
}

// ValidateVnfdInstanceBody validates the given JSON body against the parameterized
// VNFD Instance JSON schema "parameterizedVnfdInstanceSchema.json" for compliance
func ValidateVnfdInstanceBody(jsonval []byte) error {
	log.Debug()
	var schemaText = GetSchemaStringWhenGivenFilePath(SchemaParameterizedInstanceRelPath)
	ioReaderObj := strings.NewReader(schemaText)
	return json_schema_val.ValidateJSONBufAgainstSchema(jsonval, ioReaderObj, "vnfdInstanceBody.json")
}

// ValidatePaginatedVnfdsInstancesBody validates that JSON body returning the
// Vnfds adhere to the pagination format
func ValidatePaginatedVnfdsInstancesBody(jsonval []byte) error {
	log.Debug()
	var schemaText = GetSchemaStringWhenGivenFilePath(SchemaPaginatedInstancesRelPath)
	ioReaderObj := strings.NewReader(schemaText)
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
	abspath := GetAbsDIRPathGivenRelativePath(SchemaDir) + "/" + SchemaFileDefineNonParam
	nonParamDefineJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(abspath)
	if err != nil {
		return nil, err
	}

	abspath = GetAbsDIRPathGivenRelativePath(SchemaDir) + "/" + SchemaFileInputParam
	inputParamSchemaJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(abspath)
	//inputParamSchemaJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(SchemaFileInputParam)
	if err != nil {
		return nil, err
	}

	return json_schema_val.GenerateJSONSchemaFromParameterizedTemplate(parameterizedJSON, nonParamDefineJSONBuf, inputParamSchemaJSONBuf, []string{"vnfd_id", "name"}, `\${1}(.*)`)
	//return json_schema_val.GenerateJSONSchemaFromParameterizedTemplate(parameterizedJSON, nonParamDefineJSONBuf, inputParamSchemaJSONBuf, []string{"vnfd_id", "name"}, `.*\$.*`)
}

// GetAbsDIRPathGivenRelativePath returns the absolute path on the file system given the
// relative path from where this function resides
func GetAbsDIRPathGivenRelativePath(relpath string) string {
	log.Debug()
	_, fname, _, _ := runtime.Caller(0)
	var path string
	if strings.HasPrefix(relpath, "../") {
		path = filepath.Join(filepath.Dir(fname), relpath)
	} else {
		path = relpath
	}
	return path
}

// GetSchemaStringWhenGivenFilePath generates a string that needs to
// be passed to the schema validator method when compiling a json schema
func GetSchemaStringWhenGivenFilePath(relativePathOfJSONSchemaFile string) string {
	log.Debug()
	_, fname, _, _ := runtime.Caller(0)
	var path string
	if strings.HasPrefix(relativePathOfJSONSchemaFile, "../") {
		path = filepath.Join(filepath.Dir(fname), relativePathOfJSONSchemaFile)
	} else {
		path = relativePathOfJSONSchemaFile
	}

	var schemaText = `{"$ref": "` + path + `"}`
	log.Debug(schemaText)
	return schemaText
}

// GetSchemaDefinitionFileAsJSONBuf reads a Schema file and returns JSON buf
func GetSchemaDefinitionFileAsJSONBuf(schemaFileName string) ([]byte, error) {
	log.Debug()
	//bpath := GetAbsDIRPathGivenRelativePath(SchemaDir)
	//yamlText, err := ioutil.ReadFile(bpath + "/" + schemaFileName)
	yamlText, err := ioutil.ReadFile(schemaFileName)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(yamlText, &m)
	log.Debug(string(yamlText), err)

	return yamlText, err
}
