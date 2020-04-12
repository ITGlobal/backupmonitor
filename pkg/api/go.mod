module github.com/itglobal/backupmonitor/pkg/api

go 1.14

replace github.com/itglobal/backupmonitor/pkg/database => ../database

replace github.com/itglobal/backupmonitor/pkg/component => ../component

replace github.com/itglobal/backupmonitor/pkg/service => ../service

replace github.com/itglobal/backupmonitor/doc => ../../doc

replace github.com/itglobal/backupmonitor/pkg/model => ../model

replace github.com/itglobal/backupmonitor/pkg/storage => ../storage

replace github.com/itglobal/backupmonitor/pkg/util => ../util

require (
	github.com/gin-gonic/gin v1.6.2
	github.com/itglobal/backupmonitor/doc v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/spf13/viper v1.6.2
	github.com/swaggo/gin-swagger v1.2.0
)
