package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jeffprestes/copper/contracts"
)

// LogFill represents a LogFill event raised by the Exchange contract.
type LogFill struct {
	Maker                  common.Address
	Taker                  common.Address
	FeeRecipient           common.Address
	MakerToken             common.Address
	TakerToken             common.Address
	FilledMakerTokenAmount *big.Int
	FilledTakerTokenAmount *big.Int
	PaidMakerFee           *big.Int
	PaidTakerFee           *big.Int
	Tokens                 [32]byte
	OrderHash              [32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

func main() {
	//Connect to mainnet
	client, err := ethclient.Dial("https://mainnet.infura.io/QPF0qjGpH9OjFuuMrCse")
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

	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.ExchangeABI)))
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: keccak256("LogFill(address,address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes32)")
	logFillEvent := common.HexToHash("0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3")

	for _, vLog := range logs {
		fmt.Printf("Log Index: %d\n", vLog.Index)
		if vLog.Topics[0].Hex() == logFillEvent.Hex() {
			fillEvent := LogFill{}
			err := contractAbi.Unpack(&fillEvent, "LogFill", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			fillEvent.Maker = common.HexToAddress(vLog.Topics[1].Hex())
			fillEvent.FeeRecipient = common.HexToAddress(vLog.Topics[2].Hex())
			fillEvent.Tokens = vLog.Topics[3]

			fmt.Printf("Maker: %s\n", fillEvent.Maker.Hex())
			fmt.Printf("Taker: %s\n", fillEvent.Taker.Hex())
			fmt.Printf("Maker Token: %s\n", fillEvent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", fillEvent.TakerToken.Hex())
			fmt.Printf("Fee Recipient: %s\n", fillEvent.FeeRecipient.Hex())
			fmt.Printf("Filled Maker Token Amount: %s\n", fillEvent.FilledMakerTokenAmount.String())
			fmt.Printf("Filled Taker Token Amount: %s\n", fillEvent.FilledTakerTokenAmount.String())
			fmt.Printf("Paid Maker Fee: %s\n", fillEvent.PaidMakerFee.String())
			fmt.Printf("Paid Taker Fee: %s\n", fillEvent.PaidTakerFee.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(fillEvent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(fillEvent.OrderHash[:]))
		}
		fmt.Printf("\n\n")
	}
}
