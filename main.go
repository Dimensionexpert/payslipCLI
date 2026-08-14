package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/importer"
	"github.com/Dimensionexpert/payslip/internal/pdfgen"
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
	//ParseExcel
	result, err := importer.ParseExcel("July26.xlsx", 7, 2026)
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}
	log.Printf("parsed %d schools, %d employees, %d payslip records",
		len(result.Schools), len(result.Employees), len(result.Payslips))
	log.Println("parse successful")

	//InsertSchool
	schoolIDByUDISE, err := db.InsertSchools(conn, result.Schools)
	if err != nil {
		log.Fatalf("inserting schools: %v", err)
	}
	log.Printf("inserted/updated %d schools", len(schoolIDByUDISE))

	//InsertEmployee
	employeeIDByShalarthID, err := db.InsertEmployees(conn, result.Employees, schoolIDByUDISE, result.UDISEByShalarthID)
	if err != nil {
		log.Fatalf("inserting employees: %v", err)
	}
	log.Printf("inserted/updated %d employees", len(employeeIDByShalarthID))

	//InsertPayslip
	if err := db.InsertPayslipRecords(conn, result.Payslips, result.ShalarthIDByPayslipIndex, employeeIDByShalarthID); err != nil {
		log.Fatalf("inserting payslip records: %v", err)
	}
	log.Printf("inserted %d payslip records", len(result.Payslips))

	// quick manual sanity checks, remove once wired into the real pipeline
	AmountIndWords := pdfgen.AmountInWords(1224554)
	fmt.Printf("Amount is :%s", AmountIndWords)

	log.Println("import complete")
}
