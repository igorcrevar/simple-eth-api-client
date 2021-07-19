package handlers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	Models "github.com/igorcrevar/simple-eth-api-client/models"
	Modules "github.com/igorcrevar/simple-eth-api-client/modules"
)

var client *ethclient.Client
var myPrivateKey string

type tuple struct {
	name string
	fn   func(w http.ResponseWriter, r *http.Request) interface{}
}

func Register(_client *ethclient.Client, _privateKey string) {
	client = _client
	myPrivateKey = _privateKey

	register("/api/v1/eth", []tuple{
		{"chain", getChain},
		{"latest-block", getLatestBlock},
		{"get-tx", getTransaction},
		{"send-eth", sendEth},
		{"get-balance", getBalance},
	})
}

func register(prefix string, pairs []tuple) {
	httpProxyHandler := func(fn func(w http.ResponseWriter, r *http.Request) interface{}) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			result := fn(w, r)
			json.NewEncoder(w).Encode(result)
		}
	}

	if prefix == "" || prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}
	for _, pair := range pairs {
		http.HandleFunc(prefix+pair.name, httpProxyHandler(pair.fn))
	}
}

func getChain(w http.ResponseWriter, r *http.Request) interface{} {
	offsetStr := r.URL.Query().Get("offset")
	offset := new(big.Int)
	offset, _ = offset.SetString(offsetStr, 10)
	_blocks := Modules.GetBlockChain(*client, offset)
	return _blocks
}

func getLatestBlock(w http.ResponseWriter, r *http.Request) interface{} {
	_block := Modules.GetLatestBlock(*client)
	return _block
}

func getTransaction(w http.ResponseWriter, r *http.Request) interface{} {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		return &Models.Error{
			Code:    400,
			Message: "Malformed request",
		}
	}
	txHash := common.HexToHash(hash)
	_tx := Modules.GetTxByHash(*client, txHash)

	if _tx != nil {
		return _tx
	} else {
		return &Models.Error{
			Code:    404,
			Message: "Tx Not Found!",
		}
	}
}

func sendEth(w http.ResponseWriter, r *http.Request) interface{} {
	decoder := json.NewDecoder(r.Body)
	var t Models.TransferEthRequest

	err := decoder.Decode(&t)

	if err != nil {
		fmt.Println(err)
		return &Models.Error{
			Code:    400,
			Message: "Malformed request",
		}
	}

	if t.PrivateKey == "" {
		t.PrivateKey = myPrivateKey
	}

	_hash, err := Modules.TransferEth(*client, t.PrivateKey, t.To, t.Amount)

	if err != nil {
		fmt.Println(err)
		return &Models.Error{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return &Models.HashResponse{
		Hash: _hash,
	}
}

func getBalance(w http.ResponseWriter, r *http.Request) interface{} {
	address := r.URL.Query().Get("address")
	if address == "" {
		return &Models.Error{
			Code:    400,
			Message: "Malformed request",
		}
	}

	balance, err := Modules.GetAddressBalance(*client, address)

	if err != nil {
		fmt.Println(err)
		return &Models.Error{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return &Models.BalanceResponse{
		Address: address,
		Balance: balance,
		Symbol:  "Ether",
		Units:   "Wei",
	}
}
