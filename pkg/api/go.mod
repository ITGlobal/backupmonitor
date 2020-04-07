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
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/aws/aws-sdk-go v1.30.6 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/go-ini/ini v1.55.0 // indirect
	github.com/itglobal/backupmonitor/doc v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/database v0.0.0-00010101000000-000000000000 // indirect
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/storage v0.0.0-00010101000000-000000000000 // indirect
	github.com/itglobal/backupmonitor/pkg/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/m1/go-generate-password v0.0.0-20191114193340-84682ecbc3fd // indirect
	github.com/minio/minio-go v6.0.14+incompatible // indirect
	github.com/rlmcpherson/s3gof3r v0.5.0 // indirect
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/spf13/viper v1.6.2
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
)
