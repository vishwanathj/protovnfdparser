package utils

// interface defined to abstract and help with unit tests
type VnfdValidator interface {
	ValidateVnfdInstanceBody([]byte) error
	ValidateVnfdPostBody([]byte) error
	ValidatePaginatedVnfdsInstancesBody([]byte) error
}

type vnfdSchemaChecker struct{}

func (val vnfdSchemaChecker) ValidateVnfdInstanceBody(jsonval []byte) error {
	return ValidateVnfdInstanceBody(jsonval)
}

func (val vnfdSchemaChecker) ValidateVnfdPostBody(jsonval []byte) error {
	return ValidateVnfdPostBody(jsonval)
}

func (val vnfdSchemaChecker) ValidatePaginatedVnfdsInstancesBody(jsonval []byte) error {
	return ValidatePaginatedVnfdsInstancesBody(jsonval)
}

func NewVnfdValidator() VnfdValidator {
	return vnfdSchemaChecker{}
}
