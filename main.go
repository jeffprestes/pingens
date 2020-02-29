package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jeffprestes/goethereumhelper"
	"github.com/wealdtech/ethereal/cli"
	"github.com/wealdtech/go-ens"
)

var quiet bool

func main() {
	client, err := goethereumhelper.GetCustomNetworkClient("https://mainnet.infura.io/v3/1eae79e9335242369fe9b17b9413d721")
	if err != nil {
		fmt.Println("It was not possible to connect to an Ethereum node ", err)
		return
	}
	if !genericInfo(client, "code.jeffprestes.eth") {
		fmt.Println("Domain is not found")
	}
}

// genericInfo prints generic info about any ENS domain.
// It returns true if the domain exists, otherwise false
func genericInfo(client *ethclient.Client, name string) bool {
	registry, err := ens.NewRegistry(client)
	cli.ErrCheck(err, quiet, "Failed to obtain registry contract")
	controllerAddress, err := registry.Owner(name)
	cli.ErrCheck(err, quiet, "Failed to obtain controller")
	if controllerAddress == ens.UnknownAddress {
		fmt.Println("Owner not set")
		return false
	}
	controllerName, _ := ens.ReverseResolve(client, controllerAddress)
	if controllerName == "" {
		fmt.Printf("Controller is %s\n", controllerAddress.Hex())
	} else {
		fmt.Printf("Controller is %s (%s)\n", controllerName, controllerAddress.Hex())
	}

	// Resolver
	resolverAddress, err := registry.ResolverAddress(name)
	if err != nil || resolverAddress == ens.UnknownAddress {
		fmt.Println("Resolver not configured")
		return true
	}
	resolverName, _ := ens.ReverseResolve(client, resolverAddress)
	if resolverName == "" {
		fmt.Printf("Resolver is %s\n", resolverAddress.Hex())
	} else {
		fmt.Printf("Resolver is %s (%s)\n", resolverName, resolverAddress.Hex())
	}

	// Address
	address, err := ens.Resolve(client, name)
	if err == nil && address != ens.UnknownAddress {
		fmt.Printf("Domain resolves to %s\n", address.Hex())
		// Reverse resolution
		reverseDomain, err := ens.ReverseResolve(client, address)
		if err == nil && reverseDomain != "" {
			fmt.Printf("Address resolves to %s\n", reverseDomain)
		}
	}

	// Content hash
	resolver, err := ens.NewResolverAt(client, name, resolverAddress)
	if err == nil {
		bytes, err := resolver.Contenthash()
		if err == nil && len(bytes) > 0 {
			contentHash, err := ens.ContenthashToString(bytes)
			if err == nil {
				fmt.Printf("Content hash is %v\n", contentHash)
			}
		}
	}

	return true
}
