package main

import (
	"fmt"
	"net/http"

	"github.com/Nico2220/billease/internal/data"
	"github.com/Nico2220/billease/internal/pdf"
	"github.com/Nico2220/billease/internal/validator"
)

func (app *application) createInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		From     int64          `json:"from"`
		To       int64          `json:"to"`
		Services []data.Service `json:"services"`
		Vat      int            `json:"vat"`
		Currency string         `json:"currency"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, responseFormat{"error": err.Error()}, nil)
		return
	}

	invoice := data.Invoice{
		From:     input.From,
		To:       input.To,
		Services: input.Services,
		Vat:      input.Vat,
		Currency: input.Currency,
	}

	invoice.CalculateSubTotal()
	invoice.CalculateTax()
	invoice.CalculateTotal()

	v := validator.New()

	if data.ValidateInvoice(v, &invoice); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// insert invoice into db
	_, err = app.models.Invoices.Insert(invoice)

	if err != nil {
		fmt.Println("err")
		return
	}

	c1 := data.GetCompany(invoice.From)
	c2 := data.GetCompany(invoice.To)

	//create invoice
	filename := "test.pdf"
	err = pdf.New(invoice, c1, c2, filename)
	if err != nil {
		fmt.Println("err")
		return
	}

	invoice.Link = fmt.Sprintf("https://%s.s3.eu-west-1.amazonaws.com/%s", app.config.aws.bucketName, filename)

	app.awsService.UploadFile(filename)

	app.writeJSON(w, http.StatusOK, responseFormat{"invoice": invoice}, nil)

}
