module github.com/itglobal/backupmonitor/pkg/storage

replace github.com/itglobal/backupmonitor/pkg/component => ../component

replace github.com/itglobal/backupmonitor/pkg/model => ../model

replace github.com/itglobal/backupmonitor/pkg/util => ../util

go 1.14

require (
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/util v0.0.0-00010101000000-000000000000
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/spf13/viper v1.6.2
)
