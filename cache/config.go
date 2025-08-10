package cache

// TableConfig represents the configuration for database tables and their caching structure
type TableConfig struct {
	MasterTables []MasterTable `json:"MasterTables"`
	SlaveTables  []SlaveTable  `json:"SlaveTables"`
}

// MasterTable represents a master table configuration
type MasterTable struct {
	ToDB        []string `json:"toDB"`
	SqlName     string   `json:"SqlName"`
	SelectCondi string   `json:"SelectCondi"`
	Statdir     string   `json:"Statdir"`
	Keyarray    []string `json:"Keyarray"`
	Statidx     []string `json:"Statidx"`
	Slaves      []string `json:"Slaves"`
}

// SlaveTable represents a slave table configuration  
type SlaveTable struct {
	ToDB        []string `json:"toDB"`
	SqlName     string   `json:"SqlName"`
	SelectCondi string   `json:"SelectCondi"`
	Statdir     string   `json:"Statdir"`
	Keyarray    []string `json:"Keyarray"`
	Statidx     []string `json:"Statidx"`
	Slaves      []string `json:"Slaves"`
}