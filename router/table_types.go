package router

const (
	TAB_USER = "user"
)

type DBUser struct {
	Uid      uint
	Username string
	Password string
	Power    uint64
}
