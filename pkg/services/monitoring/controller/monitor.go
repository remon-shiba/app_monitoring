package ctrlMonitor

import (
	"app_monitor/pkg/config"
	scpMonitor "app_monitor/pkg/services/monitoring/script"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
)

func HealthCheckStatus() {
	var wg sync.WaitGroup
	isRunning := config.IsAppOn
	if isRunning {
		ctr := 0
		fmt.Println("HEALTH CHECK STATUS: ", time.Now().Format("2006-01-02 15:04:05"))
		for {
			wg.Add(1)
			fmt.Println("COUNTER: ", ctr)
			// SEND REQUEST
			go CheckStatusV3(config.AppList[ctr].Url, ctr, time.Second*5)
			defer wg.Done()
			ctr++
			if ctr == len(config.AppList) {
				break
			}
		}
		wg.Wait()
		// time.Sleep(time.Second * 1)
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
		// start := time.Now()
		resp, err := http.Get(url)
		// duration := time.Since(start)

		config.Mu.Lock()

		if err == nil {
			if resp.StatusCode >= 400 {
				fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, http.StatusText(resp.StatusCode))
			} else {
				fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, http.StatusText(resp.StatusCode))
			}
		} else {
			fmt.Printf("APP ID: %v | NAME: %v | URL: %v | STATUS: %v\n", config.AppList[ctr].AppId, config.AppList[ctr].Name, config.AppList[ctr].Url, "DOWN")
		}

		config.Mu.Unlock()
	}
}

func CheckStatusV3(url string, ctr int, interval time.Duration) {
	for {
		// start := time.Now()

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("ERR:", err)
		}

		client := &http.Client{
			Timeout: time.Second * time.Duration(5),
		}

		// Send the request
		resp, err := client.Do(req)
		// duration := time.Since(start)
		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				fmt.Printf("%v - APP ID: %v | STATUS: %v | ERROR: %v\n", time.Now().Format("2006-01-02 15:04:05"),
					config.AppList[ctr].AppId, "TIMEOUT", err.Error())
			}
			if strings.Contains(err.Error(), " No connection") {
				fmt.Printf("%v - APP ID: %v | STATUS: %v | ERROR: %v\n", time.Now().Format("2006-01-02 15:04:05"),
					config.AppList[ctr].AppId, "NO CONNECTION", err.Error())
			}
		} else {

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("error reading response body:", err)
			}
			resp.Body.Close()

			// Output the result
			status := "DOWN"
			if err == nil {

				if resp.StatusCode >= 400 {
					status = fmt.Sprintf("%d", resp.StatusCode)
				} else {
					status = "RUNNING"
				}

			}

			fmt.Printf("%v - APP ID: %v | STATUS: %v | RESPONSE: %v\n", time.Now().Format("2006-01-02 15:04:05"),
				config.AppList[ctr].AppId, status, string(body))
		}

	}
}
