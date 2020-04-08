module github.com/itglobal/backupmonitor

go 1.14

replace github.com/itglobal/backupmonitor/doc => ./doc

replace github.com/itglobal/backupmonitor/pkg/api => ./pkg/api

replace github.com/itglobal/backupmonitor/pkg/service => ./pkg/service

replace github.com/itglobal/backupmonitor/pkg/component => ./pkg/component

replace github.com/itglobal/backupmonitor/pkg/model => ./pkg/model

replace github.com/itglobal/backupmonitor/pkg/database => ./pkg/database

replace github.com/itglobal/backupmonitor/pkg/storage => ./pkg/storage

replace github.com/itglobal/backupmonitor/pkg/util => ./pkg/util

replace github.com/itglobal/backupmonitor/pkg/policy => ./pkg/policy

replace github.com/itglobal/backupmonitor/pkg/notify => ./pkg/notify

require (
	github.com/gin-gonic/gin v1.6.2 // indirect
	github.com/go-ini/ini v1.55.0 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/hackebrot/turtle v0.1.0 // indirect
	github.com/itglobal/backupmonitor/doc v0.0.0-00010101000000-000000000000 // indirect
	github.com/itglobal/backupmonitor/pkg/api v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/database v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000 // indirect
	github.com/itglobal/backupmonitor/pkg/notify v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/policy v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/storage v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/m1/go-generate-password v0.0.0-20191114193340-84682ecbc3fd // indirect
	github.com/minio/minio-go v6.0.14+incompatible // indirect
	github.com/slack-go/slack v0.6.3 // indirect
	github.com/spf13/viper v1.6.2
	github.com/swaggo/gin-swagger v1.2.0 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
