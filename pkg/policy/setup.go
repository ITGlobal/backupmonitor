package policy

import (
	"github.com/itglobal/backupmonitor/pkg/component"
)

// Setup configures package services
func Setup(builder component.Builder) {
	builder.AddComponent(createRetentionPolicy)
	builder.AddComponent(createNotificationPolicy)
}
