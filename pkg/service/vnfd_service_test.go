// +build unit

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vishwanathj/protovnfdparser/pkg/config"
	dalmocks "github.com/vishwanathj/protovnfdparser/pkg/datarepo/mocks"
	vsvcerr "github.com/vishwanathj/protovnfdparser/pkg/errors"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	umocks "github.com/vishwanathj/protovnfdparser/pkg/utils/mocks"
)

var createVnfdObjVcpusMissing = []byte(`{
  "virtual_links": [
    {
      "name": "mgmt_net",
      "is_management": true
    }
  ],
  "name": "vnfname",
  "vdus": [
    {
      "disk_size": "$disk_size",
      "name": "vdu1",
      "memory": "$memory",
      "vnfcs": [
        {
          "connection_points": [
            {
              "virtualLinkReference": [
                "mgmt_net"
              ],
              "ip_address": "$vdu1_vnfc1_mgmt",
              "name": "mgmtCP"
            }
          ],
          "name": "activeF5"
        }
      ],
      "image": "$image"
    }
  ]
}`)

var createVnfdObj = []byte(`{
  "virtual_links": [
    {
      "name": "mgmt_net",
      "is_management": true
    }
  ],
  "name": "vnfname",
  "vdus": [
    {
      "vcpus": "$vcpus",
      "disk_size": "$disk_size",
      "name": "vdu1",
      "memory": "$memory",
      "vnfcs": [
        {
          "connection_points": [
            {
              "virtualLinkReference": [
                "mgmt_net"
              ],
              "ip_address": "$vdu1_vnfc1_mgmt",
              "name": "mgmtCP"
            }
          ],
          "name": "activeF5"
        }
      ],
      "image": "$image"
    }
  ]
}`)

// Make sure to copy the entire JSON from an already existing VNFD in mongodb for the below variable.
// Missing some fields would cause the test to file.
var validInstanceBody = []byte(`{
    "id" : "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638451",
    "name" : "vnfdOptProps",
    "status" : "available",
    "created_at" : "2019-02-10T17:00:35-06:00",
    "vdus" : [ 
        {
            "constraints" : {
                "dedicated" : "$dedval",
                "vim_id" : "$vimval"
            },
            "disk_size" : "$disk_size",
            "high_availability" : "$haval",
            "image" : "$image",
            "memory" : "$memory",
            "name" : "vdu1",
            "scale_in_out" : {
                "default" : "$def",
                "maximum" : "$max",
                "minimum" : "$min"
            },
            "vcpus" : "$vcpus",
            "vnfcs" : [ 
                {
                    "connection_points" : [ 
                        {
                            "ip_address" : "$vdu1_vnfc1_mgmt",
                            "name" : "mgmtCP",
                            "virtuallinkreference" : [ 
                                "mgmt_net"
                            ]
                        }, 
                        {
                            "ip_address" : "$vdu1_work_net",
                            "name" : "internalCP",
                            "virtuallinkreference" : [ 
                                "worknet"
                            ]
                        }
                    ],
                    "name" : "activeF5"
                }, 
                {
                    "connection_points" : [ 
                        {
                            "ip_address" : "$vdu1_vnfc1_mgmt",
                            "name" : "mgmtCP",
                            "virtuallinkreference" : [ 
                                "mgmt_net"
                            ]
                        }, 
                        {
                            "ip_address" : "$vdu1_work_net",
                            "name" : "internalCP",
                            "virtuallinkreference" : [ 
                                "worknet"
                            ]
                        }
                    ],
                    "name" : "passiveF5"
                }
            ]
        }
    ],
    "virtual_links" : [ 
        {
            "name" : "worknet",
            "is_management" : false
        }, 
        {
            "name" : "mgmt_net",
            "is_management" : true
        }
    ]
}`)

var mockdal = dalmocks.VnfdRepository{}
var testcfg = config.GetConfigInstance()

type onCallReturnArgs struct {
	onCallMethodName    string
	onCallMethodArgType string
	retArgList          []interface{}
}

func TestVnfdService_CreateVnfd(t *testing.T) {
	mockErr := errors.New("mock error")

	// svcInstance is defined in vnfd_service.go
	svcInstance = &VnfdService{&mockdal, *testcfg}

	tests := []struct {
		desc            string
		postBody        []byte
		expectedErr     vsvcerr.VnfdsvcError
		daOnReturnArgs  *onCallReturnArgs
		valOnReturnArgs []onCallReturnArgs
	}{
		{
			desc:        "vcpu missing, post body non-compliant with schema",
			postBody:    createVnfdObjVcpusMissing,
			expectedErr: vsvcerr.VnfdsvcError{mockErr, http.StatusBadRequest},
		},
		{
			desc:        "vnfd instance body non-compliant with schema",
			postBody:    createVnfdObj,
			expectedErr: vsvcerr.VnfdsvcError{mockErr, http.StatusBadRequest},
			valOnReturnArgs: []onCallReturnArgs{
				{"ValidateVnfdPostBody", "[]uint8", []interface{}{nil}},
				{"ValidateVnfdInstanceBody", "[]uint8", []interface{}{mockErr}},
			},
		},
		{
			desc:           "insert into DB success",
			postBody:       createVnfdObj,
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusOK},
			daOnReturnArgs: &onCallReturnArgs{"InsertVnfd", "*models.Vnfd", []interface{}{nil}},
			valOnReturnArgs: []onCallReturnArgs{
				{"ValidateVnfdPostBody", "[]uint8", []interface{}{nil}},
				{"ValidateVnfdInstanceBody", "[]uint8", []interface{}{nil}},
			},
		},
		{
			desc:           "insert into DB failure",
			postBody:       createVnfdObj,
			expectedErr:    vsvcerr.VnfdsvcError{mockErr, http.StatusConflict},
			daOnReturnArgs: &onCallReturnArgs{"InsertVnfd", "*models.Vnfd", []interface{}{mockErr}},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			var vnfd models.Vnfd
			err := json.Unmarshal(tc.postBody, &vnfd)
			if err != nil {
				t.Fatal("Failed to unmarshal")
			}

			lenValArgs := len(tc.valOnReturnArgs)
			if lenValArgs != 0 {
				// vnfdVerifier is defined in vnfd_service.go
				vnfdVerifier = &umocks.VnfdValidator{}
				valMock := vnfdVerifier.(*umocks.VnfdValidator)
				for _, item := range tc.valOnReturnArgs {
					call := valMock.On(item.onCallMethodName, mock.AnythingOfType(item.onCallMethodArgType))
					for _, elem := range item.retArgList {
						call.ReturnArguments = append(call.ReturnArguments, elem)
					}
				}
			}

			if tc.daOnReturnArgs != nil {
				daMock := svcInstance.dal.(*dalmocks.VnfdRepository)
				call := daMock.On(tc.daOnReturnArgs.onCallMethodName, mock.AnythingOfType(tc.daOnReturnArgs.onCallMethodArgType))
				for _, elem := range tc.daOnReturnArgs.retArgList {
					call.ReturnArguments = append(call.ReturnArguments, elem)
				}
				call.Once()
			}

			svcerr := svcInstance.CreateVnfd(&vnfd)
			fmt.Println(svcerr.HttpCode, tc.expectedErr.HttpCode)
			assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
		})
	}
}

func TestVnfdService_GetVnfd(t *testing.T) {
	// setup test
	mockErr := errors.New("mock error")
	var vnfd models.Vnfd
	err := json.Unmarshal(validInstanceBody, &vnfd)
	if err != nil {
		t.Fatal("Failed to unmarshal")
	}
	// svcInstance is defined in vnfd_service.go
	svcInstance = &VnfdService{&mockdal, *testcfg}

	// vnfdVerifier is defined in vnfd_service.go
	vnfdVerifier = &umocks.VnfdValidator{}

	daMock := svcInstance.dal.(*dalmocks.VnfdRepository)
	valMock := vnfdVerifier.(*umocks.VnfdValidator)

	tests := []struct {
		desc            string
		nameorid        string
		expectedErr     vsvcerr.VnfdsvcError
		daOnReturnArgs  *onCallReturnArgs
		valOnReturnArgs *onCallReturnArgs
	}{
		{
			desc:            "Get by ID pass",
			nameorid:        "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638451",
			expectedErr:     vsvcerr.VnfdsvcError{nil, http.StatusOK},
			daOnReturnArgs:  &onCallReturnArgs{"FindVnfdByID", "string", []interface{}{&vnfd, nil}},
			valOnReturnArgs: &onCallReturnArgs{"ValidateVnfdInstanceBody", "[]uint8", []interface{}{nil}},
		},
		{
			desc:           "Get by ID fail",
			nameorid:       "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638461",
			expectedErr:    vsvcerr.VnfdsvcError{mockErr, http.StatusNotFound},
			daOnReturnArgs: &onCallReturnArgs{"FindVnfdByID", "string", []interface{}{nil, mockErr}},
		},
		{
			desc:            "Get by Name pass",
			nameorid:        "vnfd1",
			expectedErr:     vsvcerr.VnfdsvcError{nil, http.StatusOK},
			daOnReturnArgs:  &onCallReturnArgs{"FindVnfdByName", "string", []interface{}{&vnfd, nil}},
			valOnReturnArgs: &onCallReturnArgs{"ValidateVnfdInstanceBody", "[]uint8", []interface{}{nil}},
		},
		{
			desc:           "Get by Name fail",
			nameorid:       "vnfd1",
			expectedErr:    vsvcerr.VnfdsvcError{mockErr, http.StatusNotFound},
			daOnReturnArgs: &onCallReturnArgs{"FindVnfdByName", "string", []interface{}{nil, mockErr}},
		},
		{
			desc:            "Vnfd instance body validation failure",
			nameorid:        "vnfd12",
			expectedErr:     vsvcerr.VnfdsvcError{mockErr, http.StatusInternalServerError},
			daOnReturnArgs:  &onCallReturnArgs{"FindVnfdByName", "string", []interface{}{&vnfd, nil}},
			valOnReturnArgs: &onCallReturnArgs{"ValidateVnfdInstanceBody", "[]uint8", []interface{}{mockErr}},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			if tc.daOnReturnArgs != nil {
				call := daMock.On(tc.daOnReturnArgs.onCallMethodName, mock.AnythingOfType(tc.daOnReturnArgs.onCallMethodArgType))
				for _, item := range tc.daOnReturnArgs.retArgList {
					call.ReturnArguments = append(call.ReturnArguments, item)
				}
				call.Once()
			}
			if tc.valOnReturnArgs != nil {
				call := valMock.On(tc.valOnReturnArgs.onCallMethodName, mock.AnythingOfType(tc.valOnReturnArgs.onCallMethodArgType))
				for _, item := range tc.valOnReturnArgs.retArgList {
					call.ReturnArguments = append(call.ReturnArguments, item)
				}
				call.Once()
			}
			_, svcerr := svcInstance.GetVnfd(tc.nameorid)
			assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
		})
	}
}

func TestVnfdService_GetVnfds(t *testing.T) {
	// setup test
	mockErr := errors.New("mock error")
	var vnfd models.Vnfd
	err := json.Unmarshal(validInstanceBody, &vnfd)
	if err != nil {
		t.Fatal("Failed to unmarshal")
	}

	var vnfds []models.Vnfd
	vnfds = append(vnfds, vnfd)
	// svcInstance is defined in vnfd_service.go
	svcInstance = &VnfdService{&mockdal, *testcfg}

	type onCallMultipleReturnArgs struct {
		onCallMethodName    string
		onCallMethodArgType []string
		retArgList          []interface{}
	}

	tests := []struct {
		desc            string
		startParam      string
		limitParam      int
		sortParam       string
		expectedErr     vsvcerr.VnfdsvcError
		daOnReturnArgs  *onCallMultipleReturnArgs
		valOnReturnArgs []onCallMultipleReturnArgs
	}{
		{
			desc:           "success->set limit higher than allowed",
			startParam:     "",
			limitParam:     25,
			sortParam:      "",
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusOK},
			daOnReturnArgs: &onCallMultipleReturnArgs{"GetVnfds", []string{"string", "int", "string"}, []interface{}{vnfds, 1, nil}},
			valOnReturnArgs: []onCallMultipleReturnArgs{
				{"ValidateVnfdInstanceBody", []string{"[]uint8"}, []interface{}{nil}},
				{"ValidatePaginatedVnfdsInstancesBody", []string{"[]uint8"}, []interface{}{nil}},
			},
		},
		{
			desc:           "success->set limit within boundary",
			startParam:     "",
			limitParam:     3,
			sortParam:      "created_at",
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusOK},
			daOnReturnArgs: &onCallMultipleReturnArgs{"GetVnfds", []string{"string", "int", "string"}, []interface{}{vnfds, 1, nil}},
			valOnReturnArgs: []onCallMultipleReturnArgs{
				{"ValidateVnfdInstanceBody", []string{"[]uint8"}, []interface{}{nil}},
				{"ValidatePaginatedVnfdsInstancesBody", []string{"[]uint8"}, []interface{}{nil}},
			},
		},
		{
			desc:           "fail->retrieving vnfds from DB",
			startParam:     "",
			limitParam:     3,
			sortParam:      "",
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusInternalServerError},
			daOnReturnArgs: &onCallMultipleReturnArgs{"GetVnfds", []string{"string", "int", "string"}, []interface{}{nil, 0, mockErr}},
		},
		{
			desc:           "fail->individual instance body fails to adhere to schema",
			startParam:     "",
			limitParam:     3,
			sortParam:      "",
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusInternalServerError},
			daOnReturnArgs: &onCallMultipleReturnArgs{"GetVnfds", []string{"string", "int", "string"}, []interface{}{vnfds, 1, nil}},
			valOnReturnArgs: []onCallMultipleReturnArgs{
				{"ValidateVnfdInstanceBody", []string{"[]uint8"}, []interface{}{mockErr}},
			},
		},
		{
			desc:           "vnfds fails to adhere to paginated schema",
			startParam:     "",
			limitParam:     3,
			sortParam:      "",
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusInternalServerError},
			daOnReturnArgs: &onCallMultipleReturnArgs{"GetVnfds", []string{"string", "int", "string"}, []interface{}{vnfds, 1, nil}},
			valOnReturnArgs: []onCallMultipleReturnArgs{
				{"ValidateVnfdInstanceBody", []string{"[]uint8"}, []interface{}{nil}},
				{"ValidatePaginatedVnfdsInstancesBody", []string{"[]uint8"}, []interface{}{mockErr}},
			},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			if tc.daOnReturnArgs != nil {
				daMock := svcInstance.dal.(*dalmocks.VnfdRepository)

				call := daMock.On(tc.daOnReturnArgs.onCallMethodName)
				for _, item := range tc.daOnReturnArgs.onCallMethodArgType {
					call.Arguments = append(call.Arguments, mock.AnythingOfType(item))
				}
				for _, item := range tc.daOnReturnArgs.retArgList {
					call.ReturnArguments = append(call.ReturnArguments, item)
				}
				call.Once()
			}
			lenValArgs := len(tc.valOnReturnArgs)
			if lenValArgs != 0 {
				// vnfdVerifier is defined in vnfd_service.go
				vnfdVerifier = &umocks.VnfdValidator{}
				valMock := vnfdVerifier.(*umocks.VnfdValidator)

				for _, item := range tc.valOnReturnArgs {

					call := valMock.On(item.onCallMethodName)
					for _, args := range item.onCallMethodArgType {
						call.Arguments = append(call.Arguments, mock.AnythingOfType(args))
					}
					for _, elem := range item.retArgList {
						call.ReturnArguments = append(call.ReturnArguments, elem)
					}
				}
			}
			_, svcerr := svcInstance.GetVnfds(tc.startParam, tc.limitParam, tc.sortParam)
			assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
		})
	}
}

func TestVnfdService_GetInputParamsSchemaForVnfd(t *testing.T) {
	// svcInstance is defined in vnfd_service.go
	svcInstance = &VnfdService{&mockdal, *testcfg}

	tests := []struct {
		desc        string
		inpVnfdJson []byte
		expectedErr vsvcerr.VnfdsvcError
	}{
		{
			desc:        "success: retrieve input params schema",
			inpVnfdJson: validInstanceBody,
			expectedErr: vsvcerr.VnfdsvcError{nil, http.StatusOK},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			_, svcerr := svcInstance.GetInputParamsSchemaForVnfd(tc.inpVnfdJson)
			fmt.Println(svcerr.HttpCode, tc.expectedErr.HttpCode)
			assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
		})
	}
}
