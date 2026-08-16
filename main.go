package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/excelgen"
	"github.com/Dimensionexpert/payslip/internal/importer"
	"github.com/Dimensionexpert/payslip/internal/timer"
)

func main() {
	defer timer.Track("total execution")()

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real env vars")
	}

	excelPath := flag.String("f", "", "path to the govt xlsx export")
	month := flag.Int("m", 0, "pay period month (1-12)")
	year := flag.Int("y", 0, "pay period year")
	flag.Parse()

	if *excelPath == "" || *month == 0 || *year == 0 {
		log.Fatal("missing required flags — example: ./payslip -f=Aug26.xlsx -m=8 -y=2026")
	}

	// Open database
	done := timer.Track("database open")
	conn, err := db.OpenDB("payslip.db")
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	done()
	defer conn.Close()

	// Parse Excel
	done = timer.Track("excel parsing")
	result, err := importer.ParseExcel(*excelPath, *month, *year)
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}
	done()
	log.Printf("parsed %d schools, %d employees, %d payslip records",
		len(result.Schools), len(result.Employees), len(result.Payslips))

	// Insert schools
	done = timer.Track("school insertion")
	schoolIDByUDISE, err := db.InsertSchools(conn, result.Schools)
	if err != nil {
		log.Fatalf("inserting schools: %v", err)
	}
	done()
	log.Printf("inserted/updated %d schools", len(schoolIDByUDISE))

	// Insert employees
	done = timer.Track("employee insertion")
	employeeIDByShalarthID, err := db.InsertEmployees(conn, result.Employees, schoolIDByUDISE, result.UDISEByShalarthID)
	if err != nil {
		log.Fatalf("inserting employees: %v", err)
	}
	done()
	log.Printf("inserted/updated %d employees", len(employeeIDByShalarthID))

	// Insert payslips
	done = timer.Track("payslip insertion")
	if err := db.InsertPayslipRecords(conn, result.Payslips, result.ShalarthIDByPayslipIndex, employeeIDByShalarthID); err != nil {
		log.Fatalf("inserting payslip records: %v", err)
	}
	done()
	log.Printf("inserted %d payslip records", len(result.Payslips))
	log.Println("import complete")

	// Fetch joined school + employee + payslip data
	done = timer.Track("payslip data fetch")
	payslipData, err := db.GetPayslipData(conn)
	if err != nil {
		log.Fatalf("fetching payslip data: %v", err)
	}
	done()
	log.Printf("fetched %d payslip export records", len(payslipData))

	// Create output directories
	if err := os.MkdirAll("output/pdf", 0755); err != nil {
		log.Fatalf("creating pdf output dir: %v", err)
	}

	// Generate XLSX (sequential — fast, no benefit from concurrency here)
	done = timer.Track("xlsx generation")
	var xlsxPaths []string
	for _, entry := range payslipData {
		xlsxPath, err := excelgen.GeneratePayslip("template.xlsx", "output", entry)
		if err != nil {
			log.Fatalf("generating payslip for %s: %v", entry.Employee.Name, err)
		}
		xlsxPaths = append(xlsxPaths, xlsxPath)
	}
	done()

	// Convert to PDF using a worker pool (2 workers, given sisyphus's 4GB RAM)
	done = timer.Track("pdf conversion")
	results := concurrency.RunPDFConversion(xlsxPaths, 4)
	done()

	failures := 0
	for _, r := range results {
		if r.Err != nil {
			log.Printf("FAILED: %s: %v", r.Filepath, r.Err)
			failures++
		}
	}
	log.Printf("generated %d payslips (%d failed)", len(results), failures)
}
