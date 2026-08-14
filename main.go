package main

import (
	"log"

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

	log.Printf(
		"parsed %d schools, %d employees, %d payslip records",
		len(result.Schools),
		len(result.Employees),
		len(result.Payslips),
	)

	// Insert schools
	schoolIDByUDISE, err := db.InsertSchools(conn, result.Schools)
	if err != nil {
		log.Fatalf("inserting schools: %v", err)
	}

	log.Printf("inserted/updated %d schools", len(schoolIDByUDISE))

	// Insert employees
	employeeIDByShalarthID, err := db.InsertEmployees(
		conn,
		result.Employees,
		schoolIDByUDISE,
		result.UDISEByShalarthID,
	)
	if err != nil {
		log.Fatalf("inserting employees: %v", err)
	}

	log.Printf("inserted/updated %d employees", len(employeeIDByShalarthID))

	// Insert payslips
	if err := db.InsertPayslipRecords(
		conn,
		result.Payslips,
		result.ShalarthIDByPayslipIndex,
		employeeIDByShalarthID,
	); err != nil {
		log.Fatalf("inserting payslip records: %v", err)
	}

	log.Printf("inserted %d payslip records", len(result.Payslips))

	log.Println("import complete")

	schoolInfo, err := db.GetSchoolInfo(conn)
	if err != nil {
		log.Fatal(err)
	}

	employeeInfo, err := db.GetEmployeeInfoBySchoolID(conn, schoolInfo[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	payslipInfo, err := db.GetPayslipInfoByEmployeeID(conn, employeeInfo[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	err = excelgen.GeneratePayslip(
		"template.xlsx",
		"test.xlsx",
		schoolInfo[0],
		employeeInfo[0],
		payslipInfo[0],
	)
	if err != nil {
		log.Fatal(err)
	}
}
