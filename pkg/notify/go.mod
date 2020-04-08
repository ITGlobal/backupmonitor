module github.com/itglobal/backupmonitor/pkg/notify

go 1.14

replace github.com/itglobal/backupmonitor/pkg/component => ../component

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/hackebrot/turtle v0.1.0
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/slack-go/slack v0.6.3
	github.com/spf13/viper v1.6.2
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
