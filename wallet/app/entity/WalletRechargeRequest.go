package entity

type WalletRechargeRequest struct {
	WALLET_ID int     `json:"wallet_id"`
	Region    string  `json:"region"`
	AMOUNT    float64 `json:"amount"`
}
