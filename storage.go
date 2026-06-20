package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Category struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Unit          string `json:"unit"`
	PricePerUnit  int    `json:"price_per_unit"`
}

type Purchase struct {
	CategoryID string  `json:"category_id"`
	Supplier   string  `json:"supplier"`
	Phone      string  `json:"phone"`
	Weight     float64 `json:"weight"`
	Amount     int     `json:"amount"`
	Date       string  `json:"date"`
}

type Payout struct {
	Phone  string `json:"phone"`
	Amount int    `json:"amount"`
	Date   string `json:"date"`
}

type DataStore struct {
	CategoriesFile string
	PurchasesFile  string
	PayoutsFile    string
	mu             sync.Mutex
}

var Store = DataStore{
	CategoriesFile: "categories.json",
	PurchasesFile:  "purchases.json",
	PayoutsFile:    "payouts.json",
}

func loadJSON(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

func saveJSON(filename string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *DataStore) LoadCategories() ([]Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var categories []Category
	err := loadJSON(s.CategoriesFile, &categories)
	return categories, err
}

func (s *DataStore) SaveCategories(categories []Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return saveJSON(s.CategoriesFile, categories)
}

func (s *DataStore) LoadPurchases() ([]Purchase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var purchases []Purchase
	err := loadJSON(s.PurchasesFile, &purchases)
	return purchases, err
}

func (s *DataStore) SavePurchases(purchases []Purchase) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return saveJSON(s.PurchasesFile, purchases)
}

func (s *DataStore) LoadPayouts() ([]Payout, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var payouts []Payout
	err := loadJSON(s.PayoutsFile, &payouts)
	return payouts, err
}

func (s *DataStore) SavePayouts(payouts []Payout) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return saveJSON(s.PayoutsFile, payouts)
}
