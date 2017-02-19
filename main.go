package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	keystoreLocation = flag.String("keystore", "", "keystore path")
)

func main() {
	flag.Parse()

	ks := keystore.NewKeyStore(*keystoreLocation, keystore.StandardScryptN, keystore.StandardScryptP)

	handler := rpc.NewServer()
	pers := &PersService{ks}
	handler.RegisterName("personal", pers)
	eth := &EthService{ks}
	handler.RegisterName("eth", eth)
	consoleBackend := rpc.DialInProc(handler)

	datadir, err := ioutil.TempDir("", "ethsign")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(datadir)

	cfg := console.Config{
		DataDir: datadir,
		Client:  consoleBackend,
	}

	console, err := console.New(cfg)
	if err != nil {
		panic(err)
	}

	console.Interactive()
}
