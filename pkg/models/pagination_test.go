// +build unit

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeFirstHref(t *testing.T) {
	//c, _ := apitestutils.CreateTestContext("", nil, "")
	r := MakeFirstHref(42, "/some/path/here")
	assert.Equal(t, "http://localhost:8080/some/path/here?limit=42", r)
}

var flagtestsForTestMakeFirstHrefWhenQueryIsPassedIn = []struct {
	limit       int
	queryParams PaginationQueryParameters
	expectedURL string
}{
	{10, PaginationQueryParameters{"name"}, "http://localhost:8080/some/path/here?limit=10"},
	{10, PaginationQueryParameters{"created_at"}, "http://localhost:8080/some/path/here?limit=10&sort=created_at"},
	{10, PaginationQueryParameters{}, "http://localhost:8080/some/path/here?limit=10"},
}

func TestMakeFirstHrefWhenQueryObjectIsPassedIn(t *testing.T) {
	for _, testInput := range flagtestsForTestMakeFirstHrefWhenQueryIsPassedIn {
		//c, _ := apitestutils.CreateTestContext("", nil, "")
		r := MakeFirstHref(testInput.limit, "/some/path/here", testInput.queryParams)
		assert.Equal(t, testInput.expectedURL, r)
	}
}

func TestMakeNextHref(t *testing.T) {
	//c, _ := apitestutils.CreateTestContext("", nil, "")
	r := MakeNextHref(77, "starthere", "/some/path/here")
	assert.Equal(t, "http://localhost:8080/some/path/here?start=starthere&limit=77", r)
}

var flagtestsForTestMakeNextHrefWhenQueryIsPassedIn = []struct {
	limit       int
	start       string
	queryParams PaginationQueryParameters
	expectedURL string
}{
	{10, "EXAMPLE-ID", PaginationQueryParameters{"name"}, "http://localhost:8080/some/path/here?start=EXAMPLE-ID&limit=10"},
	{10, "EXAMPLE-ID", PaginationQueryParameters{"created_at"}, "http://localhost:8080/some/path/here?start=EXAMPLE-ID&limit=10&sort=created_at"},
	{10, "EXAMPLE-ID", PaginationQueryParameters{}, "http://localhost:8080/some/path/here?start=EXAMPLE-ID&limit=10"},
}

func TestMakeNextHrefWhenQueryObjectIsPassedIn(t *testing.T) {
	for _, testInput := range flagtestsForTestMakeNextHrefWhenQueryIsPassedIn {
		//c, _ := apitestutils.CreateTestContext("", nil, "")
		r := MakeNextHref(testInput.limit, testInput.start, "/some/path/here", testInput.queryParams)
		assert.Equal(t, testInput.expectedURL, r)
	}
}
