package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/excelgen"
	"github.com/Dimensionexpert/payslip/internal/importer"
	"github.com/Dimensionexpert/payslip/internal/pdfgen"
)

func main() {
	start := time.Now()
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

	conn, err := db.OpenDB("payslip.db")
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer conn.Close()

	// Parse Excel
	result, err := importer.ParseExcel(*excelPath, *month, *year)
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}
	log.Printf("parsed %d schools, %d employees, %d payslip records",
		len(result.Schools), len(result.Employees), len(result.Payslips))

	// Insert schools
	schoolIDByUDISE, err := db.InsertSchools(conn, result.Schools)
	if err != nil {
		log.Fatalf("inserting schools: %v", err)
	}
	log.Printf("inserted/updated %d schools", len(schoolIDByUDISE))

	// Insert employees
	employeeIDByShalarthID, err := db.InsertEmployees(conn, result.Employees, schoolIDByUDISE, result.UDISEByShalarthID)
	if err != nil {
		log.Fatalf("inserting employees: %v", err)
	}
	log.Printf("inserted/updated %d employees", len(employeeIDByShalarthID))

	// Insert payslips
	if err := db.InsertPayslipRecords(conn, result.Payslips, result.ShalarthIDByPayslipIndex, employeeIDByShalarthID); err != nil {
		log.Fatalf("inserting payslip records: %v", err)
	}
	log.Printf("inserted %d payslip records", len(result.Payslips))
	log.Println("import complete")

	// Fetch joined school+employee+payslip data in one query
	payslipData, err := db.GetPayslipData(conn)
	if err != nil {
		log.Fatalf("fetching payslip data: %v", err)
	}
	log.Printf("fetched %d payslip export records", len(payslipData))

	// Ensure the output directory exists before writing any files into it
	if err := os.MkdirAll("output/pdf", 0755); err != nil {
		log.Fatalf("creating pdf output dir: %v", err)
	}

	for _, entry := range payslipData {
		xlsxPath, err := excelgen.GeneratePayslip("template.xlsx", "output", entry)
		if err != nil {
			log.Fatalf("generating payslip for %s: %v", entry.Employee.Name, err)
		}

		if err := pdfgen.ConvertToPDF(xlsxPath, "output/pdf"); err != nil {
			log.Fatalf("converting pdf for %s: %v", entry.Employee.Name, err)
		}
	}
	log.Printf("generated %d payslips (xlsx + pdf)", len(payslipData))
	log.Printf("total time: %s", time.Since(start))
}
