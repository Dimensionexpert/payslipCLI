package pdfgen

import (
	"strings"
)

var ones = []string{
	"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine",
	"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen",
	"Seventeen", "Eighteen", "Nineteen",
}

var tens = []string{
	"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety",
}

// twoDigitWords converts a number 0-99 to words.
func twoDigitWords(n int) string {
	if n < 20 {
		return ones[n]
	}
	word := tens[n/10]
	if n%10 != 0 {
		word += " " + ones[n%10]
	}
	return word
}

// threeDigitWords converts a number 0-999 to words (used for the first group).
func threeDigitWords(n int) string {
	if n >= 100 {
		hundredsWord := ones[n/100] + " Hundred"
		remainder := n % 100
		if remainder == 0 {
			return hundredsWord
		}
		return hundredsWord + " " + twoDigitWords(remainder)
	}
	return twoDigitWords(n)
}

// AmountInWords converts a rupee amount to words using the Indian numbering
// system: groups are 3 digits (hundreds), then 2 digits each after that
// (thousand, lakh, crore) - unlike the international 3-3-3 grouping.
func AmountInWords(amount float64) string {
	n := int64(amount) // rupees only, ignoring paise for payslip purposes

	if n == 0 {
		return "Rupees Zero Only"
	}

	crore := n / 10000000
	n %= 10000000
	lakh := n / 100000
	n %= 100000
	thousand := n / 1000
	n %= 1000
	hundred := n

	var parts []string

	if crore > 0 {
		parts = append(parts, twoDigitWords(int(crore))+" Crore")
	}
	if lakh > 0 {
		parts = append(parts, twoDigitWords(int(lakh))+" Lakh")
	}
	if thousand > 0 {
		parts = append(parts, twoDigitWords(int(thousand))+" Thousand")
	}
	if hundred > 0 {
		parts = append(parts, threeDigitWords(int(hundred)))
	}

	return "Rupees " + strings.Join(parts, " ") + " Only"
}
