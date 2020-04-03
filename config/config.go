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
		Cookie struct {
			SessionId string `xml:"session_id"`
			AutoLogin string `xml:"autologin"`
		} `xml:"cookie"`
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
	Servers  []Server `xml:"servers>server"`
	JumpList []struct {
		Name string `xml:"name,attr"`
		Auth int    `xml:"auth,attr"`
		Opts []struct {
			Name string `xml:"name,attr"`
			Url  string `xml:",chardata"`
		} `xml:"opt"`
	} `xml:"jumplist>grp"`
	DetectMax int `xml:"servers>detect_max"`
	Power     []struct {
		U    int    `xml:"u,attr"`
		Desc string `xml:"desc,attr"`
	} `xml:"powers>power"`
	LeitingLog []LeitingLogConf `xml:"leitinglog"`
	DBFlush    struct {
		Open int `xml:"open"`
	} `xml:"dbflush"`
}

type LeitingLogConf struct {
	Id               string `xml:"id,attr"`
	Auth             int    `xml:"auth,attr"`
	Cmd              string `xml:"cmd"`
	LogPath          string `xml:"logpath"`
	ErrCmd           string `xml:"errcmd"`
	ErrPath          string `xml:"errpath"`
	ErrInterval      int    `xml:"err_interval"`
	ErrForceInterval int    `xml:"err_force_interval"`
	Err              []struct {
		Dir  string `xml:"dir,attr"`
		Game string `xml:"game,attr"`
	} `xml:"err"`
}

type Server struct {
	Name         string    `xml:"name,attr"`
	GamesrvAdmin string    `xml:"gamesrv_admin,attr"`
	ZonesrvAdmin string    `xml:"zonesrv_admin,attr"`
	Commands     []Command `xml:"command"`
}

type Command struct {
	Tag     string `xml:"tag,attr"`
	Name    string `xml:"name,attr"`
	Auth    int    `xml:"auth,attr"`
	Clientp string `xml:"clientp,attr"`
	Args    []struct {
		Sleep  int          `xml:"sleep,attr"`
		Output int          `xml:"output,attr"`
		Type   string       `xml:"type,attr"`
		Desc   string       `xml:"desc,attr"`
		Val    string       `xml:"v"`
		Extra  CmdArgExtras `xml:"extra"`
		Args   CmdArgExtras `xml:"args"`
	} `xml:"args"`
}

type CmdArgExtras []string

func (cae CmdArgExtras) GetVal(idx int) string {
	if len(cae) > idx {
		return cae[idx]
	}
	return ""
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
