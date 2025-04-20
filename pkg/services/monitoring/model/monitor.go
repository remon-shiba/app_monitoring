package mdlMonitor

import "time"

type (
	AppDetails struct {
		AppId     string    `json:"appId"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
