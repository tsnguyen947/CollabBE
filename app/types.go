package app

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

type User struct {
	Id       uint64
	Username string
	Rent     uint64
	Wealth   uint64
}
