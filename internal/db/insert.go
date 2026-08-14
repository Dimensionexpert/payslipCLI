package db

import (
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func InsertSchools(conn *sql.DB, schools []models.School) (map[string]int64, error) {
	idByUDISE := make(map[string]int64)
	stmt, err := conn.Prepare(`
		INSERT INTO schools (udise, school_name, ddo)
		VALUES (?, ?, ?)
		ON CONFLICT(udise) DO UPDATE SET school_name = excluded.school_name, ddo = excluded.ddo
		RETURNING id
	`)
	if err != nil {
		return nil, fmt.Errorf("preparing school insert: %w", err)
	}
	defer stmt.Close()

	for _, s := range schools {
		var id int64
		err := stmt.QueryRow(s.UDISE, s.SchoolName, s.DDO).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("inserting school %s: %w", s.UDISE, err)
		}
		idByUDISE[s.UDISE] = id
	}

	return idByUDISE, nil
}

func InsertEmployees(conn *sql.DB, employees []models.Employee, schoolIDByUDISE map[string]int64, udiseByShalarthID map[string]string) (map[string]int64, error) {
	idByShalarthID := make(map[string]int64)

	stmt, err := conn.Prepare(`
		INSERT INTO employees (school_id, shalarth_id, name, gender, designation, pan, gpf, aadhaar, mobile)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(shalarth_id) DO UPDATE SET
    	name = excluded.name,
     	gender = excluded.gender,
	    designation = excluded.designation,
	    pan = excluded.pan,
	    gpf = excluded.gpf,
	    aadhaar = excluded.aadhaar,
	    mobile = excluded.mobile
		RETURNING id
	`)
	if err != nil {
		return nil, fmt.Errorf("preparing employee insert: %w", err)
	}
	defer stmt.Close()

	for _, e := range employees {
		udise := udiseByShalarthID[e.ShalarthID]
		schoolID, ok := schoolIDByUDISE[udise]
		if !ok {
			return nil, fmt.Errorf("employee %s: no matching school for udise %s", e.ShalarthID, udise)
		}

		var id int64
		err := stmt.QueryRow(schoolID, e.ShalarthID, e.Name, e.Gender, e.Designation, e.PanNo, e.GPF, e.AadharNo, e.MobileNo).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("inserting employee %s: %w", e.ShalarthID, err)
		}
		idByShalarthID[e.ShalarthID] = id
	}

	return idByShalarthID, nil
}

func InsertPayslipRecords(conn *sql.DB, payslips []models.PayslipRecord, shalarthIDByPayslipIndex []string, employeeIDByShalarthID map[string]int64) error {
	stmt, err := conn.Prepare(`
		INSERT INTO payslip_records (
			employee_id, month, year,
			basic_pay, da, hra, ta, ta_arrears, da_arrears, basic_arrears, nps_empr_allow,
			fa, gpf_deduction, gpf_advance, pt, gis, revenue_stamp, nps_empr_contri, nps_emp_contri,
			income_tax, ngr_society_loan, home_loan,
			gross_earnings, total_deduction, net_payable
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(employee_id, month, year) DO UPDATE SET
			basic_pay = excluded.basic_pay,
			da = excluded.da,
			hra = excluded.hra,
			ta = excluded.ta,
			ta_arrears = excluded.ta_arrears,
			da_arrears = excluded.da_arrears,
			basic_arrears = excluded.basic_arrears,
			nps_empr_allow = excluded.nps_empr_allow,
			fa = excluded.fa,
			gpf_deduction = excluded.gpf_deduction,
			gpf_advance = excluded.gpf_advance,
			pt = excluded.pt,
			gis = excluded.gis,
			revenue_stamp = excluded.revenue_stamp,
			nps_empr_contri = excluded.nps_empr_contri,
			nps_emp_contri = excluded.nps_emp_contri,
			income_tax = excluded.income_tax,
			ngr_society_loan = excluded.ngr_society_loan,
			home_loan = excluded.home_loan,
			gross_earnings = excluded.gross_earnings,
			total_deduction = excluded.total_deduction,
			net_payable = excluded.net_payable
	`)
	if err != nil {
		return fmt.Errorf("preparing payslip insert: %w", err)
	}
	defer stmt.Close()

	for i, p := range payslips {
		shalarthID := shalarthIDByPayslipIndex[i]
		employeeID, ok := employeeIDByShalarthID[shalarthID]
		if !ok {
			return fmt.Errorf("payslip row %d: no matching employee for shalarth_id %s", i, shalarthID)
		}

		_, err := stmt.Exec(
			employeeID, p.Month, p.Year,
			p.BasicPay, p.DA, p.HRA, p.TA, p.TAArrears, p.DAArrears, p.BasicArrears, p.NPSEmprAllow,
			p.FA, p.GPFDeduction, p.GPFAdvance, p.PT, p.GIS, p.RevenueStamp, p.NPSEmprContri, p.NPSEmpContri,
			p.IncomeTax, p.NGRSocietyLoan, p.HomeLoan,
			p.GrossEarnings, p.TotalDeduction, p.NetPayable,
		)
		if err != nil {
			return fmt.Errorf("inserting payslip for %s: %w", shalarthID, err)
		}
	}

	return nil
}
