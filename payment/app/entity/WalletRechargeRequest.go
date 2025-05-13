package entity

type WalletRechargeRequest struct {
	WALLET_ID int     `json:"wallet_id"`
	AMOUNT    float64 `json:"amount"`
}
