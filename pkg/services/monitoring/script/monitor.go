package scpMonitor

import (
	"app_monitor/pkg/config"
	mdlMonitor "app_monitor/pkg/services/monitoring/model"
)

func GetAllApps() []mdlMonitor.AppDetails {
	appList := []mdlMonitor.AppDetails{}
	config.DBConnList[0].Raw("SELECT * FROM APP_DETAILS").Scan(&appList)
	return appList
}
