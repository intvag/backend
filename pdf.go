package main

import (
	"fmt"
	"os"

	"github.com/go-pdf/fpdf"
)

func (p Policy) GeneratePDF(person Person) (fn string, err error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetTopMargin(30)

	pdf.SetHeaderFuncMode(func() {
		pdf.Image("assets/logo.png", 10, 6, 30, 0, false, "", 0, "")
	}, true)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})

	pdf.AddPage()

	pdf.Ln(20)

	pdf.SetFont("Helvetica", "", 12)
	_, lineHt := pdf.GetFontSize()

	price := 0.0
	for _, i := range p.Items {
		price += i.Cost
	}

	pdf.Write(lineHt, person.Name)
	pdf.Ln(lineHt * 1.5)

	pdf.Write(lineHt, person.Address)
	pdf.Ln(lineHt * 1.5)

	pdf.Write(lineHt, person.Postcode)

	pdf.Ln(lineHt * 6)

	pdf.Write(lineHt, fmt.Sprintf("Dear %s, ", person.Name))

	pdf.Ln(lineHt * 3)

	pdf.Write(lineHt, fmt.Sprintf("You owe us GBP %.2f per month for your crap.", price))
	pdf.Ln(lineHt * 1.5)
	pdf.Write(lineHt, "In return we promise to fix your crap when it breaks, or even replace it.")

	pdf.Ln(lineHt * 3)

	pdf.Write(lineHt, fmt.Sprintf("Sincerely,"))
	pdf.Ln(lineHt * 1.5)

	pdf.Image("assets/signature.png", pdf.GetX()+6, pdf.GetY(), lineHt*7, 0, true, "", 0, "")
	pdf.Ln(lineHt * 1.5)

	pdf.Write(lineHt, "Some Dude, CEO International and Vague")

	f, err := os.CreateTemp("", "policy-*.pdf")
	if err != nil {
		return
	}

	return f.Name(), pdf.OutputFileAndClose(f.Name())
}
