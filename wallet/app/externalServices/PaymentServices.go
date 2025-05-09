package externalServices

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type PaymentService struct{}

type PaymentRequest struct {
	REF_NUMBER string `json:"ref_number"`
}

type PaymentResponse struct {
	Status string `json:"status"`
}

func (ps *PaymentService) SendPayment(paymentReq PaymentRequest) (PaymentResponse, error) {

	externalReq, err := json.Marshal(paymentReq)
	if err != nil {
		return PaymentResponse{}, err
	}

	externalServiceURL := os.Getenv("PAYMENT_SERVICE_URL")
	req, err := http.NewRequest("POST", externalServiceURL+"/api/v1/Payment", bytes.NewBuffer(externalReq))
	if err != nil {
		return PaymentResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PaymentResponse{}, err

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PaymentResponse{}, err
	}

	return PaymentResponse{
		Status: "success",
	}, nil

}
