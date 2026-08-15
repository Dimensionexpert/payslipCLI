package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/excelgen"
	"github.com/Dimensionexpert/payslip/internal/importer"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real env vars")
	}

	conn, err := db.OpenDB("payslip.db")
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer conn.Close()

	// Parse Excel
	result, err := importer.ParseExcel("July26.xlsx", 7, 2026)
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
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("creating output dir: %v", err)
	}

	// Generate one payslip per employee
	for _, entry := range payslipData {
		if err := excelgen.GeneratePayslip("template.xlsx", "output", entry); err != nil {
			log.Fatalf("generating payslip for %s: %v", entry.Employee.Name, err)
		}
	}
	log.Printf("generated %d payslips", len(payslipData))
}
