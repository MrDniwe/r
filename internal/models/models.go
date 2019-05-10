package models

type Article struct {
	Id     int
	Cat    int
	Access int
	Hidden int
	Header string
	Pre    string
	Text   string
	Date   int64
	Photo  string
	Views  int
	Yandex int
}
