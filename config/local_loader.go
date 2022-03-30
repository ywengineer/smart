package config

type LocalLoader struct {
	Path string `json:"path"`
}

func (ll *LocalLoader) Load() (*Conf, error) {
	return nil, nil
}
