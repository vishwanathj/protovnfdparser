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
