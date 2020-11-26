package utils

type Status struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Failed  []string `json:"failed"`
}
