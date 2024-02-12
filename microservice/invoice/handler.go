package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

type Order struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (app *application) CreateAndSendInvoice(w http.ResponseWriter, r *http.Request) {
	// receive json
	var order Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		BadRequest(w, r, err)
		return
	}

	// order := Order{
	// 	ID:        1,
	// 	Quantity:  3,
	// 	Amount:    1000,
	// 	Product:   "Widget",
	// 	FirstName: "Nayeem",
	// 	LastName:  "Nishaat",
	// 	Email:     "nayeem@nishaats",
	// 	CreatedAt: time.Now(),
	// }

	// generate a pdf invoice
	err = app.createInvoicePDF(order)
	if err != nil {
		BadRequest(w, r, err)
		return
	}

	// create mail attachment
	attachments := []string{
		fmt.Sprintf("./pdf/%d.pdf", order.ID),
	}

	// send mail with attachment
	err = app.SendMail(app.SMTP.From, order.Email, "Your Invoice", "invoice", attachments, nil)
	if err != nil {
		BadRequest(w, r, err)
		return
	}

	// delete attachments
	for _, attachment := range attachments {
		if err := os.Remove(attachment); err != nil {
			fmt.Println("Failed to remove file!", err)
		}
	}

	// send response
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = fmt.Sprintf("Invoice %d.pdf created and sent to %s.", order.ID, order.Email)

	WriteJSON(w, http.StatusCreated, resp)
}

func (app *application) createInvoicePDF(order Order) error {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)
	pdf.SetAutoPageBreak(true, 0)

	importer := gofpdi.NewImporter()

	t := importer.ImportPage(pdf, "./pdf/template.pdf", 1, "/MediaBox")

	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	// write info
	pdf.SetY(50)
	pdf.SetX(10)
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(97, 8, fmt.Sprintf("Attention: %s %s", order.FirstName, order.LastName), "", 0, "L", false, 0, "")

	pdf.Ln(5)
	pdf.CellFormat(97, 8, order.Email, "", 0, "L", false, 0, "")

	pdf.Ln(5)
	pdf.CellFormat(97, 8, order.CreatedAt.Format("2006-01-01"), "", 0, "L", false, 0, "")

	pdf.SetX(58)
	pdf.SetY(93)
	pdf.CellFormat(155, 8, order.Product, "", 0, "L", false, 0, "")

	pdf.SetX(166)
	pdf.CellFormat(20, 8, fmt.Sprint(order.Quantity), "", 0, "C", false, 0, "")

	pdf.SetX(185)
	pdf.CellFormat(20, 8, fmt.Sprintf("$%.2f", float32(order.Amount/100)), "", 0, "R", false, 0, "")

	invoicePath := fmt.Sprintf("./pdf/%d.pdf", order.ID)
	err := pdf.OutputFileAndClose(invoicePath)
	if err != nil {
		return err
	}

	return nil
}
