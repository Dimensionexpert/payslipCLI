package db

import (
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/models"
)

type PayslipExport struct {
	School   models.School
	Employee models.Employee
	Payslip  models.PayslipRecord
}

func GetPayslipData(conn *sql.DB) ([]PayslipExport, error) {
	rows, err := conn.Query(`SELECT
    s.id AS school_id,
    s.school_name,
    s.udise,
    s.ddo,

    e.id AS employee_id,
    e.school_id,
    e.shalarth_id,
    e.name,
    e.gender,
    e.designation,
    e.pan,
    e.gpf,
    e.aadhaar,
    e.mobile,

    p.id AS payslip_id,
    p.employee_id,
    p.month,
    p.year,
    p.basic_pay,
    p.da,
    p.hra,
    p.ta,
    p.ta_arrears,
    p.da_arrears,
    p.basic_arrears,
    p.nps_empr_allow,
    p.fa,
    p.gpf_deduction,
    p.gpf_advance,
    p.pt,
    p.gis,
    p.revenue_stamp,
    p.nps_empr_contri,
    p.nps_emp_contri,
    p.income_tax,
    p.ngr_society_loan,
    p.home_loan,
    p.gross_earnings,
    p.total_deduction,
    p.net_payable

FROM schools s
JOIN employees e
    ON e.school_id = s.id
JOIN payslip_records p
    ON p.employee_id = e.id

ORDER BY
    s.school_name,
    e.name;`)
	if err != nil {
		return nil, fmt.Errorf("querying payslip data: %w", err)
	}
	defer rows.Close()

	var results []PayslipExport

	for rows.Next() {
		var pe PayslipExport

		err := rows.Scan(
			&pe.School.ID, &pe.School.SchoolName, &pe.School.UDISE, &pe.School.DDO,
			&pe.Employee.ID, &pe.Employee.SchoolID, &pe.Employee.ShalarthID, &pe.Employee.Name,
			&pe.Employee.Gender, &pe.Employee.Designation, &pe.Employee.PanNo, &pe.Employee.GPF,
			&pe.Employee.AadharNo, &pe.Employee.MobileNo,
			&pe.Payslip.ID, &pe.Payslip.EmployeeID, &pe.Payslip.Month, &pe.Payslip.Year,
			&pe.Payslip.BasicPay, &pe.Payslip.DA, &pe.Payslip.HRA, &pe.Payslip.TA,
			&pe.Payslip.TAArrears, &pe.Payslip.DAArrears, &pe.Payslip.BasicArrears, &pe.Payslip.NPSEmprAllow,
			&pe.Payslip.FA, &pe.Payslip.GPFDeduction, &pe.Payslip.GPFAdvance, &pe.Payslip.PT,
			&pe.Payslip.GIS, &pe.Payslip.RevenueStamp, &pe.Payslip.NPSEmprContri, &pe.Payslip.NPSEmpContri,
			&pe.Payslip.IncomeTax, &pe.Payslip.NGRSocietyLoan, &pe.Payslip.HomeLoan,
			&pe.Payslip.GrossEarnings, &pe.Payslip.TotalDeduction, &pe.Payslip.NetPayable,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning payslip row: %w", err)
		}

		results = append(results, pe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating payslip rows: %w", err)
	}

	return results, nil
}
