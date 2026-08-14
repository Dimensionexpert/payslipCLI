package models

type School struct {
	ID         int64
	UDISE      string
	SchoolName string
	DDO        string
}

type Employee struct {
	ID          int64
	SchoolID    int64
	ShalarthID  string
	Name        string
	Gender      string
	Designation string
	PanNo       string
	GPF         string
	AadharNo    string
	MobileNo    string
}

type PayslipRecord struct {
	ID         int64
	EmployeeID int64
	Month      int
	Year       int

	// Earnings
	BasicPay     float64 // BASIC PAY
	DA           float64 // D.A.
	HRA          float64 // H.R.A.
	TA           float64 // T.A.
	TAArrears    float64 // T.A. ARREARS
	DAArrears    float64 // DA ARREARS
	BasicArrears float64 // BASIC ARREARS
	NPSEmprAllow float64 // NPS EMPR ALLOW

	// Deductions
	FA             float64 // F.A.
	GPFDeduction   float64 // G.P.F.
	GPFAdvance     float64 // GPF ADVANCE
	PT             float64 // P.T.
	GIS            float64 // G.I.S.
	RevenueStamp   float64 // REVENUE STAMP
	NPSEmprContri  float64 // NPS EMPR CONTRI
	NPSEmpContri   float64 // NPS EMP CONTRI
	IncomeTax      float64 // INCOME TAX
	NGRSocietyLoan float64 // NGR (SOCIETY LOAN)
	HomeLoan       float64 // HOME LOAN

	// Computed totals — stored verbatim from the excel, not recalculated
	// (per your earlier call: gov's numbers, gov's problem)
	GrossEarnings  float64 // GROSS EARNINGS
	TotalDeduction float64 // TOTAL DEDUCTION
	NetPayable     float64 // TOTAL NET PAYBLE
}
