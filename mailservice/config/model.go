package config

type SiteConfig struct {
	ReceiverAddress string
	CorsOrigin      string
	Fields          []string
}

type ServiceConfig struct {
	SenderAddress  string
	SenderPassword string
	ServerName     string
	ServerPort     uint16
	Sites          map[string]SiteConfig
}