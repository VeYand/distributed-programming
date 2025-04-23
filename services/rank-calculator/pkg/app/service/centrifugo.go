package service

type CentrifugoClient interface {
	Publish(channel string, data interface{}) error
}
