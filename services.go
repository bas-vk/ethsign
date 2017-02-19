package main

import (
	"bufio"
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type PersService struct {
	ks *keystore.KeyStore
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PersService) NewAccount(password string) (common.Address, error) {
	acc, err := s.ks.NewAccount(password)
	if err == nil {
		return acc.Address, nil
	}
	return common.Address{}, err
}

// UnlockAccount will unlock the account associated with the given address with the
// given password for duration seconds. If duration is nil it will unlock the account
// until its explicit locked again. It returns an indication if the account was unlocked.
func (ps *PersService) UnlockAccount(addr common.Address, passwd string, duration *uint64) (bool, error) {
	acc := accounts.Account{Address: addr}
	var err error
	if duration == nil {
		err = ps.ks.Unlock(acc, passwd)
	} else {
		err = ps.ks.TimedUnlock(acc, passwd, time.Duration((*duration))*time.Second)
	}
	return err == nil, err
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (ps *PersService) LockAccount(addr common.Address) bool {
	return ps.ks.Lock(addr) == nil
}

type EthService struct {
	ks *keystore.KeyStore
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Big     `json:"gas"`
	GasPrice hexutil.Big     `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
	Nonce    hexutil.Uint64  `json:"nonce"`
	ChainId  uint64          `json:"chainId"`
}

func (args *SendTxArgs) toTransaction() *types.Transaction {
	if args.To == nil {
		return types.NewContractCreation(uint64(args.Nonce), (*big.Int)(args.Value), args.Gas.ToInt(), args.GasPrice.ToInt(), args.Data)
	}
	return types.NewTransaction(uint64(args.Nonce), *args.To, args.Value.ToInt(), args.Gas.ToInt(), args.GasPrice.ToInt(), args.Data)
}

func (args *SendTxArgs) chainID() *big.Int {
	if args.ChainId == 0 {
		return big.NewInt(1)
	}
	return new(big.Int).SetUint64(args.ChainId)
}

func (es *EthService) Accounts() []common.Address {
	var addresses []common.Address
	for _, acc := range es.ks.Accounts() {
		addresses = append(addresses, acc.Address)
	}
	return addresses
}

// SendTransaction is here just for convenience since web3 calls this function to submit
// a transaction we can still use web3 to send transaction. In this case it will return
// the RLP encoding of the signed transaction instead of the transaction hash.
func (es *EthService) SendTransaction(args SendTxArgs) (hexutil.Bytes, error) {
	return es.SignTransaction(args)
}

// SignTransaction creates a signe transaction and returns the RLP encoded transaction.
func (es *EthService) SignTransaction(args SendTxArgs) (hexutil.Bytes, error) {
	signedTx, err := es.ks.SignTx(accounts.Account{Address: args.From}, args.toTransaction(), args.chainID())
	if err != nil {
		return hexutil.Bytes{}, err
	}

	var data bytes.Buffer
	buf := bufio.NewWriter(&data)
	if err = signedTx.EncodeRLP(buf); err != nil {
		return hexutil.Bytes{}, err
	}

	buf.Flush()
	return data.Bytes(), nil
}
