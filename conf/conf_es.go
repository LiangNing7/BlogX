package conf

import "fmt"

type ES struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	IsHttps  bool   `yaml:"is_https"`
}

func (e ES) Url() string {
	if e.IsHttps {
		return fmt.Sprintf("https://%s", e.Addr)
	}
	return fmt.Sprintf("http://%s", e.Addr)
}
