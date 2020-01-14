// +build unit

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package utils_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vishwanathj/protovnfdparser/pkg/utils"
)

var BASE_DIR = "../../test/testdata/yamlfiles/"

//var BASE_DIR_VALID_Parameterized_Input = "../../test/testdata/yamlfiles/valid/parameterizedInput/"
var BASE_DIR_VALID_Parameterized_Input = BASE_DIR + "valid/parameterizedInput/"
var BASE_DIR_VALID_Parameterized_Instance = BASE_DIR + "valid/parameterizedInstance/"
var BASE_DIR_INVALID_Parameterized_Input = BASE_DIR + "invalid/parameterizedInput/"
var BASE_DIR_INVALID_Parameterized_Instance = BASE_DIR + "invalid/parameterizedInstance/"
var BASE_DIR_VALID_NonParameterized_Input = BASE_DIR + "valid/nonParameterizedInput/"
var BASE_DIR_INVALID_NonParameterized_Input = BASE_DIR + "invalid/nonParameterizedInput/"
var BASE_DIR_VALID_Input_Param = BASE_DIR + "valid/inputParam/"
var BASE_DIR_INVALID_Input_Param = BASE_DIR + "invalid/inputParam/"
var BASE_DIR_VALID_Paginated = BASE_DIR + "valid/parameterizedPaginatedInstances/"
var BASE_DIR_INVALID_Paginated = BASE_DIR + "invalid/parameterizedPaginatedInstances/"

//var SchemaDir = "../schema/"
/*
func TestValidatePaginatedVnfdsInstancesBody_Positive(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Paginated)
	files, err := ioutil.ReadDir(bpath)
	log.WithFields(log.Fields{"Number of files": len(files), "path": bpath}).Info()

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidatePaginatedVnfdsInstancesBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

func TestValidatePaginatedVnfdsInstancesBody_Negative(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Paginated)
	files, err := ioutil.ReadDir(bpath)
	log.WithFields(log.Fields{"Number of files": len(files), "path": bpath}).Info()

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidatePaginatedVnfdsInstancesBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}

		}
	}
}
*/

func TestValidatePaginatedVnfdsInstancesBody(t *testing.T) {
	tables := []struct {
		description string
		baseDir     string
		fileType    string
		returnError bool // returnError holds false if err == nil and true if err != nil
	}{
		{"Positive: yaml file", BASE_DIR_VALID_Paginated, ".yaml", false},
		{"Positive: json file", BASE_DIR_VALID_Paginated, ".json", false},
		{"Negative: yaml file", BASE_DIR_INVALID_Paginated, ".yaml", true},
		{"Negative: json file", BASE_DIR_INVALID_Paginated, ".json", true},
	}
	for i, table := range tables {
		t.Run(fmt.Sprintf("%d:%s", i, table.description), func(t *testing.T) {
			files, err := ioutil.ReadDir(table.baseDir)

			if err != nil {
				t.Fatal(err)
			}

			for _, f := range files {
				if filepath.Ext(f.Name()) == table.fileType {
					yamlText, err := ioutil.ReadFile(table.baseDir + "/" + f.Name())
					if err != nil {
						t.Error("Error while reading VNFD File # ", err)
						// if file read fail the continue to next file.
						panic(err)
						t.Fail()
					}

					err = utils.ValidatePaginatedVnfdsInstancesBody(yamlText)

					if err == nil && table.returnError == false {
						t.Log("Success")
					} else if err != nil && table.returnError == true {
						t.Log(err)
					} else {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestValidateVnfdPostBody(t *testing.T) {
	tables := []struct {
		description string
		baseDir     string
		fileType    string
		returnError bool // returnError holds false if err == nil and true if err != nil
	}{
		{"Positive: parameterized input -> yaml files", BASE_DIR_VALID_Parameterized_Input, ".yaml", false},
		{"Positive: parameterized input -> json files", BASE_DIR_VALID_Parameterized_Input, ".json", false},
		{"Negative: parameterized input -> yaml files", BASE_DIR_INVALID_Parameterized_Input, ".yaml", true},
		{"Negative: parameterized input -> json files", BASE_DIR_INVALID_Parameterized_Input, ".json", true},
		{"Positive: non-parameterized input -> yaml files", BASE_DIR_VALID_NonParameterized_Input, ".yaml", false},
		{"Positive: non-parameterized input -> json files", BASE_DIR_VALID_NonParameterized_Input, ".json", false},
		{"Negative: non-parameterized input -> yaml files", BASE_DIR_INVALID_NonParameterized_Input, ".yaml", true},
		{"Negative: non-parameterized input -> json files", BASE_DIR_INVALID_NonParameterized_Input, ".json", true},
	}
	for i, table := range tables {
		t.Run(fmt.Sprintf("%d:%s", i, table.description), func(t *testing.T) {
			files, err := ioutil.ReadDir(table.baseDir)

			if err != nil {
				t.Fatal(err)
			}

			for _, f := range files {
				if filepath.Ext(f.Name()) == table.fileType {
					yamlText, err := ioutil.ReadFile(table.baseDir + "/" + f.Name())
					if err != nil {
						t.Error("Error while reading VNFD File # ", err)
						// if file read fail the continue to next file.
						panic(err)
						t.Fail()
					}

					err = utils.ValidateVnfdPostBody(yamlText)

					if err == nil && table.returnError == false {
						t.Log("Success")
					} else if err != nil && table.returnError == true {
						t.Log(err)
					} else {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestValidateVnfdInstanceBody(t *testing.T) {
	tables := []struct {
		description string
		baseDir     string
		fileType    string
		returnError bool // returnError holds false if err == nil and true if err != nil
	}{
		{"Positive: Valid parameterized instance -> yaml files", BASE_DIR_VALID_Parameterized_Instance, ".yaml", false},
		{"Positive: Valid parameterized instance -> json files", BASE_DIR_VALID_Parameterized_Instance, ".json", false},
		{"Negative: Invalid parameterized instance -> yaml files", BASE_DIR_INVALID_Parameterized_Instance, ".yaml", true},
		{"Negative: Invalid parameterized instance -> json files", BASE_DIR_INVALID_Parameterized_Instance, ".json", true},
	}

	for i, table := range tables {
		t.Run(fmt.Sprintf("%d:%s", i, table.description), func(t *testing.T) {
			files, err := ioutil.ReadDir(table.baseDir)

			if err != nil {
				t.Fatal(err)
			}

			for _, f := range files {
				if filepath.Ext(f.Name()) == table.fileType {
					yamlText, err := ioutil.ReadFile(table.baseDir + "/" + f.Name())
					if err != nil {
						t.Error("Error while reading VNFD File # ", err)
						// if file read fail the continue to next file.
						panic(err)
						t.Fail()
					}

					err = utils.ValidateVnfdInstanceBody(yamlText)

					if err == nil && table.returnError == false {
						t.Log("Success")
					} else if err != nil && table.returnError == true {
						t.Log(err)
					} else {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestValidateInputParamAgainstParameterizedVnfd(t *testing.T) {
	tables := []struct {
		description           string
		inputParamBaseDir     string
		inputParamFileName    string
		paramVnfdBaseDir      string
		parameterizedVnfdName string
		returnError           bool // returnError holds false if err == nil and true if err != nil
	}{
		{"Positive: inputParamConstraintsMissing", BASE_DIR_VALID_Input_Param,
			"inputParamConstraintsMissing.yaml",
			BASE_DIR_VALID_Parameterized_Input,
			"validParameterizedVNFDInputWithOptionalPropConstraintsMissing.yaml",
			false},
		{"Positive: inputParamHAMissing", BASE_DIR_VALID_Input_Param,
			"inputParamHAMissing.yaml",
			BASE_DIR_VALID_Parameterized_Input,
			"validParameterizedVNFDInputWithOptionalPropHAMissing.yaml",
			false},
		{"Positive: inputParamOptionalProps", BASE_DIR_VALID_Input_Param,
			"inputParamOptionalProps.yaml",
			BASE_DIR_VALID_Parameterized_Input,
			"validParameterizedVNFDInputWithOptionalProps.yaml",
			false},
		{"Positive: inputParamScaleMissing", BASE_DIR_VALID_Input_Param,
			"inputParamScaleMissing.yaml",
			BASE_DIR_VALID_Parameterized_Input,
			"validParameterizedVNFDInputWithOptionalPropScaleMissing.yaml",
			false},
		{"Positive: inputParamsRequiredProps", BASE_DIR_VALID_Input_Param,
			"inputParamsRequiredProps.yaml",
			BASE_DIR_VALID_Parameterized_Input,
			"validParameterizedVNFDInputWithRequiredProps.yaml",
			false},
		{"Positive: inputParamSubnetPools", BASE_DIR_VALID_Input_Param,
			"inputParamSubnetPools.yaml",
			BASE_DIR_VALID_Parameterized_Input,
			"validParameterizedVNFDInputWithOptionalProps.yaml",
			false},
		{"Negative: inputParamInvalidDedicatedConstraint", BASE_DIR_INVALID_Input_Param, "inputParamInvalidDedicatedConstraint.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidDiskSize", BASE_DIR_INVALID_Input_Param, "inputParamInvalidDiskSize.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidHA", BASE_DIR_INVALID_Input_Param, "inputParamInvalidHA.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidImageName", BASE_DIR_INVALID_Input_Param, "inputParamInvalidImageName.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidIPAddress", BASE_DIR_INVALID_Input_Param, "inputParamInvalidIPAddress.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidMemory", BASE_DIR_INVALID_Input_Param, "inputParamInvalidMemory.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidMinScale", BASE_DIR_INVALID_Input_Param, "inputParamInvalidMinScale.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidName", BASE_DIR_INVALID_Input_Param, "inputParamInvalidName.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidVcpus", BASE_DIR_INVALID_Input_Param, "inputParamInvalidVcpus.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidVimConstraint", BASE_DIR_INVALID_Input_Param, "inputParamInvalidVimConstraint.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
		{"Negative: inputParamInvalidVnfdIDFormat", BASE_DIR_INVALID_Input_Param, "inputParamInvalidVnfdIDFormat.yaml",
			BASE_DIR_VALID_Parameterized_Input, "validParameterizedVNFDInputWithOptionalProps.yaml", true},
	}

	for i, table := range tables {
		t.Run(fmt.Sprintf("%d:%s", i, table.description), func(t *testing.T) {
			vnfdpath := utils.GetAbsDIRPathGivenRelativePath(table.paramVnfdBaseDir)
			inputparampath := utils.GetAbsDIRPathGivenRelativePath(table.inputParamBaseDir)

			inparam, ierr := ioutil.ReadFile(inputparampath + "/" + table.inputParamFileName)

			if ierr != nil {
				t.Fatal(ierr)
			}

			vnfd, verr := ioutil.ReadFile(vnfdpath + "/" + table.parameterizedVnfdName)
			if verr != nil {
				t.Fatal(verr)
			}

			err := utils.ValidateInputParamAgainstParameterizedVnfd(inparam, vnfd)

			if table.returnError == true && err != nil {
				t.Log(err)
			} else if table.returnError == false && err == nil {
				//SUCCESS t.Log()
			} else {
				t.Fail()
			}
		})
	}
}

/*
func TestValidateVnfdPostBody_PositiveCases(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

func TestValidateVnfdPostBody_NegativeCases(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}

		}
	}
}


func TestValidNonParameterizedInputYaml(t *testing.T) {

	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_NonParameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

//KEEP THE COMMENTED CODE TILL JUNE 20th 2020
			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			schemaTextNonParameterizedInput := json_spec_val.GetSchemaStringWhenGivenFilePath(utils.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextNonParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				panic(err)
				fmt.Println(err)
				t.Fail()
			} else {
				t.Log("Passed")
			}
// END COMMENT
			err = utils.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}
		}
	}
}

func TestInValidNonParameterizedInputYaml(t *testing.T) {

	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_NonParameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

//KEEP THE COMMENTED CODE TILL JUNE 20th 2020
			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			schemaTextNonParameterizedInput := json_spec_val.GetSchemaStringWhenGivenFilePath(utils.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextNonParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}

			if err := schema.ValidateInterface(m); err != nil {
				t.Log(err)
			} else {
				t.Fail()
			}
//END COMMENT

			err = utils.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}
		}
	}
}

//Could not map in table driven test of as it seemed to be duplicate
func TestValidParameterizedInputYaml(t *testing.T) {

	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

//KEEP THE COMMENTED CODE TILL JUNE 20th 2020
			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			schemaTextParameterizedInput := json_spec_val.GetSchemaStringWhenGivenFilePath(utils.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				panic(err)
				fmt.Println(err)
				t.Fail()
			} else {
				t.Log("Passed")
			}
// END COMMENT

			err = utils.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

//Could not map in table driven test of as it seemed to be duplicate
func TestInValidParameterizedInputYaml(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath + BASE_DIR_INVALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}
// KEEP THE COMMENTED CODE TILL JUNE 20th 2020
			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			//if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInput)); err != nil {
			//if err := compiler.AddResource("schema.json", strings.NewReader(utilGetSchemaTextParameterizedInput())); err != nil {
			schemaTextParameterizedInput := json_spec_val.GetSchemaStringWhenGivenFilePath(utils.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				//panic(err)
				//fmt.Println(err)
				t.Logf("err")
			} else {
				t.Fail()
			}
// END COMMENT

			err = utils.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}
		}
	}
}
*/

/*func TestValidParameterizedInstanceYaml(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath + BASE_DIR_VALID_Parameterized_Instance + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidateVnfdInstanceBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}
//KEEP THE COMMENTED CODE TILL JUNE 20th 2020
			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			//if err := compiler.AddResource("schema.json", strings.NewReader(utilGetSchemaTextParameterizedInstance())); err != nil {
			//if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
			schemaTextParameterizedInstance := json_spec_val.GetSchemaStringWhenGivenFilePath(utils.SchemaParameterizedInstanceRelPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				panic(err)
				t.Fail()
			} else {
				t.Log("Passed")
			}
//END COMMENT
		}
	}
}

func TestInValidParameterizedInstanceYaml(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath + BASE_DIR_INVALID_Parameterized_Instance + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidateVnfdInstanceBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}
// START COMMENT
			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			//if err := compiler.AddResource("schema.json", strings.NewReader(utilGetSchemaTextParameterizedInstance())); err != nil {
			//if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
			schemaTextParameterizedInstance := json_spec_val.GetSchemaStringWhenGivenFilePath(utils.SchemaParameterizedInstanceRelPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				//panic(err)
				//fmt.Println(err)
				t.Logf("err")
			} else {
				t.Errorf("FAILED")
			}
//END COMMENT
		}
	}
}*/

/*
func TestValidateVnfdInstanceBody_PositiveCases(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidateVnfdInstanceBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

func TestValidateVnfdInstanceBody_NegativeCases(t *testing.T) {
	bpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = utils.ValidateVnfdInstanceBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}

		}
	}
}*/

/*
func TestPositive_ValidateInputParamAgainstParameterizedVnfd(t *testing.T) {
	tables := []struct {
		inputParamFileName    string
		parameterizedVnfdName string
	}{
		{"inputParamConstraintsMissing.yaml",
			"validParameterizedVNFDInputWithOptionalPropConstraintsMissing.yaml"},
		{"inputParamHAMissing.yaml",
			"validParameterizedVNFDInputWithOptionalPropHAMissing.yaml"},
		{"inputParamOptionalProps.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamScaleMissing.yaml",
			"validParameterizedVNFDInputWithOptionalPropScaleMissing.yaml"},
		{"inputParamsRequiredProps.yaml",
			"validParameterizedVNFDInputWithRequiredProps.yaml"},
		{"inputParamSubnetPools.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
	}

	vnfdpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	inputparampath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Input_Param)

	for _, table := range tables {
		inparam, ierr := ioutil.ReadFile(inputparampath + "/" + table.inputParamFileName)

		if ierr != nil {
			t.Fatal(ierr)
		}

		vnfd, verr := ioutil.ReadFile(vnfdpath + "/" + table.parameterizedVnfdName)
		if verr != nil {
			t.Fatal(verr)
		}

		err := utils.ValidateInputParamAgainstParameterizedVnfd(inparam, vnfd)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func TestNegative_ValidateInputParamAgainstParameterizedVnfd(t *testing.T) {
	tables := []struct {
		inputParamFileName    string
		parameterizedVnfdName string
	}{
		{"inputParamInvalidDedicatedConstraint.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidDiskSize.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidHA.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidImageName.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidIPAddress.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidMemory.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidMinScale.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidName.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidVcpus.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidVimConstraint.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidVnfdIDFormat.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
	}

	vnfdpath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	inputparampath := utils.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Input_Param)

	for _, table := range tables {
		inparam, ierr := ioutil.ReadFile(inputparampath + "/" + table.inputParamFileName)

		if ierr != nil {
			t.Fatal(ierr)
		}

		vnfd, verr := ioutil.ReadFile(vnfdpath + "/" + table.parameterizedVnfdName)
		if verr != nil {
			t.Fatal(verr)
		}

		err := utils.ValidateInputParamAgainstParameterizedVnfd(inparam, vnfd)
		if err != nil {
			t.Log(err)
		} else {
			t.Fail()
		}
	}
}
*/
