package tests

import (
	"os"
	"testing"

	"gcw/service"
	"github.com/joho/godotenv"
)

func TestStellarRPC(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Log("No .env file found, using system environment variables")
	}

	secret := os.Getenv("STELLAR_SECRET_KEY")
	contract := os.Getenv("STELLAR_CONTRACT_ADDRESS")
	rpc := os.Getenv("STELLAR_RPC_URL")
	network := os.Getenv("STELLAR_NETWORK")

	t.Logf("STELLAR_SECRET_KEY: %s", secret)
	t.Logf("STELLAR_CONTRACT_ADDRESS: %s", contract)
	t.Logf("STELLAR_RPC_URL: %s", rpc)
	t.Logf("STELLAR_NETWORK: %s", network)

	if secret == "" || contract == "" {
		t.Skip("Stellar configuration missing, skipping test")
	}

	stellarSvc := service.NewStellarService()
	balance, err := stellarSvc.GetAccountBalance()
	if err != nil {
		t.Fatalf("GetAccountBalance failed: %v", err)
	}
	t.Logf("Account Balance: %s", balance)

	hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // SHA-256 of empty string
	txHash, author, err := stellarSvc.SendAuditHash(hash)
	if err != nil {
		t.Fatalf("SendAuditHash failed: %v", err)
	}

	t.Logf("Success! TxHash: %s, Author: %s", txHash, author)
}
