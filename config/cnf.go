package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/dop251/goja"
	"github.com/ghodss/yaml"
)

type Config struct {
	Chunk  *Chunk  `json:"chunk"`
	Sec    *Sec    `json:"sec"`
	Custom *Custom `json:"custom"`
}

type Chunk struct {
	A1B string `json:"a1b"`
	BaseChunk
}

type Sec struct {
	At         string `json:"at"`
	RefererNum int    `json:"refererNum"`
	Domain     string `json:"domain"`
	Domain1122 string `json:"domain1122"`
	JsonIdUrl  string `json:"jsonIdUrl"`
	DesKey     string `json:"desKey"`
	BaseChunk
}

type BaseChunk struct {
	Path         string `json:"path"`
	ClientKey    string `json:"clientKey"`
	ClientSecret string `json:"clientSecret"`
}

type Custom struct {
	Tokens      []string `json:"tokens"`
	PreDuration int      `json:"preDuration"`
	SecTime     int      `json:"secTime"`
	LogPath     string   `json:"logPath"`
}

func Load() (cnf *Config, err error) {
	all, err := ioutil.ReadFile("./config/cnf.yaml")
	if err != nil {
		return
	}
	b, err := yaml.YAMLToJSON(all)
	if err != nil {
		return
	}

	cnf = new(Config)

	if err = json.Unmarshal(b, cnf); err != nil {
		return
	}

	if err = load(cnf); err != nil {
		return
	}
	return
}

func load(c *Config) error {
	if len(c.Custom.Tokens) == 0 {
		return errors.New("请在配置文件中custom-tokens中添加token")
	}
	if c.Sec.ClientKey == "" {
		key := js(c.Chunk.A1B, c.Chunk.ClientKey)
		if key == "" {
			return errors.New("clientKey是空")
		}
		c.Sec.ClientKey = key
	}

	if c.Sec.ClientSecret == "" {
		secret := js(c.Chunk.A1B, c.Chunk.ClientSecret)
		if secret == "" {
			return errors.New("clientSecret是空")
		}
		c.Sec.ClientSecret = secret
	}

	if c.Sec.Path == "" {
		path := js(c.Chunk.A1B, c.Chunk.Path)
		if path == "" {
			return errors.New("path是空")
		}
		c.Sec.Path = path
	}
	return nil
}

func js(a1b, key string) string {
	if a1b == "" || key == "" {
		return ""
	}
	goJs := goja.New()

	runString, err := goJs.RunString(a1b + " " + key)
	if err != nil {
		return ""
	}
	return runString.String()
}
