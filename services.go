package main

import (
	"bufio"
	"bytes"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type PersService struct {
	am *accounts.Manager
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PersService) NewAccount(password string) (common.Address, error) {
	acc, err := s.am.NewAccount(password)
	if err == nil {
		return acc.Address, nil
	}
	return common.Address{}, err
}

// UnlockAccount will unlock the account associated with the given address with the
// given password for duration seconds. If duration is nil it will unlock the account
// until its explicit locked again. It returns an indication if the account was unlocked.
func (ps *PersService) UnlockAccount(addr common.Address, passwd string, duration *rpc.HexNumber) (bool, error) {
	acc := accounts.Account{Address: addr}
	var err error
	if duration == nil {
		err = ps.am.Unlock(acc, passwd)
	} else {
		err = ps.am.TimedUnlock(acc, passwd, time.Duration(duration.Int())*time.Second)
	}
	return err == nil, err
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (ps *PersService) LockAccount(addr common.Address) bool {
	return ps.am.Lock(addr) == nil
}

type EthService struct {
	am *accounts.Manager
}

type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      rpc.HexNumber   `json:"gas"`
	GasPrice rpc.HexNumber   `json:"gasPrice"`
	Value    rpc.HexNumber   `json:"value"`
	Data     string          `json:"data"`
	Nonce    rpc.HexNumber   `json:"nonce"`
}

func (es *EthService) Accounts() []common.Address {
	var addresses []common.Address
	for _, acc := range es.am.Accounts() {
		addresses = append(addresses, acc.Address)
	}
	return addresses
}

// SendTransaction is here just for convenience since web3 calls this function to submit
// a transaction we can still use web3 to send transaction. In this case it will return
// the RLP encoding of the signed transaction instead of the transaction hash.
func (es *EthService) SendTransaction(args SendTxArgs) (rpc.HexBytes, error) {
	return es.SignTransaction(args)
}

// SignTransaction creates a signe transaction and returns the RLP encoded transaction.
func (es *EthService) SignTransaction(args SendTxArgs) (rpc.HexBytes, error) {
	var tx *types.Transaction
	if args.To == nil {
		tx = types.NewContractCreation(args.Nonce.Uint64(), args.Value.BigInt(), args.Gas.BigInt(), args.GasPrice.BigInt(), common.FromHex(args.Data))
	} else {
		tx = types.NewTransaction(args.Nonce.Uint64(), *args.To, args.Value.BigInt(), args.Gas.BigInt(), args.GasPrice.BigInt(), common.FromHex(args.Data))
	}

	signature, err := es.am.Sign(args.From, tx.SigHash().Bytes())
	if err != nil {
		return rpc.HexBytes{}, err
	}

	signedTx, err := tx.WithSignature(signature)
	if err != nil {
		return rpc.HexBytes{}, err
	}

	var data bytes.Buffer
	buf := bufio.NewWriter(&data)
	if err = signedTx.EncodeRLP(buf); err != nil {
		return rpc.HexBytes{}, err
	}

	buf.Flush()
	return data.Bytes(), nil
}
