package repository

type Database interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
