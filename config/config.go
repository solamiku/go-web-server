package config

import (
	"encoding/xml"
	"fmt"
	"os"
)

//Config strcut
type Config struct {
	Server struct {
		Listen     string `xml:"listen"`
		Compress   int    `xml:"compress"`
		EncryptKey string `xml:"encrypt_key"`
		PaddingKey string `xml:"padding_key"`
		Https      struct {
			Open int    `xml:"open"`
			Crt  string `xml:"crt"`
			Key  string `xml:"key"`
		} `xml:"https"`
	} `xml:"server"`
	Db struct {
		Src     string `xml:"src"`
		MaxConn int    `xml:"max_conn"`
	} `xml:"db"`
	Template struct {
		Views []struct {
			Src        string   `xml:"src,attr"`
			Components []string `xml:"component"`
		} `xml:"views>view"`
	} `xml:"template"`
}

//G export global config
var G *Config

// LoadCfg load config file to export var Config
func LoadCfg(fnconfig string) error {
	local := Config{}

	f, err := os.Open(fnconfig)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = xml.NewDecoder(f).Decode(&local); err != nil {
		err = fmt.Errorf("(%s)%v", fnconfig, err)
		return err
	}

	G = &local

	return nil
}
