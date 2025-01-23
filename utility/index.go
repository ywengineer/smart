package utility

// RdbProperties rational database configuration properties
type RdbProperties struct {
	Name       string           `json:"name" yaml:"name"` // mysql or postgres
	Username   string           `json:"username" yaml:"username"`
	Password   string           `json:"password" yaml:"password"`
	Host       string           `json:"host" yaml:"host"`
	Port       int              `json:"port" yaml:"port"`
	Database   string           `json:"database" yaml:"database"`
	Parameters string           `json:"parameters" yaml:"parameters"`
	Pool       DbPoolProperties `json:"pool" yaml:"pool"`
}

type DbPoolProperties struct {
	MaxIdleCon          int   `json:"max_idle_con" yaml:"max-idle-con"`
	MaxOpenCon          int   `json:"max_open_con" yaml:"max-open-con"`
	MaxLifeTimeInMinute int64 `json:"max_life_time_minute" yaml:"max-life-time-minute"`
}
