package structs

// Status is the struct for the status of the MAAS plugin
type Status struct {
	Running      bool   `json:"running"`
	Failed       bool   `json:"failed"`
	ErrDetails   string `json:"err_details"`
	MachineCount int    `json:"machine_count"`
}
