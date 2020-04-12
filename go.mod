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
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hackebrot/go-repr v0.1.0 // indirect
	github.com/itglobal/backupmonitor/pkg/api v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/database v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/notify v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/policy v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/storage v0.0.0-00010101000000-000000000000
	github.com/smartystreets/assertions v0.0.0-20190116191733-b6c0e53d7304 // indirect
	github.com/spf13/viper v1.6.2
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
