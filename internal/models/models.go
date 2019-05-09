package models

type Article struct {
	id     int
	cat    int
	access int
	hidden int
	header string
	pre    string
	text   string
	date   int64
	photo  string
	views  int
	yandex int
}
