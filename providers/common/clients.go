package common

const (
	// Providers
	TGProvider ProviderType = "tg"
	VKProvider ProviderType = "vk"
)

type ProviderType string

type Message struct {
	Provider ProviderType
	ChatID   int
	Content  string
}

type Provider interface {
	SendMessage(Message) error
	GetMessages() <-chan Message
}
