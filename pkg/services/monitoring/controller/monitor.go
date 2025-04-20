package ctrlMonitor

import (
	"app_monitor/pkg/config"
	scpMonitor "app_monitor/pkg/services/monitoring/script"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
)

func HealthCheckStatus() {
	time.Sleep(5 * time.Second)
	isRunning := config.IsAppOn
	if isRunning {
		ctr := 0
		for {
			// SEND REQUEST
			go CheckStatusV2(config.AppList[ctr].Url, ctr)
			ctr++
			if ctr == len(config.AppList) {
				break
			}
		}
		HealthCheckStatus()
	}
}

func GetAllRegisteredApps() bool {
	config.AppList = scpMonitor.GetAllApps()
	return len(config.AppList) > 0
}

func CheckStatus(url string, ctr int) {
	// SEND REQUEST
	response, status, err := utils_v1.SendRequestWithCode(config.AppList[ctr].Url, "GET", []byte(""), map[string]string{}, 5)
	marResp, _ := json.Marshal(response)
	if response != nil {
		fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, "RUNNING")
	}
	if err != nil {
		fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, "DOWN")
	}

	fmt.Println("RESPONSE:", string(marResp))
	fmt.Println("RESPONSE STATUS:", *status)
}

func CheckStatusV2(url string, ctr int) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		config.Mu.Lock()

		if err == nil {
			if resp.StatusCode >= 400 {
				fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v | DURATION: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, resp.StatusCode, duration)
			} else {
				fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v | DURATION: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, "RUNNING", duration)
			}
		} else {
			fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v | DURATION: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, "DOWN", duration)
		}

		config.Mu.Unlock()
		<-ticker.C
	}
}
