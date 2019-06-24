package rsd

type Service struct {
	Type   string            `json:"type"`
	Job    string            `json:"job"`
	URL    string            `json:"url"`
	Labels map[string]string `json:"labels"`
}

type Services struct {
	Services []Service `json:"services"`
}

type SDConfigs []SDConfig

type SDConfig struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
