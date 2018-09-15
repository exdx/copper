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
			Indexedmaker           string //address
			Taker                  string //address
			FeeRecipient           string //address
			MakerToken             string //address
			TakerToken             string //address
			FilledMakerTokenAmount uint
			FilledTakerTokenAmount uint
			PaidMakerFee           uint
			PaidTakerFee           uint
			Indexedtokens          [32]byte
			OrderHash              [32]byte
		}{}
		err := contractAbi.Unpack(&fillEvent, "LogFill", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(string(fillEvent.indexedmaker))
		// fmt.Println(string(fillEvent.taker))
	}

}
