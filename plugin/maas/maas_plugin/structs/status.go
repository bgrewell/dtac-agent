package structs

type Status struct {
	Running      bool   `json:"running"`
	Failed       bool   `json:"failed"`
	ErrDetails   string `json:"err_details"`
	MachineCount int    `json:"machine_count"`
}
