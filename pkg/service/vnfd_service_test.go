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
	dalmocks "github.com/vishwanathj/protovnfdparser/pkg/dataaccess/mocks"
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

/*func TestVnfdService_CreateVnfd(t *testing.T) {
	mockErr := errors.New("mock error")

	// svcInstance is defined in vnfd_service.go
	svcInstance = &VnfdService{&mockdal, *testcfg}

	tests := []struct {
		desc                string
		postBody            []byte
		postBodyErr         bool
		postBodySuccess     bool
		instanceBodyErr     bool
		instanceBodySuccess bool
		insertDBErr         bool
		insertDBSuccess     bool
		expectedErr         vsvcerr.VnfdsvcError
	}{
		{
			desc:        "vcpu missing, post body non-compliant with schema",
			postBody:    createVnfdObjVcpusMissing,
			expectedErr: vsvcerr.VnfdsvcError{mockErr, http.StatusBadRequest},
		},
		{
			desc:            "vnfd instance body non-compliant with schema",
			postBody:        createVnfdObj,
			postBodySuccess: true,
			instanceBodyErr: true,
			expectedErr:     vsvcerr.VnfdsvcError{mockErr, http.StatusBadRequest},
		},
		{
			desc:                "insert into DB passes",
			postBody:            createVnfdObj,
			postBodySuccess:     true,
			instanceBodySuccess: true,
			insertDBSuccess:     true,
			expectedErr:         vsvcerr.VnfdsvcError{nil, http.StatusOK},
		},
		{
			desc:        "insert into DB fails",
			postBody:    createVnfdObj,
			insertDBErr: true,
			expectedErr: vsvcerr.VnfdsvcError{mockErr, http.StatusConflict},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			var vnfd models.Vnfd
			err := json.Unmarshal(tc.postBody, &vnfd)
			if err != nil {
				t.Fatal("Failed to unmarshal")
			}

			// svcInstance is defined in vnfd_service.go
			//svcInstance = &VnfdService{&mockdal, *testcfg}

			// pre setup, prior to invoking CreateVnfd
			if tc.postBodyErr || tc.instanceBodyErr || tc.postBodySuccess || tc.instanceBodySuccess {
				// vnfdVerifier is defined in vnfd_service.go
				vnfdVerifier = &umocks.VnfdValidator{}
				if tc.postBodyErr {
					// []uint8 used instead of []byte due to https://github.com/stretchr/testify/issues/387
					vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdPostBody", mock.AnythingOfType("[]uint8")).Return(mockErr)
				}
				if tc.postBodySuccess {
					vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdPostBody", mock.AnythingOfType("[]uint8")).Return(nil)
				}
				if tc.instanceBodyErr {
					vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdInstanceBody", mock.AnythingOfType("[]uint8")).Return(mockErr)
				}
				if tc.instanceBodySuccess {
					vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdInstanceBody", mock.AnythingOfType("[]uint8")).Return(nil)
				}
			}

			if tc.insertDBErr {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("InsertVnfd", mock.AnythingOfType("*models.Vnfd")).Return(mockErr).Once()
			}
			if tc.insertDBSuccess {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("InsertVnfd", mock.AnythingOfType("*models.Vnfd")).Return(nil).Once()
			}

			svcerr := svcInstance.CreateVnfd(&vnfd)
			fmt.Println(svcerr.HttpCode, tc.expectedErr.HttpCode)
			assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
		})
	}
}

/*
func TestVnfdService_GetVnfd2(t *testing.T) {
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

	tests := []struct {
		desc             string
		nameorid         string
		findByIDErr      bool
		findByIDPass     bool
		findByNameErr    bool
		findByNamePass   bool
		instanceBodyFail bool
		expectedErr      vsvcerr.VnfdsvcError
	}{
		{
			desc:         "Get by ID pass",
			nameorid:     "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638451",
			findByIDPass: true,
			expectedErr:  vsvcerr.VnfdsvcError{nil, http.StatusOK},
		},
		{
			desc:        "Get by ID fail",
			nameorid:    "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638461",
			findByIDErr: true,
			expectedErr: vsvcerr.VnfdsvcError{mockErr, http.StatusNotFound},
		},
		{
			desc:           "Get by Name pass",
			nameorid:       "vnfd1",
			findByNamePass: true,
			expectedErr:    vsvcerr.VnfdsvcError{nil, http.StatusOK},
		},
		{
			desc:          "Get by Name fail",
			nameorid:      "vnfd1",
			findByNameErr: true,
			expectedErr:   vsvcerr.VnfdsvcError{mockErr, http.StatusNotFound},
		},
		{
			desc:             "Vnfd instance body validation failure",
			nameorid:         "vnfd12",
			instanceBodyFail: true,
			expectedErr:      vsvcerr.VnfdsvcError{mockErr, http.StatusInternalServerError},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			if tc.findByIDPass {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("FindVnfdByID", mock.AnythingOfType("string")).Return(&vnfd, nil).Once()
				vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdInstanceBody", mock.AnythingOfType("[]uint8")).Return(nil).Once()
			}
			if tc.findByNamePass {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("FindVnfdByName", mock.AnythingOfType("string")).Return(&vnfd, nil).Once()
				vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdInstanceBody", mock.AnythingOfType("[]uint8")).Return(nil).Once()
			}
			if tc.findByIDErr {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("FindVnfdByID", mock.AnythingOfType("string")).Return(nil, mockErr).Once()
			}
			if tc.findByNameErr {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("FindVnfdByName", mock.AnythingOfType("string")).Return(nil, mockErr).Once()
			}

			if tc.instanceBodyFail {
				svcInstance.dal.(*dalmocks.VnfdRepository).On("FindVnfdByName", mock.AnythingOfType("string")).Return(&vnfd, nil).Once()
				vnfdVerifier.(*umocks.VnfdValidator).On("ValidateVnfdInstanceBody", mock.AnythingOfType("[]uint8")).Return(mockErr).Once()
			}

			_, svcerr := svcInstance.GetVnfd(tc.nameorid)
			assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
		})
	}
}*/

func TestVnfdService_CreateVnfd(t *testing.T) {
	mockErr := errors.New("mock error")

	// svcInstance is defined in vnfd_service.go
	svcInstance = &VnfdService{&mockdal, *testcfg}

	type onCallReturnArgs struct {
		onCallMethodName    string
		onCallMethodArgType string
		retArgList          []interface{}
	}

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

	type onCallReturnArgs struct {
		onCallMethodName    string
		onCallMethodArgType string
		retArgList          []interface{}
	}

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
