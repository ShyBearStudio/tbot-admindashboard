package tbot

type TBot interface {
	Start() error
	Stop() error
}
