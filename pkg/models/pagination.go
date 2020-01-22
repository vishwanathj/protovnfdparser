package models

import (
	"bytes"
	"strconv"

	"github.com/vishwanathj/protovnfdparser/pkg/config"
	"github.com/vishwanathj/protovnfdparser/pkg/constants"
)

// Link captures the URL info
type Link struct {
	Href string `json:"href" binding:"required,url"`
}

// PaginatedVnfds Struct for holding paginated Vnfd responses
type PaginatedVnfds struct {
	Limit      int    `json:"limit" binding:"required,min=1,max=10"`
	TotalCount int    `json:"total_count" binding:"required,gte=0"`
	First      *Link  `json:"first" binding:"required"`
	Next       *Link  `json:"next,omitempty"`
	Vnfds      []Vnfd `json:"vnfds,omitempty"`
}

// PgnQueryParams contains the query parameters a user can specify in a GET
// list request
type PgnQueryParams struct {
	OrderBy string
}

// MakeFirstHref -- make href for First link in collection
func MakeFirstHref(cfg config.PaginationConfig, limit int, apipath string, vars ...PgnQueryParams) string {
	var buffer bytes.Buffer
	buffer.WriteString(cfg.BaseURI + apipath + "?" + constants.PaginationURLLimit + "=" + strconv.Itoa(limit))
	//currently one only one PgnQueryParams object is expected.
	if len(vars) > 0 {
		orderBy := vars[0].OrderBy
		if orderBy != cfg.DefaultOrderBy && orderBy != "" {
			buffer.WriteString("&" + constants.PaginationURLSort + "=")
			buffer.WriteString(orderBy)
		}
	}
	return buffer.String()
}

// MakeNextHref -- make href for next link in collection
func MakeNextHref(cfg config.PaginationConfig, limit int, start string, apipath string, vars ...PgnQueryParams) string {
	var buffer bytes.Buffer
	buffer.WriteString(cfg.BaseURI + apipath + "?" + constants.PaginationURLStart + "=" + start)
	buffer.WriteString("&" + constants.PaginationURLLimit + "=" + strconv.Itoa(limit))
	//currently one only one PgnQueryParams object is expected.
	if len(vars) > 0 {
		orderBy := vars[0].OrderBy
		if orderBy != cfg.DefaultOrderBy && orderBy != "" {
			buffer.WriteString("&" + constants.PaginationURLSort + "=")
			buffer.WriteString(orderBy)
		}
	}
	return buffer.String()
}
