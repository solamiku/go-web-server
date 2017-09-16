package config

import (
	"encoding/xml"
	"fmt"
	"os"
)

//Config strcut
type Config struct {
	Server struct {
		Listen   string `xml:"listen"`
		Compress int    `xml:"compress"`
	} `xml:"server"`
	Db struct {
		Src     string `xml:"src"`
		MaxConn int    `xml:"max_conn"`
	} `xml:"db"`
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
