package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
)

func AddCategory(id, name, unit string, pricePerUnit int) error {
	if id == "" {
		return errors.New("品类ID不能为空")
	}
	if name == "" {
		return errors.New("品类名称不能为空")
	}
	if unit == "" {
		return errors.New("单位不能为空")
	}
	if pricePerUnit <= 0 {
		return errors.New("每单位回收价必须是正整数")
	}

	categories, err := Store.LoadCategories()
	if err != nil {
		return fmt.Errorf("加载品类数据失败: %v", err)
	}

	for _, cat := range categories {
		if cat.ID == id {
			return fmt.Errorf("品类ID %s 已存在", id)
		}
	}

	categories = append(categories, Category{
		ID:           id,
		Name:         name,
		Unit:         unit,
		PricePerUnit: pricePerUnit,
	})

	if err := Store.SaveCategories(categories); err != nil {
		return fmt.Errorf("保存品类数据失败: %v", err)
	}

	fmt.Printf("品类添加成功: %s %s (单价: %d分/%s)\n", id, name, pricePerUnit, unit)
	return nil
}

func PurchaseItem(categoryID, supplier, phone string, weight float64, date string) error {
	if categoryID == "" {
		return errors.New("品类ID不能为空")
	}
	if supplier == "" {
		return errors.New("供应商姓名不能为空")
	}
	if phone == "" {
		return errors.New("供应商电话不能为空")
	}
	const epsilon = 1e-9
	if weight <= epsilon {
		return errors.New("重量必须是正数")
	}
	if date == "" {
		return errors.New("日期不能为空")
	}

	categories, err := Store.LoadCategories()
	if err != nil {
		return fmt.Errorf("加载品类数据失败: %v", err)
	}

	var category Category
	found := false
	for _, cat := range categories {
		if cat.ID == categoryID {
			category = cat
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("品类ID %s 不存在", categoryID)
	}

	amount := int(math.Floor(weight * float64(category.PricePerUnit)))
	if amount <= 0 {
		return errors.New("计算金额必须为正数，请检查重量和单价")
	}

	purchases, err := Store.LoadPurchases()
	if err != nil {
		return fmt.Errorf("加载收购记录失败: %v", err)
	}

	purchases = append(purchases, Purchase{
		CategoryID: categoryID,
		Supplier:   supplier,
		Phone:      phone,
		Weight:     weight,
		Amount:     amount,
		Date:       date,
	})

	if err := Store.SavePurchases(purchases); err != nil {
		return fmt.Errorf("保存收购记录失败: %v", err)
	}

	fmt.Printf("收购记录添加成功:\n")
	fmt.Printf("  品类: %s (%s)\n", category.Name, category.Unit)
	fmt.Printf("  供应商: %s (%s)\n", supplier, phone)
	fmt.Printf("  重量: %v %s\n", weight, category.Unit)
	fmt.Printf("  单价: %d分/%s\n", category.PricePerUnit, category.Unit)
	fmt.Printf("  金额: %d分 (%.2f元)\n", amount, float64(amount)/100.0)
	fmt.Printf("  日期: %s\n", date)
	return nil
}

func AddPayout(phone string, amount int, date string) error {
	if phone == "" {
		return errors.New("供应商电话不能为空")
	}
	if amount <= 0 {
		return errors.New("付款金额必须是正整数")
	}
	if date == "" {
		return errors.New("日期不能为空")
	}

	result, err := Store.AddPayoutAtomic(phone, amount, date)
	if err != nil {
		return err
	}

	fmt.Printf("付款记录添加成功:\n")
	fmt.Printf("  电话: %s\n", phone)
	fmt.Printf("  本次付款: %d分 (%.2f元)\n", amount, float64(amount)/100.0)
	fmt.Printf("  累计应得: %d分 (%.2f元)\n", result.TotalEarned, float64(result.TotalEarned)/100.0)
	fmt.Printf("  累计已付: %d分 (%.2f元)\n", result.NewPaid, float64(result.NewPaid)/100.0)
	fmt.Printf("  剩余未付: %d分 (%.2f元)\n", result.TotalEarned-result.NewPaid, float64(result.TotalEarned-result.NewPaid)/100.0)
	fmt.Printf("  日期: %s\n", date)
	return nil
}

func SupplierBalance(phone string) error {
	if phone == "" {
		return errors.New("供应商电话不能为空")
	}

	totalEarned, err := calculateSupplierEarned(phone)
	if err != nil {
		return err
	}

	totalPaid, err := calculateSupplierPaid(phone)
	if err != nil {
		return err
	}

	unpaid := totalEarned - totalPaid

	purchases, err := Store.LoadPurchases()
	if err != nil {
		return fmt.Errorf("加载收购记录失败: %v", err)
	}

	var supplierName string
	for _, p := range purchases {
		if p.Phone == phone {
			supplierName = p.Supplier
			break
		}
	}

	if supplierName == "" {
		supplierName = "未知"
	}

	fmt.Printf("供应商余额查询:\n")
	fmt.Printf("  姓名: %s\n", supplierName)
	fmt.Printf("  电话: %s\n", phone)
	fmt.Printf("  累计应得: %d分 (%.2f元)\n", totalEarned, float64(totalEarned)/100.0)
	fmt.Printf("  累计已付: %d分 (%.2f元)\n", totalPaid, float64(totalPaid)/100.0)
	fmt.Printf("  剩余未付: %d分 (%.2f元)\n", unpaid, float64(unpaid)/100.0)
	return nil
}

var monthPattern = regexp.MustCompile(`^\d{4}-\d{2}$`)

func MonthlySummary(month string) error {
	if month == "" {
		return errors.New("月份不能为空")
	}
	if !monthPattern.MatchString(month) {
		return fmt.Errorf("月份格式不正确，应为 YYYY-MM，例如: 2024-03")
	}

	purchases, err := Store.LoadPurchases()
	if err != nil {
		return fmt.Errorf("加载收购记录失败: %v", err)
	}

	categories, err := Store.LoadCategories()
	if err != nil {
		return fmt.Errorf("加载品类数据失败: %v", err)
	}

	catMap := make(map[string]Category)
	for _, cat := range categories {
		catMap[cat.ID] = cat
	}

	type SummaryItem struct {
		CategoryID    string
		CategoryName  string
		Unit          string
		TotalWeight   float64
		TotalAmount   int
	}

	summary := make(map[string]*SummaryItem)
	totalWeight := 0.0
	totalAmount := 0

	for _, p := range purchases {
		if strings.HasPrefix(p.Date, month) {
			item, exists := summary[p.CategoryID]
			if !exists {
				cat, ok := catMap[p.CategoryID]
				catName := p.CategoryID
				unit := ""
				if ok {
					catName = cat.Name
					unit = cat.Unit
				}
				item = &SummaryItem{
					CategoryID:   p.CategoryID,
					CategoryName: catName,
					Unit:         unit,
				}
				summary[p.CategoryID] = item
			}
			item.TotalWeight += p.Weight
			item.TotalAmount += p.Amount
			totalWeight += p.Weight
			totalAmount += p.Amount
		}
	}

	fmt.Printf("%s 月度汇总:\n", month)
	fmt.Printf("%-10s %-15s %-15s %-15s\n", "品类ID", "品类名称", "总重量", "总金额")
	fmt.Println(strings.Repeat("-", 55))

	if len(summary) == 0 {
		fmt.Println("  (本月无收购记录)")
	} else {
		for _, item := range summary {
			fmt.Printf("%-10s %-15s %-15s %-15s\n",
				item.CategoryID,
				item.CategoryName,
				fmt.Sprintf("%.2f %s", item.TotalWeight, item.Unit),
				fmt.Sprintf("%d分(%.2f元)", item.TotalAmount, float64(item.TotalAmount)/100.0))
		}
	}

	fmt.Println(strings.Repeat("-", 55))
	fmt.Printf("%-10s %-15s %-15s %-15s\n",
		"合计", "",
		fmt.Sprintf("%.2f", totalWeight),
		fmt.Sprintf("%d分(%.2f元)", totalAmount, float64(totalAmount)/100.0))
	return nil
}

func calculateSupplierEarned(phone string) (int, error) {
	purchases, err := Store.LoadPurchases()
	if err != nil {
		return 0, fmt.Errorf("加载收购记录失败: %v", err)
	}

	total := 0
	for _, p := range purchases {
		if p.Phone == phone {
			total += p.Amount
		}
	}
	return total, nil
}

func calculateSupplierPaid(phone string) (int, error) {
	payouts, err := Store.LoadPayouts()
	if err != nil {
		return 0, fmt.Errorf("加载付款记录失败: %v", err)
	}

	total := 0
	for _, p := range payouts {
		if p.Phone == phone {
			total += p.Amount
		}
	}
	return total, nil
}
