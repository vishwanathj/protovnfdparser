package mongo

import (
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"gopkg.in/mgo.v2"
)

// VnfdKey the key to be used for the VNFD document collection
const VnfdKey = "name"

// VnfdID the key that holds the Vnfd ID
const VnfdID = "_id"

type vnfdMgoModel struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Created_At string `json:"created_at"`
	Vdus       []*struct {
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

func vnfdMgoModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{VnfdKey},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func (v *vnfdMgoModel) toModelVnfd() *models.Vnfd {
	return &models.Vnfd{
		ID:            v.ID,
		Name:          v.Name,
		Status:        v.Status,
		CreatedAt:     v.Created_At,
		Vdus:          v.Vdus,
		Virtual_Links: v.Virtual_Links,
	}
}

func toVnfdMgoModel(v *models.Vnfd) *vnfdMgoModel {
	return &vnfdMgoModel{
		ID:            v.ID,
		Name:          v.Name,
		Status:        v.Status,
		Created_At:    v.CreatedAt,
		Vdus:          v.Vdus,
		Virtual_Links: v.Virtual_Links,
	}
}
