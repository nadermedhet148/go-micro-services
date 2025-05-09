package externalServices

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type NotificationService struct{}

type NotificationRequest struct {
	REF_NUMBER string `json:"ref_number"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}

func (ps *NotificationService) SendNotification(NotificationReq NotificationRequest) (NotificationResponse, error) {

	externalReq, err := json.Marshal(NotificationReq)
	if err != nil {
		return NotificationResponse{}, err
	}

	externalServiceURL := os.Getenv("Notification_SERVICE_URL")
	req, err := http.NewRequest("POST", externalServiceURL+"/api/v1/Notification", bytes.NewBuffer(externalReq))
	if err != nil {
		return NotificationResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return NotificationResponse{}, err

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NotificationResponse{}, err
	}

	return NotificationResponse{
		Status: "success",
	}, nil

}
