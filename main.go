package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/excelgen"
	"github.com/Dimensionexpert/payslip/internal/importer"
	"github.com/Dimensionexpert/payslip/internal/timer"
)

const (
	inboxDir  = "/srv/payslip/inbox"
	outputDir = "/srv/payslip/output"
)

func parseMonthYearFromFilename(filename string) (int, int, error) {
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	if len(base) != 4 {
		return 0, 0, fmt.Errorf("expected filename like 0726.xlsx (MMYY), got %q", filename)
	}
	month, err := strconv.Atoi(base[0:2])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month in filename %q", filename)
	}
	yearSuffix, err := strconv.Atoi(base[2:4])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year in filename %q", filename)
	}
	return month, 2000 + yearSuffix, nil
}

func main() {
	defer timer.Track("total execution")()

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real env vars")
	}

	filename := flag.String("f", "", "filename in inbox, e.g 0726.xlsx")
	flag.Parse()

	if *filename == "" {
		log.Fatal("missing required flag — example: ./payslip -f=0726.xlsx")
	}

	month, year, err := parseMonthYearFromFilename(*filename)
	if err != nil {
		log.Fatalf("bad filename: %v", err)
	}

	excelPath := filepath.Join(inboxDir, *filename)

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
	result, err := importer.ParseExcel(excelPath, month, year)
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
	if err := os.MkdirAll(filepath.Join(outputDir, "pdf"), 0755); err != nil {
		log.Fatalf("creating output dir: %v", err)
	}

	// Generate XLSX (sequential — fast, no benefit from concurrency here)
	done = timer.Track("xlsx generation")
	var xlsxPaths []string
	for _, entry := range payslipData {
		xlsxPath, err := excelgen.GeneratePayslip("template.xlsx", outputDir, entry)
		if err != nil {
			log.Fatalf("generating payslip for %s: %v", entry.Employee.Name, err)
		}
		xlsxPaths = append(xlsxPaths, xlsxPath)
	}
	done()

	// Convert to PDF using a worker pool
	done = timer.Track("pdf conversion")
	results := concurrency.RunPDFConversion(xlsxPaths, 4, filepath.Join(outputDir, "pdf"))
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
