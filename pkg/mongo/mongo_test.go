// +build integration

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package mongo_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

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

func TestVnfdService(t *testing.T) {
	t.Run("CreateVnfd", createVnfd_single_vnfd_insert_for_get_test)
	t.Run("getVnfdByName", getVnfdByName)
	t.Run("getByVnfdID", getByVnfdID)
	t.Run("getByVnfdIDFail", getByVnfdIDFail)
	t.Run("getVnfds", getVnfds)
	t.Run("getVnfdsFail", getVnfdsFail)
	t.Run("incorrectURL", incorrectURL_Negative)
}

func getVnfds(t *testing.T) {
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
	if err != nil {
		panic(err)
		t.Fail()
	}
	pvnfds, err := vnfdService.GetVnfds("", 1, "created_at")
	if pvnfds.Vnfds == nil && pvnfds.Next == nil &&
		pvnfds.TotalCount == 0 && pvnfds.Limit == 1 && pvnfds.First != nil {
		t.Log("Success")
	} else {
		t.Fail()
	}

	var m interface{}
	err = yaml.Unmarshal(insertVnfdOne, &m)
	if err != nil {
		panic(err)
		t.Fail()
	}
	e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
	if e != nil {
		t.Fatal()
	}
	pvnfds, err = vnfdService.GetVnfds("", 1, "created_at")
	if pvnfds.Vnfds != nil && pvnfds.Next != nil &&
		pvnfds.TotalCount == 1 && pvnfds.Limit == 1 && pvnfds.First != nil {
		t.Log("Success")
	} else {
		t.Fail()
	}
}

func getVnfdsFail(t *testing.T) {
	//mongoUrl := fmt.Sprintf("%s:%d", *mgoIpPtr, *mgoPortPtr)
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
	if err != nil {
		panic(err)
		t.Fail()
	}
	pvnfds, err := vnfdService.GetVnfds("", 1, "created_at")
	if pvnfds.Vnfds == nil && pvnfds.Next == nil &&
		pvnfds.TotalCount == 0 && pvnfds.Limit == 1 && pvnfds.First != nil {
		t.Log("Success")
	} else {
		t.Fail()
	}

	var m interface{}
	err = yaml.Unmarshal(insertVnfdOneInvalid, &m)
	if err != nil {
		panic(err)
		t.Fail()
	}
	e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
	if e != nil {
		t.Fatal()
	}
	pvnfds, err = vnfdService.GetVnfds("", 1, "created_at")

	if err != nil {
		t.Log("Success")
		t.Log(err)
	} else {
		t.Fail()
	}
}

func getByVnfdID(t *testing.T) {
	//mongoUrl := fmt.Sprintf("%s:%d", *mgoIpPtr, *mgoPortPtr)
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	var m interface{}
	err = yaml.Unmarshal(insertVnfdOne, &m)
	if err != nil {
		panic(err)
		t.Fail()
	}
	e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
	if e != nil {
		t.Fatal()
	}

	vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
	if err != nil {
		panic(err)
		t.Fail()
	}
	testVnfdID := "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638451"

	obj, err := vnfdService.GetByVnfdID(testVnfdID)
	if err != nil {
		t.Error("vnfdService.GetByVnfdname FAILED")
	} else {
		t.Log(obj.ID)
	}
}

func getByVnfdIDFail(t *testing.T) {
	//mongoUrl := fmt.Sprintf("%s:%d", *mgoIpPtr, *mgoPortPtr)
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	var m interface{}
	err = yaml.Unmarshal(insertVnfdOneInvalid, &m)
	if err != nil {
		panic(err)
		t.Fail()
	}
	e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
	if e != nil {
		t.Fatal()
	}

	vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
	if err != nil {
		panic(err)
		t.Fail()
	}
	testVnfdID := "VNFD-50c270ff-47c4-4d66-8a6f-f24de7638461"

	obj, err := vnfdService.GetByVnfdID(testVnfdID)
	if err != nil {
		t.Log("Success")
		t.Log(err)
	} else {
		t.Fail()
		t.Log(obj)
	}
}

func getVnfdByName(t *testing.T) {
	//mongoUrl := fmt.Sprintf("%s:%d", *mgoIpPtr, *mgoPortPtr)
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	var m interface{}
	err = yaml.Unmarshal(insertVnfdOne, &m)
	if err != nil {
		panic(err)
		t.Fail()
	}
	e := session.GetCollection(dbName, vnfdCollectionName).Insert(m)
	if e != nil {
		t.Fatal()
	}

	vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
	if err != nil {
		panic(err)
		t.Fail()
	}
	testVnfdName := "vnfdOptProps"

	obj, err := vnfdService.GetByVnfdname(testVnfdName)
	if err != nil {
		t.Error("vnfdService.GetByVnfdname FAILED")
	} else {
		t.Log(obj.Name)
	}
}

func createVnfd_single_vnfd_insert_for_get_test(t *testing.T) {
	//Arrange
	//mongoUrl := fmt.Sprintf("%s:%d", *mgoIpPtr, *mgoPortPtr)
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	vnfdService, err := service.GetVnfdServiceInstance(*testcfg)
	if err != nil {
		t.Fatal("Failed to create dataacess layer object")
	}

	testVnfdname := "vnfname"

	var vnfd models.Vnfd

	err = json.Unmarshal(createVnfdObj, &vnfd)

	err = vnfdService.CreateVnfd(&vnfd)

	//Assert
	if err != nil {
		t.Errorf("Unable to create vnfd: %s", err)
	}
	var results []models.Vnfd
	session.GetCollection(dbName, vnfdCollectionName).Find(nil).All(&results)

	count := len(results)

	if count != 1 {
		//t.Error("Incorrect number of results. Expected `1`, got: `%i`", count)
		t.Error("Incorrect number of results.  ", count)
		t.Fatal()
	}
	if results[0].Name != vnfd.Name {
		t.Errorf("Incorrect Vnfdname. Expected `%s`, Got: `%s`", testVnfdname, results[0].Name)
		t.Fatal()
	}
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
