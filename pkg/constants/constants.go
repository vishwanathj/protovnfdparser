package constants

// ENVCONFIG_PREFIX is the value for prefix to be used with envconfig module
const ENVCONFIG_PREFIX = "VNFDSVC"

// VnfdKey the key to be used for the VNFD document collection
const VnfdKey = "name"

// VnfdID the key that holds the Vnfd ID
const VnfdID = "_id"

// IDPrefix VNFD ID prefix
const VnfdIDPrefix = "VNFD-"

// Available the vnfd status at creation time
const Available = "available"

// VnfdIDPattern regexp pattern for the VnfdId
const VnfdIDPattern = "^VNFD-[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"

// PaginationURLSort
const PaginationURLSort = "sort"

// PaginationURLLimit
const PaginationURLLimit = "limit"

// PaginationURLStart
const PaginationURLStart = "start"

const ApiPathVnfds = "/vnfds"