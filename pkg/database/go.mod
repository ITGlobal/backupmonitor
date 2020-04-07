module github.com/itglobal/backupmonitor/pkg/database

go 1.14

replace github.com/itglobal/backupmonitor/pkg/model => ../model

replace github.com/itglobal/backupmonitor/pkg/component => ../component

require (
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000
	github.com/jinzhu/gorm v1.9.12
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/spf13/viper v1.6.2
)
