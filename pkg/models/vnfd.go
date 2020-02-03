package models

import (
	"time"

	"github.com/vishwanathj/protovnfdparser/pkg/constants"
	"github.com/vishwanathj/protovnfdparser/pkg/errors"

	"github.com/satori/go.uuid"
)

// Vnfd Struct for holding Vnfd data
// Refer https://www.sohamkamani.com/blog/golang/2018-07-19-golang-omitempty/
// for the reason Vdus, Constraints, HA, ScaleInOut, VLs are declared as pointers
type Vnfd struct {
	//ID        string `json:"id,omitempty" bson:"_id"`
	// remove the above comment if there are failures in creating and retrieving VNFDs
	// below is a refactor attempting to get rid of `"bson:" _id"' that is mongo specific
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Status    string `json:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Vdus      []*struct {
		Constraints *struct {
			Dedicated interface{} `json:"dedicated,omitempty"`
			Vim_ID    string      `json:"vim_id,omitempty"`
		} `json:"constraints,omitempty"`
		Disk_Size         interface{} `json:"disk_size"`
		High_Availability *string     `json:"high_availability,omitempty"`
		Image             string      `json:"image"`
		Memory            interface{} `json:"memory"`
		Name              string      `json:"name"`
		Scale_In_Out      *struct {
			Default interface{} `json:"default"`
			Maximum interface{} `json:"maximum"`
			Minimum interface{} `json:"minimum"`
		} `json:"scale_in_out,omitempty"`
		Vcpus interface{} `json:"vcpus"`
		Vnfcs []struct {
			Connection_Points []struct {
				IP_Address           string   `json:"ip_address"`
				Name                 string   `json:"name"`
				VirtualLinkReference []string `json:"virtualLinkReference"`
			} `json:"connection_points"`
			Name string `json:"name"`
		} `json:"vnfcs"`
	} `json:"vdus"`
	Virtual_Links []*struct {
		Name          string `json:"name"`
		Is_management bool   `json:"is_management,omitempty"`
	} `json:"virtual_links"`
}

// SetCreationTimeAttributes sets the ID, creation timestamp and initial status
func (v *Vnfd) SetCreationTimeAttributes() {
	id := uuid.NewV4()
	v.CreatedAt = time.Now().Format(time.RFC3339)
	v.ID = constants.VnfdIDPrefix + id.String()
	v.Status = constants.Available
}

// VnfdService defines methods that implement the VnfdService interface
type VnfdService interface {
	CreateVnfd(v *Vnfd) errors.VnfdsvcError
	GetByVnfdname(vnfdname string) (*Vnfd, errors.VnfdsvcError)
	GetByVnfdID(vnfdID string) (*Vnfd, errors.VnfdsvcError)
	GetVnfds(start string, limit int, sort string) (PaginatedVnfds, errors.VnfdsvcError)
	GetInputParamsSchemaForVnfd(vnfdjson []byte) ([]byte, errors.VnfdsvcError)
	GetHealth() string
	//GetReadiness() (string)
}

/*
// Refer https://www.sohamkamani.com/blog/golang/2018-07-19-golang-omitempty/
// for the reason Vdus, Constraints, HA, ScaleInOut, VLs are declared as pointers
type Vnfd struct {
	ID 			string 			`json:"id,omitempty"`
	Name		string			`json:"name"`
	Status		string			`json:"status"`
	CreatedAt	string			`json:"created_at"`
	Vdus []*struct {
		Constraints *struct {
			Dedicated  string `json:"dedicated,omitempty"`
			Vim_ID     string `json:"vim_id,omitempty"`
		} `json:"constraints,omitempty"`
		DiskSize         	string `json:"disk_size"`
		High_Availability 	*string `json:"high_availability,omitempty"`
		Image            	string `json:"image"`
		Memory           	string `json:"memory"`
		Name             	string `json:"name"`
		Scale_In_Out       *struct {
			Default string `json:"default"`
			Maximum string `json:"maximum"`
			Minimum string `json:"minimum"`
		} `json:"scale_in_out,omitempty"`
		Vcpus string `json:"vcpus"`
		Vnfcs []struct {
			Connection_Points []struct {
				IP_Address           string   `json:"ip_address"`
				Name                 string   `json:"name"`
				VirtualLinkReference []string `json:"virtualLinkReference"`
			} `json:"connection_points"`
			Name string `json:"name"`
		} `json:"vnfcs"`
	} `json:"vdus"`
	VirtualLink []*struct {
		Name 			string 		`json:"name"`
		Is_management 	bool		`json:"is_management,omitempty"`
	} `json:"virtual_links"`
}
*/

/*
type Constraint struct {
	Dedicated  string `json:"dedicated,omitempty"`
	Vim_ID     string `json:"vim_id,omitempty"`
}

type ScaleInOut struct {
	Default string `json:"default"`
	Maximum string `json:"maximum"`
	Minimum string `json:"minimum"`
}

type ConnectionPoint struct {
	IP_Address           string   `json:"ip_address"`
	Name                 string   `json:"name"`
	VirtualLinkReference []string `json:"virtualLinkReference"`
}

type Vnfc struct {
	Connection_Points []ConnectionPoint `json:"connection_points"`
	Name string `json:"name"`
}

type Vdu struct {
	Constraints 		*Constraint `json:"constraints,omitempty"`
	DiskSize         	string `json:"disk_size"`
	High_Availability 	*string `json:"high_availability,omitempty"`
	Image            	string `json:"image"`
	Memory           	string `json:"memory"`
	Name             	string `json:"name"`
	Scale_In_Out       *ScaleInOut `json:"scale_in_out,omitempty"`
	Vcpus 				string `json:"vcpus"`
	Vnfcs 				[]Vnfc `json:"vnfcs"`
}

type VirtualLink struct {
	Name 			string 		`json:"name"`
	Is_management 	bool		`json:"is_management,omitempty"`
}

// Refer https://www.sohamkamani.com/blog/golang/2018-07-19-golang-omitempty/
// for the reason Vdus, Constraints, HA, ScaleInOut, VLs are declared as pointers
type Vnfd struct {
	ID 			string 			`json:"id,omitempty"`
	Name		string			`json:"name"`
	Status		string			`json:"status"`
	CreatedAt	string			`json:"created_at"`
	Vdus 		[]*Vdu 			`json:"vdus"`
	VirtualLink []*VirtualLink `json:"virtual_links"`
}*/
