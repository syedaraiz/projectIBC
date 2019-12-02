package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	a2 "github.com/syedaraiz/projectIBC/assignment02IBC"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	satoshiBool := flag.Bool("s", false, "a bool")
	normalBool := flag.Bool("o", true, "a bool")

	flag.Parse()

	if *satoshiBool && *normalBool {
		fmt.Println("cannot run as satoshi and normal at sametime")
		os.Exit(0)
	} else if !*satoshiBool && !*normalBool {
		fmt.Println("specify type")
		os.Exit(0)
	}

	myPortNumber := flag.Args()[0]
	satoshiPortNumber := flag.Args()[1]

	if *satoshiBool {
		noOfNode, err := strconv.Atoi(flag.Args()[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("1. New transaction")
		go a2.Satoshi(satoshiPortNumber, myPortNumber, noOfNode)
		go a2.CreateTransaction(myPortNumber, satoshiPortNumber)
		go a2.Normal(myPortNumber, satoshiPortNumber)
	} else if *normalBool {
		fmt.Println("1. New transaction")
		go a2.CreateTransaction(myPortNumber, satoshiPortNumber)
		go a2.Normal(myPortNumber, satoshiPortNumber)
	}
	var channel chan int
	<-channel
}
