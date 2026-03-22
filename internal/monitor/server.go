package monitor

type Server struct{
	Name string `json:"name"`
	Host string `json:"host"`
	Type string `json:"type"`
	Interval int `json:"interval"`
	Timeout int `json:"timeout"`
	Enabled bool `json:"enabled"`
}