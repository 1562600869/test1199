package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "add-category":
		err = cmdAddCategory(args)
	case "purchase":
		err = cmdPurchase(args)
	case "payout":
		err = cmdPayout(args)
	case "supplier-balance":
		err = cmdSupplierBalance(args)
	case "monthly":
		err = cmdMonthly(args)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Printf("未知命令: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func cmdAddCategory(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: add-category <ID> <名称> --unit <单位> --price-per-unit <单价(分)>")
	}

	id := args[0]
	name := args[1]
	args = args[2:]

	unit := ""
	pricePerUnitStr := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--unit":
			if i+1 >= len(args) {
				return fmt.Errorf("--unit 需要参数值")
			}
			unit = args[i+1]
			i++
		case "--price-per-unit":
			if i+1 >= len(args) {
				return fmt.Errorf("--price-per-unit 需要参数值")
			}
			pricePerUnitStr = args[i+1]
			i++
		}
	}

	if unit == "" {
		return fmt.Errorf("缺少 --unit 参数")
	}
	if pricePerUnitStr == "" {
		return fmt.Errorf("缺少 --price-per-unit 参数")
	}

	pricePerUnit, err := strconv.Atoi(pricePerUnitStr)
	if err != nil {
		return fmt.Errorf("--price-per-unit 必须是整数: %v", err)
	}

	return AddCategory(id, name, unit, pricePerUnit)
}

func cmdPurchase(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: purchase <品类ID> --supplier <姓名> --phone <电话> --weight <重量> --date <日期>")
	}

	categoryID := args[0]
	args = args[1:]

	supplier := ""
	phone := ""
	weightStr := ""
	date := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--supplier":
			if i+1 >= len(args) {
				return fmt.Errorf("--supplier 需要参数值")
			}
			supplier = args[i+1]
			i++
		case "--phone":
			if i+1 >= len(args) {
				return fmt.Errorf("--phone 需要参数值")
			}
			phone = args[i+1]
			i++
		case "--weight":
			if i+1 >= len(args) {
				return fmt.Errorf("--weight 需要参数值")
			}
			weightStr = args[i+1]
			i++
		case "--date":
			if i+1 >= len(args) {
				return fmt.Errorf("--date 需要参数值")
			}
			date = args[i+1]
			i++
		}
	}

	if supplier == "" {
		return fmt.Errorf("缺少 --supplier 参数")
	}
	if phone == "" {
		return fmt.Errorf("缺少 --phone 参数")
	}
	if weightStr == "" {
		return fmt.Errorf("缺少 --weight 参数")
	}
	if date == "" {
		return fmt.Errorf("缺少 --date 参数")
	}

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return fmt.Errorf("--weight 必须是数字: %v", err)
	}

	return PurchaseItem(categoryID, supplier, phone, weight, date)
}

func cmdPayout(args []string) error {
	supplierPhone := ""
	amountStr := ""
	date := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--supplier":
			if i+1 >= len(args) {
				return fmt.Errorf("--supplier 需要参数值")
			}
			supplierPhone = args[i+1]
			i++
		case "--amount":
			if i+1 >= len(args) {
				return fmt.Errorf("--amount 需要参数值")
			}
			amountStr = args[i+1]
			i++
		case "--date":
			if i+1 >= len(args) {
				return fmt.Errorf("--date 需要参数值")
			}
			date = args[i+1]
			i++
		}
	}

	if supplierPhone == "" {
		return fmt.Errorf("缺少 --supplier 参数(电话)")
	}
	if amountStr == "" {
		return fmt.Errorf("缺少 --amount 参数")
	}
	if date == "" {
		return fmt.Errorf("缺少 --date 参数")
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return fmt.Errorf("--amount 必须是整数: %v", err)
	}

	return AddPayout(supplierPhone, amount, date)
}

func cmdSupplierBalance(args []string) error {
	phone := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--phone":
			if i+1 >= len(args) {
				return fmt.Errorf("--phone 需要参数值")
			}
			phone = args[i+1]
			i++
		}
	}

	if phone == "" {
		return fmt.Errorf("缺少 --phone 参数")
	}

	return SupplierBalance(phone)
}

func cmdMonthly(args []string) error {
	month := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--month":
			if i+1 >= len(args) {
				return fmt.Errorf("--month 需要参数值")
			}
			month = args[i+1]
			i++
		}
	}

	if month == "" {
		return fmt.Errorf("缺少 --month 参数")
	}

	return MonthlySummary(month)
}

func printUsage() {
	fmt.Println("社区废品回收站管理工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run . <命令> [参数]")
	fmt.Println("  或编译后: ./recyclestation <命令> [参数]")
	fmt.Println()
	fmt.Println("可用命令:")
	fmt.Println()
	fmt.Println("  add-category <ID> <名称> --unit <单位> --price-per-unit <单价(分)>")
	fmt.Println("      添加新品类，单价为整数分")
	fmt.Println("      示例: go run . add-category C001 废纸 --unit kg --price-per-unit 2")
	fmt.Println()
	fmt.Println("  purchase <品类ID> --supplier <姓名> --phone <电话> --weight <重量> --date <日期>")
	fmt.Println("      添加收购记录，重量支持小数，自动计算金额(重量×单价，向下取整到分)")
	fmt.Println("      示例: go run . purchase C001 --supplier 王大爷 --phone 138xxxx --weight 15.5 --date 2024-03-20")
	fmt.Println()
	fmt.Println("  payout --supplier <电话> --amount <金额(分)> --date <日期>")
	fmt.Println("      向供应商付款，金额必须是正整数，累计付款不能超过累计应得金额")
	fmt.Println("      示例: go run . payout --supplier 138xxxx --amount 150 --date 2024-03-20")
	fmt.Println()
	fmt.Println("  supplier-balance --phone <电话>")
	fmt.Println("      查询某供应商的累计应得、已付和未付金额")
	fmt.Println("      示例: go run . supplier-balance --phone 138xxxx")
	fmt.Println()
	fmt.Println("  monthly --month <YYYY-MM>")
	fmt.Println("      查询某月各品类的回收重量和金额汇总")
	fmt.Println("      示例: go run . monthly --month 2024-03")
}
