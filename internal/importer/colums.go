package importer

// Column indices in the government xlsx export (0-indexed, per GetRows).
const (
	colUDISE       = 1
	colSchool      = 2 //school name
	colDDO         = 3
	colName        = 5 //employee Name
	colShalarthID  = 6
	colGender      = 7
	colDesignation = 8
	colGPF         = 9
	colPAN         = 10
	colAadhaar     = 11
	colMobile      = 12

	colBasicPay      = 13
	colDA            = 14
	colHRA           = 15
	colTA            = 16
	colTAArrears     = 17
	colDAArrears     = 18
	colBasicArrears  = 19
	colNPSEmprAllow  = 20
	colGrossEarnings = 21

	colFA             = 22
	colGPFDeduction   = 23
	colGPFAdvance     = 24
	colPT             = 25
	colGIS            = 26
	colRevenueStamp   = 27
	colNPSEmprContri  = 28
	colNPSEmpContri   = 29
	colIncomeTax      = 30
	colNGRSocietyLoan = 31
	colHomeLoan       = 32
	colTotalDeduction = 33
	colNetPayable     = 34
)
