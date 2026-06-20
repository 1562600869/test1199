package main

import (
	"encoding/json"
	"fmt"
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

type AddPayoutTxResult struct {
	TotalEarned int
	TotalPaid   int
	NewPaid     int
}

func (s *DataStore) AddPayoutAtomic(phone string, amount int, date string) (*AddPayoutTxResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var purchases []Purchase
	if err := loadJSON(s.PurchasesFile, &purchases); err != nil {
		return nil, fmt.Errorf("加载收购记录失败: %v", err)
	}

	totalEarned := 0
	for _, p := range purchases {
		if p.Phone == phone {
			totalEarned += p.Amount
		}
	}

	var payouts []Payout
	if err := loadJSON(s.PayoutsFile, &payouts); err != nil {
		return nil, fmt.Errorf("加载付款记录失败: %v", err)
	}

	totalPaid := 0
	for _, p := range payouts {
		if p.Phone == phone {
			totalPaid += p.Amount
		}
	}

	if totalPaid+amount > totalEarned {
		return nil, fmt.Errorf("累计付款(%d分)超过累计应得(%d分)，当前未付金额为%d分",
			totalPaid+amount, totalEarned, totalEarned-totalPaid)
	}

	payouts = append(payouts, Payout{
		Phone:  phone,
		Amount: amount,
		Date:   date,
	})

	if err := saveJSON(s.PayoutsFile, payouts); err != nil {
		return nil, fmt.Errorf("保存付款记录失败: %v", err)
	}

	return &AddPayoutTxResult{
		TotalEarned: totalEarned,
		TotalPaid:   totalPaid,
		NewPaid:     totalPaid + amount,
	}, nil
}
