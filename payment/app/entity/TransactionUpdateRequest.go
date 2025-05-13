package entity

type TransactionUpdateRequest struct {
	REF_NUMBER string `json:"ref_number"`
	STATUS     string `json:"status"`
}
