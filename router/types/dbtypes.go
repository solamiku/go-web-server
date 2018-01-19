package types

const (
	TAB_USER = "user"
)

type DBUser struct {
	Uid      uint
	Username string
	Passwd   string
	Power    uint64
}
