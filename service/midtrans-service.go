package service

import (
	"log"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransService struct {
	Client coreapi.Client
}

func NewMidtransService() *MidtransService {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	envStr := os.Getenv("MIDTRANS_ENVIRONMENT")

	envType := midtrans.Sandbox
	if envStr == "production" {
		envType = midtrans.Production
	}

	var client coreapi.Client
	client.New(serverKey, envType)

	return &MidtransService{
		Client: client,
	}
}

// GenerateQRIS creates a new Charge request specifically for QRIS.
// Returns the qr_string or an error.
func (s *MidtransService) GenerateQRIS(orderID string, amount int64) (string, error) {
	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomExpiry: &coreapi.CustomExpiry{
			ExpiryDuration: 60, // 60 minutes expiry
			Unit:           "minute",
		},
	}

	res, err := s.Client.ChargeTransaction(req)
	if err != nil {
		log.Println("Midtrans Charge Error:", err)
		return "", err
	}

	// Midtrans returns the QR string in the Actions array for GoPay/QRIS
	if len(res.Actions) > 0 {
		for _, action := range res.Actions {
			if action.Name == "generate-qr-code" {
				return action.URL, nil
			}
		}
	}

	// Fallback to QRString if populated
	if res.QRString != "" {
		return res.QRString, nil
	}

	return "", nil
}
