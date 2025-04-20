package config

import (
	mdlMonitor "app_monitor/pkg/services/monitoring/model"
	"sync"
)

var (
	AppS3cr3tK3y = "c9UvsjTg7BkPIw8ucByygvMfRx0XtDN5"
	IsAppOn      = true
	AppList      []mdlMonitor.AppDetails
	Mu           sync.RWMutex
)
