// +build unit

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package models_test

import (
	"fmt"
	"github.com/vishwanathj/protovnfdparser/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
)

var testcfg = config.GetConfigInstance()

func TestMakeFirstHref(t *testing.T) {
	pgCfg := testcfg.PgntConfig
	baseURI := pgCfg.BaseURI
	tests := []struct {
		desc        string
		limit       int
		apipath     string
		query       models.PaginationQueryParameters
		expectedURL string
	}{
		{
			desc:        "test with limit only and no query params",
			limit:       42,
			apipath:     "/api/v1/vnfds",
			expectedURL: baseURI + "/api/v1/vnfds?limit=42",
		},
		{
			desc:        "test with limit and `DefaultOrderBy` query",
			limit:       10,
			apipath:     "/api/v1/vnfds",
			query:       models.PaginationQueryParameters{OrderBy: "name"},
			expectedURL: baseURI + "/api/v1/vnfds?limit=10",
		},
		{
			desc:        "test with limit and non ``DefaultOrderBy` query",
			limit:       10,
			apipath:     "/api/v1/vnfds",
			query:       models.PaginationQueryParameters{OrderBy: "created_at"},
			expectedURL: baseURI + "/api/v1/vnfds?limit=10&sort=created_at",
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			r := models.MakeFirstHref(*pgCfg, tc.limit, tc.apipath, tc.query)
			assert.Equal(t, tc.expectedURL, r)
		})
	}
}

func TestMakeNextHref(t *testing.T) {
	pgCfg := testcfg.PgntConfig
	baseURI := pgCfg.BaseURI
	tests := []struct {
		desc        string
		limit       int
		start       string
		apipath     string
		query       models.PaginationQueryParameters
		expectedURL string
	}{
		{
			desc:        "test with limit, start and no query params",
			limit:       42,
			start:       "EXAMPLE-ID",
			apipath:     "/api/v1/vnfds",
			expectedURL: baseURI + "/api/v1/vnfds?start=EXAMPLE-ID&limit=42",
		},
		{
			desc:        "test with limit, start and `DefaultOrderBy` query",
			limit:       10,
			start:       "VNFD-NAME",
			apipath:     "/api/v1/vnfds",
			query:       models.PaginationQueryParameters{OrderBy: "name"},
			expectedURL: baseURI + "/api/v1/vnfds?start=VNFD-NAME&limit=10",
		},
		{
			desc:        "test with limit and non ``DefaultOrderBy` query",
			limit:       10,
			start:       "VNFD-ID",
			apipath:     "/api/v1/vnfds",
			query:       models.PaginationQueryParameters{OrderBy: "created_at"},
			expectedURL: baseURI + "/api/v1/vnfds?start=VNFD-ID&limit=10&sort=created_at",
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tc.desc), func(t *testing.T) {
			r := models.MakeNextHref(*pgCfg, tc.limit, tc.start, tc.apipath, tc.query)
			assert.Equal(t, tc.expectedURL, r)
		})
	}
}
