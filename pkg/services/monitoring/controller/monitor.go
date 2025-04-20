package ctrlMonitor

import (
	"app_monitor/pkg/config"
	mdlMonitor "app_monitor/pkg/services/monitoring/model"
	scpMonitor "app_monitor/pkg/services/monitoring/script"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
)

func HealthCheckStatus() {
	isRunning := config.IsAppOn
	if isRunning {
		ctr := 0
		for {
			// SEND REQUEST
			go CheckStatusV3(config.AppList[ctr].Url, ctr, time.Second*5)
			ctr++
			if ctr == len(config.AppList) {
				break
			}
		}
		time.Sleep(time.Second * 5)
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

func CheckStatusV2(url string, ctr int, interval time.Duration) {
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
	}
}

func CheckStatusV3(url string, ctr int, interval time.Duration) {
	for {
		start := time.Now()

		var app mdlMonitor.AppDetails

		// Lock only while reading shared config
		config.Mu.RLock()
		if ctr >= 0 && ctr < len(config.AppList) {
			app = config.AppList[ctr]
		}
		config.Mu.RUnlock()

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("ERR:", err)
		}

		client := &http.Client{
			Timeout: time.Second * time.Duration(5),
		}

		// Send the request
		resp, err := client.Do(req)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v | DURATION: %v\n",
				app.AppId, app.Name, app.Url, "TIMEOUT", duration)
		} else {
			defer resp.Body.Close()

			// Output the result
			status := "DOWN"
			if err == nil {

				if resp.StatusCode >= 400 {
					status = fmt.Sprintf("%d", resp.StatusCode)
				} else {
					status = "RUNNING"
				}

			}

			fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v | DURATION: %v\n",
				app.AppId, app.Name, app.Url, status, duration)
		}

	}
}
