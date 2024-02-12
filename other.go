package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Customer represents a customer in the database
type Customer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Order represents an order in the database
type Order struct {
	Item   string    `json:"item"`
	Amount float64   `json:"amount"`
	Time   time.Time `json:"time"`
}

var customers []Customer
var orders []Order

func main() {
	router := mux.NewRouter()

	// Customers endpoints
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", createCustomer).Methods("POST")

	// Orders endpoints
	router.HandleFunc("/orders", getOrders).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all customers
func getCustomers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(customers)
}

// Create a new customer
func createCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	json.NewDecoder(r.Body).Decode(&customer)
	customers = append(customers, customer)
	json.NewEncoder(w).Encode(customer)
}

// Get all orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(orders)
}

// Create a new order
func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)

	// Send SMS alert to the customer
	sendSMSAlert(order)
}

// Send SMS alert to the customer
func sendSMSAlert(order Order) {
	// Implement your SMS gateway integration here
	fmt.Printf("Sending SMS alert to customer: Order for %s created\n", order.Item)
}
