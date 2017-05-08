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
	lightKDF         = flag.Bool("lightkdf", false, "reduce key-derivation RAM & CPU usage at some expense of KDF strength")
)

func main() {
	flag.Parse()

	n, p := keystore.StandardScryptN, keystore.StandardScryptP
	if *lightKDF {
		n, p = keystore.LightScryptN, keystore.LightScryptP
	}

	ks := keystore.NewKeyStore(*keystoreLocation, n, p)

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
