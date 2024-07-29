package blockchain

// Block defines the structure for the blockchain node
type Block struct {
	Pos       int
	Data      Entreprenuer
	Timestamp string
	Hash      string
	PrevHash  string
}

// Blockchain defines the structure for the blockchain
type Blockchain struct {
	blocks []*Block
}

// Entreprenuer define the structure for the entreprenuer
type Entreprenuer struct {
	FirstName  string   `json:"first_name"`
	SecondName string   `json:"second_name"`
	Location   string   `json:"location"`
	Business   Business `json:"business"`
	Phone      string   `json:"phone"`
	NationalID string   `json:"national_id"`
}

// Business define the structure for the business
type Business struct {
	BusinessID    string `json:"business_id"`
	Status        string `json:"status"`
	BusinessValue string `json:"business_value"`
	Name          string `json:"name"`
	Address       string `json:"address"`
}

// BlockchainInstance declares a global blockchain instance
var BlockChain *Blockchain
