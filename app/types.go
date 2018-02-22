package app

type User struct {
	Id            uint64
	Username      string
	Rent          uint64
	Wealth        uint64
	EncryptedPass string
}

type Budget struct {
	Id     uint64
	UserId uint64
	Other  string
}
