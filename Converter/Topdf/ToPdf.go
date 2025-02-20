package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/phpdave11/gofpdf"
	"github.com/unidoc/unioffice/document"
)

func main() {
	// Define the input file and its format
	inputFile := "D:/OneDrive - VietNam National University - HCM INTERNATIONAL UNIVERSITY/aoisjdaojiasodj.txt" // Change this to the desired input file
	format := strings.ToLower(strings.Split(inputFile, ".")[1])

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	switch format {
	case "txt":
		convertTxtToPdf(inputFile, pdf)
	case "csv":
		convertCsvToPdf(inputFile, pdf)
	case "xlsx":
		convertXlsxToPdf(inputFile, pdf)
	case "doc", "docx":
		convertDocToPdf(inputFile, pdf)
	default:
		fmt.Println("Unsupported file format:", format)
		return
	}

	// Save the PDF to a file
	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		fmt.Println("Error saving PDF file:", err)
		return
	}

	fmt.Println("PDF file created successfully.")
}

func convertTxtToPdf(inputFile string, pdf *gofpdf.Fpdf) {
	txtFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening text file:", err)
		return
	}
	defer txtFile.Close()

	scanner := bufio.NewScanner(txtFile)
	for scanner.Scan() {
		line := scanner.Text()
		pdf.Cell(0, 10, line)
		pdf.Ln(10)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading text file:", err)
		return
	}
}

func convertCsvToPdf(inputFile string, pdf *gofpdf.Fpdf) {
	csvFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer csvFile.Close()

	scanner := bufio.NewScanner(csvFile)
	for scanner.Scan() {
		line := scanner.Text()
		pdf.Cell(0, 10, line)
		pdf.Ln(10)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}
}

func convertXlsxToPdf(inputFile string, pdf *gofpdf.Fpdf) {
	xlFile, err := excelize.OpenFile(inputFile)
	if err != nil {
		fmt.Println("Error opening XLSX file:", err)
		return
	}

	for _, sheetName := range xlFile.GetSheetMap() {
		err := xlFile.GetRows(sheetName)
		rows := xlFile.GetRows(sheetName)
		if err != nil {
			fmt.Println("Error reading sheet:", err)
			return
		}

		for _, row := range rows {
			for _, cell := range row {
				pdf.Cell(40, 10, cell)
			}
			pdf.Ln(10)
		}
	}
}

func convertDocToPdf(inputFile string, pdf *gofpdf.Fpdf) {
	doc, err := document.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening DOC/DOCX file:", err)
		return
	}

	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			pdf.Cell(0, 10, run.Text())
		}
		pdf.Ln(10)
	}
}
