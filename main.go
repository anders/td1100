package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/barcode"
)

type barcodeList struct {
	Code        string
	Width       float64
	Description string
}

var defaultList = []barcodeList{
	{"[FNC3]$P,Ae,P\r", 70, "Återställ till EU-fabriksinställningar"},
	{"[FNC3]$P,HA35,P\r", 80, "Simulera USB-tangentbord (tryck)"},
	{"[FNC3]$P\r", 50, "Programmeringsläge (tryck)"},
	{"[FNC3]$CUSSE01\r", 70, "Aktivera USB-viloläge (tryck)"},
	{"[FNC3]$CSNRM04\r", 70, "Konstant kodläsning (tryck)"},
	{"[FNC3]$CPLEN01\r", 70, "Aktivera UK Plessey-läsning (tryck)"},
	{"[FNC3]$CPLCC01\r", 70, "Aktivera UK Plessey-kontrollsifferuträkning (tryck)"},
	{"[FNC3]$CPLCT01\r", 70, "Inkludera UK Plessey-kontrollsiffror i läsning (tryck)"},
	{"[FNC3]$C3BEN00\r", 70, "Deaktivera EAN13-läsning (ISBN, ISSN) (tryck)"},
	{"[FNC3]$C8BEN00\r", 70, "Deaktivera EAN8-läsning (ISBN, ISSN) (tryck)"},
	{"[FNC3]$P\r", 50, "Avsluta programmeringsläge (tryck)"},
}

var (
	writeList  bool
	listPath   string
	outputPath string
)

func init() {
	flag.BoolVar(&writeList, "writeList", false, "write default list to file")
	flag.StringVar(&listPath, "list", "", "path to JSON list of barcodes (optional)")
	flag.StringVar(&outputPath, "output", "", "path to PDF output")
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	flag.Parse()

	if writeList {
		if listPath == "" {
			flag.Usage()
			return
		}
		b, err := json.MarshalIndent(defaultList, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		b = append(b, '\n')
		if err := ioutil.WriteFile(listPath, b, 0644); err != nil {
			log.Fatal(err)
		}
		return
	}

	if outputPath == "" {
		flag.Usage()
		return
	}

	var barcodes []barcodeList
	if listPath != "" {
		b, err := ioutil.ReadFile(listPath)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(b, &barcodes); err != nil {
			log.Fatal(err)
		}
	} else {
		barcodes = defaultList
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetFont("Helvetica", "", 12)

	u := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(0, 30,
			u("Programmering av TD1100-streckkodsläsare"),
			// border, ln
			"0", 0,
			// align, fill,
			"L", false,
			//link, link str
			0, "")
		pdf.Line(10, 30, 210-30, 30)

		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(0, 10, u(fmt.Sprintf("sida %d av {nb}", pdf.PageNo())), "0", 0, "R", false, 0, "")

		pdf.Ln(30)
	})

	pdf.AliasNbPages("")
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 12)
	pdf.Write(5, u("Läs streckkoderna i tur och ordning. Vänta ett par sekunder mellan varje kodläsning."))
	pdf.Ln(20)

	drawBC := func(code string, w, h float64) {
		code = strings.Replace(code, "[FNC3]", string(code128.FNC3), 1)
		bc, err := code128.Encode(code)
		if err != nil {
			log.Fatal(err)
		}
		key := barcode.Register(bc)
		barcode.Barcode(pdf, key, 10+pdf.GetX(), pdf.GetY(), w, h-10, false)
		pdf.Ln(h + 5)
	}

	for i, bc := range barcodes {
		pdf.Cell(100, 0, u(fmt.Sprintf("%d. %s", i+1, bc.Description)))
		pdf.Ln(5)
		drawBC(bc.Code, bc.Width, 30)
	}

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		log.Fatal(err)
	}
}
