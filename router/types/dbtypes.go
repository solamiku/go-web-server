package types

const (
	TAB_USER        = "user"
	TAB_DBFLUSHCFG  = "dbflushcfg"
	TAB_DBFLUSHTMPL = "dbflushtmpl"
)

type DBUser struct {
	Uid      uint
	Username string
	Passwd   string
	Power    uint64
}

type Dbflushcfg struct {
	Id   int `xorm:"pk autoincr""`
	Info string
	Dest string
}

type Dbflushtmpl struct {
	Id     int `xorm:"pk autoincr"`
	Info   string
	Tmpl   string
	Affect string
}
