package db

import (
	"database/sql"
	"fmt"
	"sort"
)

type SchoolResult struct {
	ID         int
	SchoolName string
	UDISE      string
	DDO        string
}

type EmployeeResult struct {
	ID          int
	SchoolID    int
	ShalarthID  string
	Name        string
	Gender      string
	Designation string
	PAN         string
	GPF         string
	Aadhaar     string
	Mobile      string
}

type PayslipResult struct {
	ID             int
	EmployeeID     int
	Month          int
	Year           int
	BasicPay       float64
	DA             float64
	HRA            float64
	TA             float64
	TAArrears      float64
	DAArrears      float64
	BasicArrears   float64
	NPSEmprAllow   float64
	FA             float64
	GPFDeduction   float64
	GPFAdvance     float64
	PT             float64
	GIS            float64
	RevenueStamp   float64
	NPSEmprContri  float64
	NPSEmpContri   float64
	IncomeTax      float64
	NGRSocietyLoan float64
	HomeLoan       float64
	GrossEarnings  float64
	TotalDeduction float64
	NetPayable     float64
}

func GetSchoolInfo(conn *sql.DB) ([]SchoolResult, error) {
	stmt, err := conn.Prepare(`SELECT id, udise, school_name, ddo FROM schools`)
	if err != nil {
		return nil, fmt.Errorf("preparing school info: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("querying school info: %w", err)
	}
	defer rows.Close()

	var results []SchoolResult
	var udise, schoolName, ddo string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &udise, &schoolName, &ddo)
		if err != nil {
			return nil, fmt.Errorf("scanning school row: %w", err)
		}

		results = append(results, SchoolResult{
			ID:         id,
			UDISE:      udise,
			SchoolName: schoolName,
			DDO:        ddo,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating school rows: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results, nil
}

func GetEmployeeInfoBySchoolID(conn *sql.DB, schoolId int) ([]EmployeeResult, error) {
	stmt, err := conn.Prepare(`SELECT * FROM employees WHERE school_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("preparing school info: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(schoolId)
	if err != nil {
		return nil, fmt.Errorf("querying school info: %w", err)
	}
	defer rows.Close()

	var results []EmployeeResult
	var id, schoolID int
	var shalarthID, name, gender, designation, pan, gpf, aadhaar, mobile string

	for rows.Next() {
		err := rows.Scan(
			&id,
			&schoolID,
			&shalarthID,
			&name,
			&gender,
			&designation,
			&pan,
			&gpf,
			&aadhaar,
			&mobile,
		)

		if err != nil {
			return nil, fmt.Errorf("iterating school rows: %w", err)
		}

		results = append(results, EmployeeResult{
			ID:          id,
			SchoolID:    schoolID,
			ShalarthID:  shalarthID,
			Name:        name,
			Gender:      gender,
			Designation: designation,
			PAN:         pan,
			GPF:         gpf,
			Aadhaar:     aadhaar,
			Mobile:      mobile,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating school rows: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})
	return results, nil

}

func GetPayslipInfoByEmployeeID(conn *sql.DB, employeeID int) ([]PayslipResult, error) {
	rows, err := conn.Query(`
		SELECT
			id,
			employee_id,
			month,
			year,
			basic_pay,
			da,
			hra,
			ta,
			ta_arrears,
			da_arrears,
			basic_arrears,
			nps_empr_allow,
			fa,
			gpf_deduction,
			gpf_advance,
			pt,
			gis,
			revenue_stamp,
			nps_empr_contri,
			nps_emp_contri,
			income_tax,
			ngr_society_loan,
			home_loan,
			gross_earnings,
			total_deduction,
			net_payable
		FROM payslip_records
		WHERE employee_id = ?
	`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("querying payslip records: %w", err)
	}
	defer rows.Close()

	var results []PayslipResult

	for rows.Next() {
		var p PayslipResult

		err := rows.Scan(
			&p.ID,
			&p.EmployeeID,
			&p.Month,
			&p.Year,
			&p.BasicPay,
			&p.DA,
			&p.HRA,
			&p.TA,
			&p.TAArrears,
			&p.DAArrears,
			&p.BasicArrears,
			&p.NPSEmprAllow,
			&p.FA,
			&p.GPFDeduction,
			&p.GPFAdvance,
			&p.PT,
			&p.GIS,
			&p.RevenueStamp,
			&p.NPSEmprContri,
			&p.NPSEmpContri,
			&p.IncomeTax,
			&p.NGRSocietyLoan,
			&p.HomeLoan,
			&p.GrossEarnings,
			&p.TotalDeduction,
			&p.NetPayable,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning payslip row: %w", err)
		}

		results = append(results, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating payslip rows: %w", err)
	}

	return results, nil
}
