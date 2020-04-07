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

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/aws/aws-sdk-go v1.30.6 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gin-gonic/gin v1.6.2 // indirect
	github.com/go-ini/ini v1.55.0 // indirect
	github.com/go-openapi/spec v0.19.7 // indirect
	github.com/go-openapi/swag v0.19.8 // indirect
	github.com/go-swagger/go-swagger v0.23.0 // indirect
	github.com/itglobal/backupmonitor/pkg/api v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/database v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000 // indirect
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/storage v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/m1/go-generate-password v0.0.0-20191114193340-84682ecbc3fd // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/minio/minio-go v6.0.14+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rlmcpherson/s3gof3r v0.5.0 // indirect
	github.com/sarulabs/di v2.0.0+incompatible
	github.com/spf13/viper v1.6.2
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14 // indirect
	github.com/swaggo/gin-swagger v1.2.0 // indirect
	github.com/swaggo/swag v1.6.5
	github.com/urfave/cli v1.22.4 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/tools v0.0.0-20200403190813-44a64ad78b9b // indirect
	gopkg.in/tylerb/graceful.v1 v1.2.15 // indirect
)
