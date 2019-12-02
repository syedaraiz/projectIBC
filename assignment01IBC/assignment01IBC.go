package assignment01IBC
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)
type Transaction struct {
		Sender string
		Receiver string
		Amount float64
	}
type Block struct {
	Transactions []Transaction
	PrevHash string
	PrevPointer *Block
}
func InsertBlock(block *Block, chainHead *Block) *Block {
	block.PrevPointer = chainHead
	chainHead = block
	return chainHead
}
func CalculateHash(block *Block) string {
	str := block.PrevHash
	for _, tx := range block.Transactions {
		str += tx.Sender + tx.Receiver + strconv.FormatFloat(tx.Amount, 'f', -1, 64)	}
		hash := sha256.Sum256([]byte(str))
		return hex.EncodeToString(hash[:])
	}
	func ListBlocks(chainHead *Block) {
		fmt.Println("--------------------------------------------------------------------------")
		curr := chainHead
		for curr != nil {
			for _, tx := range curr.Transactions {
				fmt.Println("Sender: " + tx.Sender)
				fmt.Println("Receiver: " + tx.Receiver)
				fmt.Println("Amount:", tx.Amount)		}
				fmt.Println("PrevHash: " + curr.PrevHash)
				fmt.Println("--------------------------------------------------------------------------")
				curr = curr.PrevPointer
				}
				fmt.Println()
		}
	func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {}
		// curr := chainHead
		// for curr != nil {
		// 	if curr.Transaction == oldTrans {
		// 		curr.Transaction = newTrans
		// 	}	// 	curr = curr.PrevPointer
		// }}func VerifyChain(chainHead *Block) {	var isValid bool = true	curr := chainHead	for curr.PrevPointer != nil {		if curr.PrevHash != CalculateHash(curr.PrevPointer) {			isValid = false		}		curr = curr.PrevPointer	} 	if isValid {		fmt.Println("Blockchain is valid")	} else {		fmt.Println("Blockchain is invalid")	}}
