package config

type NacosLoader struct {
	Path string `json:"path"`
}

func (nl *NacosLoader) Load() (*Conf, error) {
	return nil, nil
}
