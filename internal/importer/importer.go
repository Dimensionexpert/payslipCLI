package importer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func OpenExcel(path string) (*excelize.File, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening excel file: %w", err)
	}
	return f, nil
}

func cleanCell(s string) string {
	return strings.TrimSpace(s)
}

func cleanNumericString(s string) string {
	s = cleanCell(s)
	if strings.ContainsAny(s, "eE") {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return strconv.FormatFloat(f, 'f', 0, 64)
		}
	}
	return s
}

// parseFloat converts a cell's string value to float64. A blank cell
// becomes 0 (some earning/deduction columns are legitimately empty for a
// given teacher). Anything non-blank that fails to parse is a real error —
// we don't guess, we report and stop.
func parseFloat(s string) (float64, error) {
	s = cleanCell(s)
	if s == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", s, err)
	}
	return v, nil
}

type ParseResult struct {
	Schools   []models.School
	Employees []models.Employee
	Payslips  []models.PayslipRecord

	UDISEByShalarthID        map[string]string
	ShalarthIDByPayslipIndex []string
}

func ParseExcel(path string, month, year int) (*ParseResult, error) {
	f, err := OpenExcel(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}

	schoolsByUDISE := make(map[string]models.School)
	employeeByShalarthID := make(map[string]models.Employee)
	udiseByShalarthID := make(map[string]string)

	var payslips []models.PayslipRecord
	var shalarthIDByPayslipIndex []string

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1

		udise := cleanCell(row[colUDISE])
		if udise == "" {
			return nil, fmt.Errorf("row %d: missing UDISE", rowNum)
		}

		shalarthID := cleanCell(row[colShalarthID])
		if shalarthID == "" {
			return nil, fmt.Errorf("row %d: missing ShalarthID", rowNum)
		}

		pan := cleanCell(row[colPAN])
		if pan == "" {
			return nil, fmt.Errorf("row %d: missing Pan", rowNum)
		}

		aadhar := cleanNumericString(row[colAadhaar])
		if aadhar == "" {
			return nil, fmt.Errorf("row %d: missing Aadhar", rowNum)
		}

		if _, exist := schoolsByUDISE[udise]; !exist {
			schoolsByUDISE[udise] = models.School{
				UDISE:      udise,
				SchoolName: cleanCell(row[colSchool]),
				DDO:        cleanCell(row[colDDO]),
			}
		}

		if _, exist := employeeByShalarthID[shalarthID]; !exist {
			employeeByShalarthID[shalarthID] = models.Employee{
				ShalarthID:  shalarthID,
				Name:        cleanCell(row[colName]),
				Gender:      cleanCell(row[colGender]),
				Designation: cleanCell(row[colDesignation]),
				PanNo:       pan,
				GPF:         cleanCell(row[colGPF]),
				AadharNo:    aadhar,
				MobileNo:    cleanCell(row[colMobile]),
			}
		}

		udiseByShalarthID[shalarthID] = udise

		payslip, err := parsePayslipRow(row, month, year)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}
		payslips = append(payslips, payslip)
		shalarthIDByPayslipIndex = append(shalarthIDByPayslipIndex, shalarthID)
	}

	schools := make([]models.School, 0, len(schoolsByUDISE))
	for _, s := range schoolsByUDISE {
		schools = append(schools, s)
	}

	employees := make([]models.Employee, 0, len(employeeByShalarthID))
	for _, e := range employeeByShalarthID {
		employees = append(employees, e)
	}

	return &ParseResult{
		Schools:                  schools,
		Employees:                employees,
		Payslips:                 payslips,
		UDISEByShalarthID:        udiseByShalarthID,
		ShalarthIDByPayslipIndex: shalarthIDByPayslipIndex,
	}, nil
}

// parsePayslipRow builds one PayslipRecord from a single row's money
// columns. month/year are passed in since the sheet doesn't state them.
func parsePayslipRow(row []string, month, year int) (models.PayslipRecord, error) {
	var p models.PayslipRecord
	p.Month = month
	p.Year = year

	fields := []struct {
		col int
		dst *float64
	}{
		{colBasicPay, &p.BasicPay},
		{colDA, &p.DA},
		{colHRA, &p.HRA},
		{colTA, &p.TA},
		{colTAArrears, &p.TAArrears},
		{colDAArrears, &p.DAArrears},
		{colBasicArrears, &p.BasicArrears},
		{colNPSEmprAllow, &p.NPSEmprAllow},
		{colFA, &p.FA},
		{colGPFDeduction, &p.GPFDeduction},
		{colGPFAdvance, &p.GPFAdvance},
		{colPT, &p.PT},
		{colGIS, &p.GIS},
		{colRevenueStamp, &p.RevenueStamp},
		{colNPSEmprContri, &p.NPSEmprContri},
		{colNPSEmpContri, &p.NPSEmpContri},
		{colIncomeTax, &p.IncomeTax},
		{colNGRSocietyLoan, &p.NGRSocietyLoan},
		{colHomeLoan, &p.HomeLoan},
		{colGrossEarnings, &p.GrossEarnings},
		{colTotalDeduction, &p.TotalDeduction},
		{colNetPayable, &p.NetPayable},
	}

	for _, f := range fields {
		v, err := parseFloat(row[f.col])
		if err != nil {
			return p, err
		}
		*f.dst = v
	}

	return p, nil
}
