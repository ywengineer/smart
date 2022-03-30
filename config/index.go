package config

type Conf struct {
	Network           string `json:"network" yaml:"network"`
	Address           string `json:"address" yaml:"address"`
	Workers           int    `json:"workers" yaml:"workers"`
	WorkerLoadBalance int    `json:"load_balance" yaml:"load_balance"`
}
