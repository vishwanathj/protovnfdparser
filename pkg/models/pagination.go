package models

import (
	"bytes"
	"strconv"
)

/*func baseURI() string {
	var ip string
	var port int
	if flag.Lookup("ip") != nil {
		ip = flag.Lookup("ip").Value.(flag.Getter).Get().(string)
	} else {
		ip = "localhost"
	}
	if flag.Lookup("port") != nil {
		port = flag.Lookup("port").Value.(flag.Getter).Get().(int)
	} else {
		port = 8080
	}
	url := fmt.Sprintf("http://%s:%d", ip, port)
	return url
}*/
// BaseURI  The base URI that's used by REST assured when making requests if
// a non-fully qualified URI is used in the request
const BaseURI = "http://localhost:8080"

// DefaultLimit is the count of results that would be returned when no value
// for Limit attribute is specified
const DefaultLimit = 5

// MaxLimit is the maximum value that is supported for the the Limit attribute
const MaxLimit = 10

// MinLimit is the minimum value that is supported for the the Limit attribute
const MinLimit = 1

// Link captures the URL info
type Link struct {
	Href string `json:"href" binding:"required,url"`
}

// PaginatedCollection captures the attributes required for pagination results
type PaginatedCollection struct {
	Limit      int   `json:"limit" binding:"required,min=1,max=10"`
	TotalCount int   `json:"total_count" binding:"required,gte=0"`
	First      *Link `json:"first" binding:"required"`
	Next       *Link `json:"next,omitempty"`
}

// DefaultOrderBy will store the property definitions should be ordered by,
// by default
const DefaultOrderBy = "name"

// PaginationQueryParameters contains the query parameters a user can specify in a GET
// list request
type PaginationQueryParameters struct {
	OrderBy string
}

// MakeFirstHref -- make href for First link in collection
/*
At the moment only one PaginationQueryParameters object is expected. This
object will store additional query parameters (such as sort) which should be
appended to the Next Href.
*/
func MakeFirstHref(limit int, apipath string, vars ...PaginationQueryParameters) string {
	var buffer bytes.Buffer
	//buffer.WriteString(ctx.Value(riaascontext.RiaasBaseURI).(string))
	buffer.WriteString(BaseURI)
	//buffer.WriteString(baseURI())
	buffer.WriteString(apipath)
	buffer.WriteString("?limit=")
	buffer.WriteString(strconv.Itoa(limit))
	if len(vars) > 0 {
		orderBy := vars[0].OrderBy
		if orderBy != DefaultOrderBy && orderBy != "" {
			buffer.WriteString("&sort=")
			buffer.WriteString(orderBy)
		}
	}
	return buffer.String()
}

// MakeNextHref -- make href for next link in collection
/*
At the moment only one PaginationQueryParameters object is expected. This
object will store additional query parameters (such as sort) which should be
appended to the Next Href.
*/
func MakeNextHref(limit int, start string, apipath string, vars ...PaginationQueryParameters) string {
	var buffer bytes.Buffer
	//buffer.WriteString(ctx.Value(riaascontext.RiaasBaseURI).(string))
	buffer.WriteString(BaseURI)
	//buffer.WriteString(baseURI())
	buffer.WriteString(apipath)
	buffer.WriteString("?start=")
	buffer.WriteString(start)
	buffer.WriteString("&limit=")
	buffer.WriteString(strconv.Itoa(limit))
	if len(vars) > 0 {
		orderBy := vars[0].OrderBy
		if orderBy != DefaultOrderBy && orderBy != "" {
			buffer.WriteString("&sort=")
			buffer.WriteString(orderBy)
		}
	}
	return buffer.String()
}
