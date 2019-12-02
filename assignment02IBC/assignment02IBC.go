package assignment02IBC

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	a1 "github.com/syedaraiz/projectIBC/assignment01IBC"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func isPresent(selectedNodes []string, curr string) bool {
	for _, s := range selectedNodes {
		if curr == s {
			return true
		}
	}
	return false
}
func getSelectedNodes(ports []string, n int) []string {

	nPorts := len(ports)

	if n+1 == nPorts {
		return ports[:nPorts-1]
	} else if n == 0 {
		return ports[1:]
	}

	var selectedNodes []string

	for i := 0; i < n; {
		idx := randInt(0, nPorts)
		if !isPresent(selectedNodes, ports[idx]) && idx != n {
			selectedNodes = append(selectedNodes, ports[idx])
			i++
		}
	}

	return selectedNodes
}

func sendSelectedNodes(selectedNodes []string, encoder *gob.Encoder) {

	err := encoder.Encode(selectedNodes)

	if err != nil {
		log.Println(err)
	}
}

func getPort(decoder *gob.Decoder) string {

	var recievedPort string

	err := decoder.Decode(&recievedPort)

	if err != nil {
		log.Println(err)
	}
	return recievedPort
}

func getMessageType(decoder *gob.Decoder) string {
	return getPort(decoder)
}

func sendBlockChain(chainHead *a1.Block, encoder *gob.Encoder) {

	err := encoder.Encode(chainHead)

	if err != nil {
		log.Println(err)
	}
}

func coinBaseTransaction(port string, chainHead *a1.Block) *a1.Block {
	tx := a1.Transaction{
		Sender:   "coinbase",
		Receiver: port,
		Amount:   100,
	}

	var txs []a1.Transaction

	txs = append(txs, tx)

	var block *a1.Block
	if chainHead == nil {
		block = &a1.Block{
			Transactions: txs,
			PrevHash:     "0",
			PrevPointer:  chainHead,
		}
	} else {
		block = &a1.Block{
			Transactions: txs,
			PrevHash:     a1.CalculateHash(chainHead),
			PrevPointer:  chainHead,
		}
	}

	return block
}

func newConnenction(portsChan chan []string, chainChan chan *a1.Block, ports []string, chainHead *a1.Block, decoder *gob.Decoder) {

	ports = append(ports, getPort(decoder))

	log.Println("A node has joined:", ":"+ports[len(ports)-1])

	connection, err := net.Dial("tcp", ":"+ports[len(ports)-1])

	if err != nil {
		log.Println(err)
		return
	}

	selectedNodes := getSelectedNodes(ports, len(ports)-1)

	encoder := gob.NewEncoder(connection)

	sendSelectedNodes(selectedNodes, encoder)

	block := coinBaseTransaction(ports[0], chainHead)

	chainHead = a1.InsertBlock(block, chainHead)

	sendBlockChain(chainHead, encoder)

	log.Println("Sending blockchain")

	log.Println("Updating other peers")

	boardcastBlock(ports[:len(ports)-1], block)

	connection.Close()

	portsChan <- ports

	chainChan <- chainHead

}

func selectMiner(ports []string, experiencePoints []float32, stakeCoins []float32, isAllowed []bool, decoder *gob.Decoder) {

	var tx a1.Transaction

	decoder.Decode(&tx)

	log.Println("Transaction recieved:", tx)
	//Selection Of Miner Before
	//	idx := randInt(0, len(ports)-1) //Selected Miner
	idx := 0
	//Selection Of Miner After
	minerCriteria := ((experiencePoints[0] * 60) / 100) + ((stakeCoins[0] * 40) / 100)
	for i := 0; i < len(ports); i++ {
		nodeCriterialValue := ((experiencePoints[i] * 60) / 100) + ((stakeCoins[i] * 40) / 100)
		if nodeCriterialValue > minerCriteria && isAllowed[i] {
			minerCriteria = nodeCriterialValue
			idx = i
		}
	}
	isAllowed[idx] = false
	log.Println("Selected miner:", ":"+ports[idx])
	experiencePoints[idx]++
	time.AfterFunc(1*time.Hour, func() { isAllowed[idx] = true })

	connection, err := net.Dial("tcp", ":"+ports[idx])

	if err != nil {
		log.Println(err)
		return
	}

	encoder := gob.NewEncoder(connection)

	encoder.Encode("verify")
	encoder.Encode(tx)

}
func Satoshi(satoshiPortNumber string, myPortNumber string, noOfNode int) {

	var err error
	var ln net.Listener
	var connection net.Conn
	var chainHead *a1.Block
	var ports []string
	//Changes
	var experiencePoints []float32
	var stakeCoins []float32
	var isAllowed []bool

	var messageType string

	fmt.Println("satochi is listening at:", ":"+satoshiPortNumber)

	ln, err = net.Listen("tcp", ":"+satoshiPortNumber)

	for i := 0; i < noOfNode+1; {

		connection, err = ln.Accept()

		decoder := gob.NewDecoder(connection)

		messageType = getMessageType(decoder)

		if messageType == "connect" {

			ports = append(ports, getPort(decoder))
			experiencePoints = append(experiencePoints, 0)
			stakeCoins = append(stakeCoins, 0)
			isAllowed = append(isAllowed, true)

			log.Println("A node has joined:", ":"+ports[len(ports)-1])

			connection.Close()

			if i == 0 {

				chainHead = a1.InsertBlock(coinBaseTransaction(ports[0], chainHead), nil)

				i++

				continue
			}
			chainHead = a1.InsertBlock(coinBaseTransaction(ports[0], chainHead), chainHead)

			i++
		} else if messageType == "transaction" {

			log.Println("Cannot make transaction till required nodes are connected")

		} else {
			log.Println("Invalid message type")
		}

	}

	//sending block chain to all connected pears

	for i, port := range ports {

		connection, err = net.Dial("tcp", ":"+port)

		if err != nil {
			log.Println(err)
			continue
		}

		selectedNodes := getSelectedNodes(ports, i)

		encoder := gob.NewEncoder(connection)

		sendSelectedNodes(selectedNodes, encoder)

		sendBlockChain(chainHead, encoder)

		connection.Close()
	}

	chainChan := make(chan *a1.Block)
	portsChan := make(chan []string)
	for {

		connection, err = ln.Accept()

		decoder := gob.NewDecoder(connection)

		messageType = getMessageType(decoder)

		if messageType == "connect" {
			go newConnenction(portsChan, chainChan, ports, chainHead, decoder)
		} else if messageType == "transaction" {
			go selectMiner(ports, experiencePoints, stakeCoins, isAllowed, decoder)
		} else {
			log.Println("Invalid message type")
		}
		if messageType == "connect" {
			ports = <-portsChan
			chainHead = <-chainChan
		}
	}
}

func connectToSatoshi(myPortNumber string, satoshiPortNumber string) ([]string, *a1.Block, net.Listener) {

	connection, err := net.Dial("tcp", ":"+satoshiPortNumber)

	for err != nil {
		connection, err = net.Dial("tcp", ":"+satoshiPortNumber)
		if err != nil {
			log.Println(err)
		}
	}

	encoder := gob.NewEncoder(connection)

	encoder.Encode("connect")

	encoder.Encode(myPortNumber)

	connection.Close()

	var ports []string

	ln, err := net.Listen("tcp", ":"+myPortNumber)

	connection, err = ln.Accept()

	decoder := gob.NewDecoder(connection)
	err = decoder.Decode(&ports)

	if err != nil {
		log.Println(err)
	}

	var chainHead *a1.Block

	err = decoder.Decode(&chainHead)

	return ports, chainHead, ln
}

func CreateTransaction(myPortNumber string, satoshiPortNumber string) {

	for {
		var choice int
		_, err := fmt.Scanf("%d", &choice)

		fmt.Print("Receiver: ")
		var receiver string
		_, err = fmt.Scanf("%s", &receiver)

		fmt.Print("Amount: ")
		var amount float64
		_, err = fmt.Scanf("%f", &amount)

		fmt.Println(myPortNumber, receiver, amount)

		connection, err := net.Dial("tcp", ":"+satoshiPortNumber)

		if err != nil {
			log.Println(err)
		}

		encoder := gob.NewEncoder(connection)
		encoder.Encode("transaction")

		tx := a1.Transaction{
			Sender:   myPortNumber,
			Receiver: receiver,
			Amount:   amount,
		}
		encoder.Encode(tx)
	}
}

func boardcastBlock(ports []string, block *a1.Block) {
	for _, port := range ports {

		connection, err := net.Dial("tcp", ":"+port)
		if err != nil {
			log.Println(err)
			continue
		}

		encoder := gob.NewEncoder(connection)

		encoder.Encode("block")
		encoder.Encode(block)

		connection.Close()
	}
}
func isBoardcast(chainHead *a1.Block, block *a1.Block) bool {
	return a1.CalculateHash(chainHead) == a1.CalculateHash(block)
}
func isValidTransaction(newTx a1.Transaction, chainHead *a1.Block) bool {
	curr := chainHead
	var sum float64
	sum = 0
	for curr != nil {
		for _, tx := range curr.Transactions {
			if newTx.Sender == tx.Receiver {
				sum += tx.Amount
			} else if newTx.Sender == tx.Sender {
				sum -= tx.Amount
			}
		}
		curr = curr.PrevPointer
	}
	if newTx.Amount <= sum {
		return true
	}
	return false
}
func Normal(myPortNumber string, satoshiPortNumber string) {
	ports, chainHead, ln := connectToSatoshi(myPortNumber, satoshiPortNumber)

	log.Println("Your connected nodes are:", ports)
	a1.ListBlocks(chainHead)
	log.Println("Blockchain revcieved as above")

	for {
		connection, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		decoder := gob.NewDecoder(connection)

		messageType := getMessageType(decoder)

		log.Println("Message type:", messageType)

		if messageType == "block" {
			var block *a1.Block
			err = decoder.Decode(&block)
			if !isBoardcast(chainHead, block) {
				chainHead = a1.InsertBlock(block, chainHead)
				log.Println("broadcasting")
				boardcastBlock(ports, block)
			} else {
				log.Println("already Broadcasted")
			}

		} else if messageType == "verify" {
			var tx a1.Transaction
			err = decoder.Decode(&tx)

			block := coinBaseTransaction(myPortNumber, chainHead)

			if isValidTransaction(tx, chainHead) {
				block.Transactions = append(block.Transactions, tx)
				block.PrevHash = a1.CalculateHash(chainHead)
			}
			boardcastBlock(ports, block)
		}
		connection.Close()
	}
}
