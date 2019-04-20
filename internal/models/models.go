package models

type Article struct {
	active bool   `json:"active" bson:"active"`
	title  string `json:"title bson:title"`
}
