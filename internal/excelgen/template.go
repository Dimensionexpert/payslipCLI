package excelgen

import (
	"log"

	"github.com/Dimensionexpert/payslip/internal/db"
	"github.com/Dimensionexpert/payslip/internal/pdfgen"
	"github.com/xuri/excelize/v2"
)

func GeneratePayslip(templatePath string, outputPath string, school db.SchoolResult, employee db.EmployeeResult, payslip db.PayslipResult) error {

	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		log.Printf("Error opening file:%v", f)
		return err
	}
	defer f.Close()

	f.SetCellValue("Sheet1", "C8", employee.Name)
	f.SetCellValue("Sheet1", "H8", employee.ID)

	f.SetCellValue("Sheet1", "C9", employee.Designation)
	f.SetCellValue("Sheet1", "H9", school.UDISE)

	f.SetCellValue("Sheet1", "C10", school.SchoolName)
	f.SetCellValue("Sheet1", "H10", employee.PAN)

	f.SetCellValue("Sheet1", "C11", employee.GPF)
	f.SetCellValue("Sheet1", "H11", employee.ShalarthID)

	f.SetCellValue("Sheet1", "D16", payslip.BasicPay)
	f.SetCellValue("Sheet1", "D17", payslip.DA)
	f.SetCellValue("Sheet1", "D18", payslip.HRA)
	f.SetCellValue("Sheet1", "D19", payslip.TA)
	f.SetCellValue("Sheet1", "D20", payslip.TAArrears)
	f.SetCellValue("Sheet1", "D21", payslip.DAArrears)
	f.SetCellValue("Sheet1", "D22", payslip.BasicArrears)
	f.SetCellValue("Sheet1", "D23", payslip.NPSEmprAllow)

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

	f.SetCellValue("Sheet1", "D27", payslip.GrossEarnings)
	f.SetCellValue("Sheet1", "H27", payslip.TotalDeduction)
	f.SetCellValue("Sheet1", "A31", payslip.NetPayable)

	f.SetCellValue("Sheet1", "D30", pdfgen.AmountInWords(payslip.NetPayable))

	if err := f.SaveAs(outputPath); err != nil {
		return err
	}

	return nil
}
