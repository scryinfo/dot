package conns

type Config struct {
	Scheme   string          `json:"scheme"`
	Services []ServiceConfig `json:"services"`
}

type ServiceConfig struct {
	Name    string   `json:"name"`
	Addrs   []string `json:"addrs"`
	Balance string   `json:"balance"` // round or first, the default value is round
}
