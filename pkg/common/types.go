package common

type MessageHandler func(data []byte)

type Service interface {
	Close()
}
