// +build integration

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package mongo_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/vishwanathj/protovnfdparser/pkg/constants"

	"github.com/stretchr/testify/assert"
	vsvcerror "github.com/vishwanathj/protovnfdparser/pkg/errors"

	"github.com/vishwanathj/protovnfdparser/pkg/config"

	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"github.com/vishwanathj/protovnfdparser/pkg/service"

	"github.com/ghodss/yaml"
	"github.com/vishwanathj/protovnfdparser/pkg/mongo"
)

var mongoUrl string
var dbName string
var vnfdCollectionName string
var testcfg = config.GetConfigInstance()

func init() {
	dbName = testcfg.MongoDBConfig.MongoDBName
	vnfdCollectionName = testcfg.MongoDBConfig.MongoColName
	mongoUrl = fmt.Sprintf("%s:%d", testcfg.MongoDBConfig.MongoIP, testcfg.MongoDBConfig.MongoPort)

}

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
var insertVnfdOne = []byte(`{
    "_id" : "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638451",
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

// Make sure to copy the entire JSON from an already existing VNFD in mongodb for the below variable.
// Missing some fields would cause the test to file.
// In this case the status is set to "" which is invalid
var insertVnfdOneInvalid = []byte(`{
    "_id" : "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638461",
    "name" : "vnfdOptProps",
    "status" : "",
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

func TestGetVnfds(t *testing.T) {
	tests := []struct {
		desc        string
		startParam  string
		limitParam  int
		sortParam   string
		jsonInput   []byte
		expectedErr *vsvcerror.VnfdsvcError
		expectedOut *models.PaginatedVnfds
	}{
		{
			desc:        "zero vnfd retrieval success",
			startParam:  "",
			limitParam:  1,
			sortParam:   "created_at",
			jsonInput:   nil,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusOK},
			expectedOut: &models.PaginatedVnfds{
				Limit:      1,
				TotalCount: 0,
				First:      &models.Link{Href: models.MakeFirstHref(*testcfg.PgntConfig, 1, constants.ApiPathVnfds, models.PgnQueryParams{OrderBy: "created_at"})},
				Next:       nil,
				Vnfds:      nil,
			},
		},
		{
			desc:        "single vnfd retrieval success",
			startParam:  "",
			limitParam:  1,
			sortParam:   "created_at",
			jsonInput:   insertVnfdOne,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusOK},
			expectedOut: &models.PaginatedVnfds{
				Limit:      1,
				TotalCount: 1,
				First:      &models.Link{Href: models.MakeFirstHref(*testcfg.PgntConfig, 1, constants.ApiPathVnfds, models.PgnQueryParams{OrderBy: "created_at"})},
				Next:       &models.Link{Href: models.MakeNextHref(*testcfg.PgntConfig, 1, "vnfdOptProps", constants.ApiPathVnfds, models.PgnQueryParams{OrderBy: "created_at"})},
			},
		},
		{
			desc:        "fail to retrieve vnfd that does not adhere to schema",
			startParam:  "",
			limitParam:  1,
			sortParam:   "created_at",
			jsonInput:   insertVnfdOneInvalid,
			expectedErr: &vsvcerror.VnfdsvcError{errors.New("I[#] S[#/vnfdInstance/required] missing properties: \"status\""), http.StatusInternalServerError},
			expectedOut: nil,
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			// check if mongodb needs to be pre-populated for the test
			if tc.jsonInput != nil {
				session, err := mongo.NewSession(mongoUrl)
				if err != nil {
					log.Fatalf("Unable to connect to mongo: %s", err)
				}
				defer func() {
					session.DropDatabase(dbName)
					session.Close()
				}()
				var m interface{}
				err = yaml.Unmarshal(tc.jsonInput, &m)
				if err != nil {
					panic(err)
					t.Fail()
				}
				e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
				if e != nil {
					t.Fatal()
				}
			}

			vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
			if err != nil {
				panic(err)
				t.Fatal()
			}

			pvnfds, svcerr := vnfdService.GetVnfds(tc.startParam, tc.limitParam, tc.sortParam)

			if tc.expectedErr != nil {
				//assert.Equal(t, ssvcerr.OrigError, tc.expectedErr.OrigError, "error string comparison")
				assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
				//assert.EqualValues(t, svcerr, *tc.expectedErr, "error comparison")
			}
			if tc.expectedOut != nil {
				assert.Equal(t, tc.expectedOut.Limit, pvnfds.Limit, "Limit comparison")
				assert.Equal(t, tc.expectedOut.TotalCount, pvnfds.TotalCount, "TotalCount comparison")
				assert.Equal(t, tc.expectedOut.First, pvnfds.First, "Link First comparison")
				assert.Equal(t, tc.expectedOut.Next, pvnfds.Next, "Link Next comparison")
				assert.Equal(t, tc.expectedOut.TotalCount, len(pvnfds.Vnfds), "Vnfds count comparison")
				//assert.EqualValues(t, pvnfds, *tc.expectedOut, "output comparison")
			}
		})
	}
}

func TestGetVnfd(t *testing.T) {
	tests := []struct {
		desc        string
		nameorid    string
		jsonInput   []byte
		expectedErr *vsvcerror.VnfdsvcError
	}{
		{
			desc:        "get vnfd by id successfully",
			nameorid:    "vnfdOptProps",
			jsonInput:   insertVnfdOne,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusOK},
		},
		{
			desc:        "get vnfd by name successfully",
			nameorid:    "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638451",
			jsonInput:   insertVnfdOne,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusOK},
		},
		{
			desc:        "get vnfd by id fails for non existing name",
			nameorid:    "vnfdProps",
			jsonInput:   insertVnfdOne,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusNotFound},
		},
		{
			desc:        "get vnfd by name fails for non existing id",
			nameorid:    "VNFD-30c270ff-47c4-4d66-8a6f-f24de7638451",
			jsonInput:   insertVnfdOne,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusNotFound},
		},
		{
			desc:        "get vnfd by name fails due to non-compliant with schema",
			nameorid:    "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638461",
			jsonInput:   insertVnfdOneInvalid,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusInternalServerError},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			// check if mongodb needs to be pre-populated for the test
			if tc.jsonInput != nil {
				session, err := mongo.NewSession(mongoUrl)
				if err != nil {
					log.Fatalf("Unable to connect to mongo: %s", err)
				}
				defer func() {
					session.DropDatabase(dbName)
					session.Close()
				}()
				var m interface{}
				err = yaml.Unmarshal(tc.jsonInput, &m)
				if err != nil {
					panic(err)
					t.Fail()
				}
				e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
				if e != nil {
					t.Fatal()
				}
			}

			vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
			if err != nil {
				panic(err)
				t.Fatal()
			}

			vnfd, svcerr := vnfdService.GetVnfd(tc.nameorid)
			if tc.expectedErr != nil {
				assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
			}
			if vnfd != nil && tc.nameorid != vnfd.Name && tc.nameorid != vnfd.ID {
				t.Fail()
			}
		})
	}
}

func TestCreateVnfd(t *testing.T) {
	tests := []struct {
		desc        string
		jsonInput   []byte
		expectedErr *vsvcerror.VnfdsvcError
	}{
		{
			desc:        "create vnfd successfully",
			jsonInput:   createVnfdObj,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusOK},
		},
		{
			desc:        "create vnfd fails as json does not adhere to schema",
			jsonInput:   insertVnfdOneInvalid,
			expectedErr: &vsvcerror.VnfdsvcError{nil, http.StatusBadRequest},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
			if err != nil {
				panic(err)
				t.Fatal()
			}

			var vnfd models.Vnfd

			err = json.Unmarshal(tc.jsonInput, &vnfd)
			if err != nil {
				t.Fatal("Failed to unmarshal")
			}

			svcerr := vnfdService.CreateVnfd(&vnfd)

			if tc.expectedErr != nil {
				assert.Equal(t, svcerr.HttpCode, tc.expectedErr.HttpCode, "error code comparison")
			} else {
				session, err := mongo.NewSession(mongoUrl)
				if err != nil {
					log.Fatalf("Unable to connect to mongo: %s", err)
				}
				defer func() {
					session.DropDatabase(dbName)
					session.Close()
				}()

				var results []models.Vnfd
				session.GetCollection(dbName, vnfdCollectionName).Find(nil).All(&results)
				assert.Equal(t, len(results), 1, "vnfd count comparison")
				assert.Equal(t, results[0].Name, vnfd.Name, "vnfd names comparison")
				assert.Equal(t, len(vnfd.ID), 0, "verify ID not present prior to creation")
				assert.NotEqual(t, vnfd.ID, results[0].ID, "comparing IDs to be not same")
			}
		})
	}
}

func TestVnfdService(t *testing.T) {
	t.Run("incorrectURL", incorrectURL_Negative)
}

func incorrectURL_Negative(t *testing.T) {
	session, err := mongo.NewSession("")
	if session == nil && err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}

/*
Below can be an example for testing out pagination; Note the use of the Replace function
var createVnfdObj = []byte(`{
  "virtual_links": [
    {
      "name": "mgmt_net",
      "is_management": true
    }
  ],
  "name": "$VNFDNAME"
}`)

func main() {
	content := []byte(strings.Replace(string(createVnfdObj), "$VNFDNAME", "VISH", 1))
	fmt.Println(string(content))
}
*/

/*
Below is an example on how to initialize a nested [] struct

vnfd := &root.Vnfd{

	Name: "vnfname",
	Vdus: []struct {
		Constraints struct {
			Dedicated string `json:"dedicated",omitempty`
			Vim_ID     string `json:"vim_id",omitempty`
		} `json:"constraints",omitempty`
		DiskSize         	string `json:"disk_size"`
		High_Availability 	string `json:"high_availability",omitempty`
		Image            	string `json:"image"`
		Memory           	string `json:"memory"`
		Name             	string `json:"name"`
		Scale_In_Out       struct {
			Default string `json:"default"`
			Maximum string `json:"maximum"`
			Minimum string `json:"minimum"`
		} `json:"scale_in_out",omitempty`
		Vcpus string `json:"vcpus"`
		Vnfcs []struct {
			Connection_Points []struct {
				IP_Address           string   `json:"ip_address"`
				Name                 string   `json:"name"`
				VirtualLinkReference []string `json:"virtualLinkReference"`
			} `json:"connection_points"`
			Name string `json:"name"`
		} `json:"vnfcs"`
	}{
		{DiskSize: "$disksize", Image: "$image", Vcpus: "$vcpus", Name: "vdu1", Memory: "$mem",
		Vnfcs: []struct {
			Connection_Points []struct {
				IP_Address           string   `json:"ip_address"`
				Name                 string   `json:"name"`
				VirtualLinkReference []string `json:"virtualLinkReference"`
			} `json:"connection_points"`
			Name string `json:"name"`
		}{ {
			Name: "mgmtCP",
			Connection_Points: []struct {
				IP_Address           string   `json:"ip_address"`
				Name                 string   `json:"name"`
				VirtualLinkReference []string `json:"virtualLinkReference"`}{
					{IP_Address: "$vdu1_vnfc1_mgmt_net_subnet", Name: "mgmtCP1", VirtualLinkReference: []string{"mgmt_net"}},
			}} }}},
	VirtualLink: []struct {
		Name 			string 		`json:"name"`
		Is_management 	bool		`json:"is_management",omitempty`
	} {
		{Name: "mgmt_net", Is_management: true}}}
*/
