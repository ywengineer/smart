package config

type Loader interface {
	Load() (*Conf, error)
}

type ValueLoader struct {
	Conf *Conf
}

func (vl *ValueLoader) Load() (*Conf, error) {
	return vl.Conf, nil
}
