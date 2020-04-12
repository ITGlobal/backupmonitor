module github.com/itglobal/backupmonitor/pkg/policy

replace github.com/itglobal/backupmonitor/pkg/component => ../component

replace github.com/itglobal/backupmonitor/pkg/database => ../database

replace github.com/itglobal/backupmonitor/pkg/model => ../model

replace github.com/itglobal/backupmonitor/pkg/service => ../service

replace github.com/itglobal/backupmonitor/pkg/storage => ../storage

replace github.com/itglobal/backupmonitor/pkg/util => ../util

replace github.com/itglobal/backupmonitor/pkg/notify => ../notify

go 1.14

require (
	github.com/itglobal/backupmonitor/pkg/component v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/database v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/model v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/notify v0.0.0-00010101000000-000000000000
	github.com/itglobal/backupmonitor/pkg/service v0.0.0-00010101000000-000000000000
	github.com/sarulabs/di v2.0.0+incompatible
)
