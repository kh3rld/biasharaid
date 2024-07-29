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

// Define the structure for the business
type Business struct {
	BusinessID    string `json:"business_id"`
	Status        string `json:"status"`
	BusinessValue string `json:"business_value"`
	Name          string `json:"name"`
	Address       string `json:"address"`
}
