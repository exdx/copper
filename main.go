package main

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	exchange "github.com/copper/contracts"
)

func main() {
	//Connect to mainnet
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	//Input address of 0x v1 contract and build event query
	contractAddress := common.HexToAddress("0x12459C951127e0c374FF9105DdA097662A027093")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6338488),
		ToBlock:   big.NewInt(6338504),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	//Run query against ethereum client
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(exchange.ExchangeABI)))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		fillEvent := struct {
			indexedmaker           string //address
			taker                  string //address
			feeRecipient           string //address
			makerToken             string //address
			takerToken             string //address
			filledMakerTokenAmount uint
			filledTakerTokenAmount uint
			paidMakerFee           uint
			paidTakerFee           uint
			indexedtokens          [32]byte
			orderHash              [32]byte
		}{}
		err := contractAbi.Unpack(&fillEvent, "LogFill", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(string(fillEvent.indexedmaker))
		// fmt.Println(string(fillEvent.taker))
	}

}
