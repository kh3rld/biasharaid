package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// Entrepreneur defines the structure for the entrepreneur
type Entrepreneur struct {
	FirstName  string   `json:"first_name"`
	SecondName string   `json:"second_name"`
	Location   string   `json:"location"`
	Business   Business `json:"business"`
	Phone      string   `json:"phone"`
	NationalID string   `json:"national_id"`
	IsGenesis  bool     `json:"is_genesis"`
}

// Business defines the structure for the business
type Business struct {
	BusinessID    string `json:"business_id"`
	Status        string `json:"status"`
	BusinessValue string `json:"business_value"`
	Name          string `json:"name"`
	Address       string `json:"address"`
}

// Block defines the structure for the blockchain node
type Block struct {
	Pos       int
	Data      Entrepreneur
	Timestamp string
	Hash      string
	PrevHash  string
}

// CreateBlock creates a new block with the given data and previous hash
func (b *Block) CreateBlock(prevBlock *Block, person Entrepreneur) *Block {
	block := &Block{
		Pos:       prevBlock.Pos + 1,
		Data:      person,
		Timestamp: time.Now().String(),
		PrevHash:  prevBlock.Hash,
	}

	block.Hash = block.GenerateHash()

	return block
}

// GenerateHash generates a SHA-256 hash for the block
func (b *Block) GenerateHash() string {
	bytes, _ := json.Marshal(b.Data)
	data := string(b.Pos) + b.Timestamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// ValidateHash validates a given hash
func (b *Block) ValidateHash(hash string) bool {
	return b.GenerateHash() == hash
}

// Blockchain defines the structure for the blockchain
type Blockchain struct {
	blocks []*Block
}

// BlockchainInstance declares a global blockchain instance
var BlockchainInstance *Blockchain

// InitializeBlockchain initializes a new blockchain with the genesis block
func InitializeBlockchain() *Blockchain {
	genesisBlock := &Block{
		Pos:       0,
		Data:      Entrepreneur{IsGenesis: true},
		Timestamp: time.Now().String(),
		PrevHash:  "",
	}
	genesisBlock.Hash = genesisBlock.GenerateHash()

	return &Blockchain{[]*Block{genesisBlock}}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data Entrepreneur) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := prevBlock.CreateBlock(prevBlock, data)

	if validBlock(newBlock, prevBlock) {
		bc.blocks = append(bc.blocks, newBlock)
	}
}

// validBlock checks if a new block is valid based on the previous block's hash
func validBlock(newBlock, prevBlock *Block) bool {
	if prevBlock.Hash != newBlock.PrevHash {
		return false
	}

	if !newBlock.ValidateHash(newBlock.Hash) {
		return false
	}

	if prevBlock.Pos+1 != newBlock.Pos {
		return false
	}
	return true
}