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
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/itglobal/backupmonitor/pkg/api v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/database v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/notify v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/policy v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/storage v0.0.0-00010101000000-000000000000
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/smartystreets/assertions v0.0.0-20190116191733-b6c0e53d7304 // indirect
	github.com/spf13/viper v1.6.2
	github.com/swaggo/swag v1.6.7 // indirect
	github.com/urfave/cli v1.22.4 // indirect
	github.com/urfave/cli/v2 v2.2.0 // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200701151220-7cb253f4c4f8 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
