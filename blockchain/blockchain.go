package blockchain

// Define the structure for the entreprenuer
type Entreprenuer struct {
	FirstName  string   `json:"first_name"`
	SecondName string   `json:"second_name"`
	Location   string   `json:"location"`
	Business   Business `json:"business"`
	Phone      string   `json:"phone"`
	NationalID string   `json:"national_id"`
}
