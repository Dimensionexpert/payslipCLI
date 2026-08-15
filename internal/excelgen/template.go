package excelgen

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/pdfgen"
)

func sanitizeFilename(s string) string {
	return strings.Join(strings.Fields(s), "_")
}

func GeneratePayslip(templatePath, outputDir string, data db.PayslipExport) error {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return fmt.Errorf("opening template: %w", err)
	}
	defer f.Close()

	school := data.School
	employee := data.Employee
	payslip := data.Payslip

	// Identity block
	f.SetCellValue("Sheet1", "C8", employee.Name)
	f.SetCellValue("Sheet1", "H8", employee.ID)
	f.SetCellValue("Sheet1", "C9", employee.Designation)
	f.SetCellValue("Sheet1", "H9", school.UDISE)
	f.SetCellValue("Sheet1", "C10", school.SchoolName)
	f.SetCellValue("Sheet1", "H10", employee.PanNo)
	f.SetCellValue("Sheet1", "C11", employee.GPF)
	f.SetCellValue("Sheet1", "G11", employee.ShalarthID)

	// Earnings
	f.SetCellValue("Sheet1", "D16", payslip.BasicPay)
	f.SetCellValue("Sheet1", "D17", payslip.DA)
	f.SetCellValue("Sheet1", "D18", payslip.HRA)
	f.SetCellValue("Sheet1", "D19", payslip.TA)
	f.SetCellValue("Sheet1", "D20", payslip.TAArrears)
	f.SetCellValue("Sheet1", "D21", payslip.DAArrears)
	f.SetCellValue("Sheet1", "D22", payslip.BasicArrears)
	f.SetCellValue("Sheet1", "D23", payslip.NPSEmprAllow)

	// Deductions
	f.SetCellValue("Sheet1", "H16", payslip.FA)
	f.SetCellValue("Sheet1", "H17", payslip.GPFDeduction)
	f.SetCellValue("Sheet1", "H18", payslip.GPFAdvance)
	f.SetCellValue("Sheet1", "H19", payslip.PT)
	f.SetCellValue("Sheet1", "H20", payslip.GIS)
	f.SetCellValue("Sheet1", "H21", payslip.RevenueStamp)
	f.SetCellValue("Sheet1", "H22", payslip.NPSEmprContri)
	f.SetCellValue("Sheet1", "H23", payslip.NPSEmpContri)
	f.SetCellValue("Sheet1", "H24", payslip.IncomeTax)
	f.SetCellValue("Sheet1", "H25", payslip.NGRSocietyLoan)
	f.SetCellValue("Sheet1", "H26", payslip.HomeLoan)

	// Totals
	f.SetCellValue("Sheet1", "D27", payslip.GrossEarnings)
	f.SetCellValue("Sheet1", "H27", payslip.TotalDeduction)
	f.SetCellValue("Sheet1", "A30", fmt.Sprintf("₹  %.0f", payslip.NetPayable))
	f.SetCellValue("Sheet1", "D30", pdfgen.AmountInWords(payslip.NetPayable))

	// A33 ("PLACE") left untouched — confirmed fixed office name in the template
	f.SetCellValue("Sheet1", "E32", time.Now().Format("02-Jan-06"))

	filename := fmt.Sprintf("%s_%s.xlsx",
		sanitizeFilename(school.SchoolName),
		sanitizeFilename(employee.Name),
	)
	outputPath := filepath.Join(outputDir, filename)

	if err := f.SaveAs(outputPath); err != nil {
		return fmt.Errorf("saving output: %w", err)
	}
	return nil
}
