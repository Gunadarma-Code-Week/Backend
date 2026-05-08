package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type StellarService interface {
	SendAuditHash(hash string) (string, string, error)
	GetAccountBalance() (string, error)
	GetPublicKey() string
}

type stellarService struct {
	horizonClient   *horizonclient.Client
	rpcURL          string
	secret          string
	network         string
	contractID      string
	functionName    string
	publicKey       string
}

// Structs for JSON-RPC
type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Transaction string `json:"transaction"`
	} `json:"params"`
}

type SimResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		TransactionData string   `json:"transactionData"`
		MinResourceFee  string   `json:"minResourceFee"`
		Events          []string `json:"events"`
		Results         []struct {
			Auth []string `json:"auth"`
			XDR  string   `json:"xdr"`
		} `json:"results"`
		Error string `json:"error"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type SendResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Hash   string `json:"hash"`
		Status string `json:"status"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewStellarService() StellarService {
	net := os.Getenv("STELLAR_NETWORK")
	if net == "" {
		net = "testnet"
	}

	rpcURL := os.Getenv("STELLAR_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://stellar-soroban-public.nodies.app"
	}

	horizonURL := os.Getenv("STELLAR_HORIZON_URL")
	if horizonURL == "" {
		horizonURL = "https://horizon.stellar.org"
		if net == "testnet" {
			horizonURL = "https://horizon-testnet.stellar.org"
		}
	}

	fmt.Printf("Stellar Service initialized: network=%s, rpc=%s, horizon=%s, contract=%s\n", 
		net, rpcURL, horizonURL, os.Getenv("STELLAR_CONTRACT_ADDRESS"))

	secret := os.Getenv("STELLAR_SECRET_KEY")
	publicKey := ""
	if secret != "" {
		kp, err := keypair.ParseFull(secret)
		if err == nil {
			publicKey = kp.Address()
		}
	}

	return &stellarService{
		horizonClient: &horizonclient.Client{HorizonURL: horizonURL},
		rpcURL:        rpcURL,
		secret:        secret,
		network:       net,
		contractID:    os.Getenv("STELLAR_CONTRACT_ADDRESS"),
		functionName:  "record_hash",
		publicKey:     publicKey,
	}
}

func (s *stellarService) GetPublicKey() string {
	return s.publicKey
}

func (s *stellarService) SendAuditHash(auditHash string) (string, string, error) {
	if s.secret == "" {
		return "", "", fmt.Errorf("STELLAR_SECRET_KEY not configured")
	}
	if s.contractID == "" {
		return "", "", fmt.Errorf("STELLAR_CONTRACT_ADDRESS not configured")
	}

	kp, err := keypair.ParseFull(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("invalid stellar secret: %w", err)
	}

	authorAddress := kp.Address()

	// 1. Ambil detail akun untuk Sequence Number via Horizon
	accountRequest := horizonclient.AccountRequest{AccountID: authorAddress}
	sourceAccount, err := s.horizonClient.AccountDetail(accountRequest)
	if err != nil {
		return "", authorAddress, fmt.Errorf("failed to fetch account detail: %w", err)
	}

	// Hash data (menggunakan auditHash yang diberikan)
	// AuditHash kita asumsikan sudah berupa hex string, kita convert balik ke bytes
	var hashData [32]byte
	fmt.Sscanf(auditHash, "%x", &hashData)
	hashBytesXDR := xdr.ScBytes(hashData[:])
	hashSCVal := xdr.ScVal{Type: xdr.ScValTypeScvBytes, Bytes: &hashBytesXDR}

	// Alamat Kontrak
	contractBytes, _ := strkey.Decode(strkey.VersionByteContract, s.contractID)
	var contractHash xdr.ContractId
	copy(contractHash[:], contractBytes)
	contractSCAddr := xdr.ScAddress{Type: xdr.ScAddressTypeScAddressTypeContract, ContractId: &contractHash}

	// 3. Bangun Transaksi Awal untuk Simulasi
	invokeOp := &txnbuild.InvokeHostFunction{
		HostFunction: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeContract: &xdr.InvokeContractArgs{
				ContractAddress: contractSCAddr,
				FunctionName:    xdr.ScSymbol(s.functionName),
				Args:            []xdr.ScVal{hashSCVal},
			},
		},
	}

	txSim, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{invokeOp},
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		return "", authorAddress, fmt.Errorf("failed to build sim tx: %w", err)
	}

	txBase64, _ := txSim.Base64()

	// 4. Simulasi Transaksi via JSON-RPC
	var simResp SimResponse
	if err := s.callRPC("simulateTransaction", txBase64, &simResp); err != nil {
		return "", authorAddress, fmt.Errorf("simulation failed: %w", err)
	}
	if simResp.Error != nil {
		return "", authorAddress, fmt.Errorf("RPC simulation error: %s", simResp.Error.Message)
	}
	if simResp.Result.Error != "" {
		return "", authorAddress, fmt.Errorf("VM simulation error: %s", simResp.Result.Error)
	}

	// 5. Ekstraksi Data Simulasi (Footprint & Auth)
	var txData xdr.SorobanTransactionData
	if err := xdr.SafeUnmarshalBase64(simResp.Result.TransactionData, &txData); err != nil {
		return "", authorAddress, fmt.Errorf("failed to unmarshal transaction data: %w", err)
	}

	minResourceFee, _ := strconv.ParseInt(simResp.Result.MinResourceFee, 10, 64)
	totalFee := int64(txnbuild.MinBaseFee) + minResourceFee

	var authEntries []xdr.SorobanAuthorizationEntry
	if len(simResp.Result.Results) > 0 {
		for _, authB64 := range simResp.Result.Results[0].Auth {
			var authEntry xdr.SorobanAuthorizationEntry
			xdr.SafeUnmarshalBase64(authB64, &authEntry)
			authEntries = append(authEntries, authEntry)
		}
	}

	// 6. Bangun Transaksi Final dengan Auth & SorobanData
	invokeOpFinal := &txnbuild.InvokeHostFunction{
		HostFunction: invokeOp.HostFunction,
		Auth:         authEntries,
	}

	txFinal, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false, // Seq sudah naik di txSim
			Operations:           []txnbuild.Operation{invokeOpFinal},
			BaseFee:              totalFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)

	txFinalBase64, _ := txFinal.Base64()

	// 7. HACKING XDR MANUAL (Injeksi Footprint)
	var env xdr.TransactionEnvelope
	xdr.SafeUnmarshalBase64(txFinalBase64, &env)
	if env.Type == xdr.EnvelopeTypeEnvelopeTypeTx {
		env.V1.Tx.Ext.V = 1
		env.V1.Tx.Ext.SorobanData = &txData
	}

	injectedB64, _ := xdr.MarshalBase64(env)
	genericTx, _ := txnbuild.TransactionFromXDR(injectedB64)
	txFinalObj, _ := genericTx.Transaction()

	// 8. Sign dan Kirim
	networkPassphrase := network.PublicNetworkPassphrase
	if s.network == "testnet" {
		networkPassphrase = network.TestNetworkPassphrase
	}

	txSigned, err := txFinalObj.Sign(networkPassphrase, kp)
	signedB64, _ := txSigned.Base64()

	var sendResp SendResponse
	if err := s.callRPC("sendTransaction", signedB64, &sendResp); err != nil {
		return "", authorAddress, fmt.Errorf("send transaction failed: %w", err)
	}
	if sendResp.Error != nil {
		return "", authorAddress, fmt.Errorf("RPC send error: %s", sendResp.Error.Message)
	}

	return sendResp.Result.Hash, authorAddress, nil
}

func (s *stellarService) GetAccountBalance() (string, error) {
	if s.secret == "" {
		return "0", fmt.Errorf("STELLAR_SECRET_KEY not configured")
	}

	kp, err := keypair.ParseFull(s.secret)
	if err != nil {
		return "0", err
	}

	accountRequest := horizonclient.AccountRequest{AccountID: kp.Address()}
	account, err := s.horizonClient.AccountDetail(accountRequest)
	if err != nil {
		return "0", err
	}

	for _, balance := range account.Balances {
		if balance.Asset.Type == "native" {
			return balance.Balance, nil
		}
	}

	return "0", nil
}

func (s *stellarService) callRPC(method string, txB64 string, target interface{}) error {
	reqBody := RPCRequest{
		JSONRPC: "2.0",
		ID:      int(time.Now().Unix()),
		Method:  method,
	}
	reqBody.Params.Transaction = txB64

	jsonBytes, _ := json.Marshal(reqBody)
	resp, err := http.Post(s.rpcURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return json.Unmarshal(bodyBytes, target)
}
