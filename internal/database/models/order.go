package models

type Order struct {
    ID        int     `json:"id"`
    ProductID int     `json:"product_id"`
    Number    int     `json:"number"`
    Amount    float64 `json:"amount"`
    Time      string  `json:"time"`
    // Add other fields as needed
}
