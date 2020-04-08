module github.com/itglobal/backupmonitor/pkg/storage

replace github.com/itglobal/backupmonitor/pkg/component => ../component

replace github.com/itglobal/backupmonitor/pkg/model => ../model

replace github.com/itglobal/backupmonitor/pkg/util => ../util

go 1.14

require (
	github.com/go-ini/ini v1.55.0 // indirect
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/spf13/viper v1.6.2
)
