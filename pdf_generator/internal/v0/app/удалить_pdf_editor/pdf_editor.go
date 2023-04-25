package delete_pdf_editor

import (
	// "github.com/manyakrus/starter/common/v0/contextmain"
	// stopapp "github.com/manyakrus/starter/common/v0/stopapp"

	pdf "github.com/adrg/go-wkhtmltopdf"
	logger "github.com/manyakrus/starter/common/v0/logger"
	"os"
)

//// log - глобальный логгер
var log = logger.GetLog()

func init() {

}

func CreatePDF_from_XLSX(filename, filename2 string, map1 map[string]interface{}) error {
	var err error

	//// Initialize library.
	//if err := pdf.Init(); err != nil {
	//	log.Fatal(err)
	//}
	//defer pdf.Destroy()
	//
	//// Create object from file.
	//object, err := pdf.NewObject(filename)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//object.Header.ContentCenter = "[title]"
	//object.Header.DisplaySeparator = true
	//
	////// Create object from URL.
	////object2, err := pdf.NewObject("https://google.com")
	////if err != nil {
	////	log.Fatal(err)
	////}
	//object.Footer.ContentLeft = "[date]"
	//object.Footer.ContentCenter = "Sample footer information"
	//object.Footer.ContentRight = "[page]"
	//object.Footer.DisplaySeparator = true
	//
	////// Create object from reader.
	////inFile, err := os.Open(filename)
	////if err != nil {
	////	log.Fatal(err)
	////}
	////defer inFile.Close()
	////
	////object3, err := pdf.NewObjectFromReader(inFile)
	////if err != nil {
	////	log.Fatal(err)
	////}
	////object3.Zoom = 1.5
	////object3.TOC.Title = "Table of Contents"
	//
	//// Create converter.
	//converter, err := pdf.NewConverter()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer converter.Destroy()
	//
	//// Add created objects to the converter.
	//converter.Add(object)
	////converter.Add(object2)
	////converter.Add(object3)
	//
	//// Set converter options.
	//converter.Title = "Sample document"
	//converter.PaperSize = pdf.A4
	//converter.Orientation = pdf.Landscape
	//converter.MarginTop = "1cm"
	//converter.MarginBottom = "1cm"
	//converter.MarginLeft = "10mm"
	//converter.MarginRight = "10mm"
	//
	//// Convert objects and save the output PDF document.
	//outFile, err := os.Create(filename2)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer outFile.Close()
	//
	//if err := converter.Run(outFile); err != nil {
	//	log.Fatal(err)
	//}

	//file, err := pdf_io.Open(filename)
	//if err != nil {
	//	panic(err)
	//}
	//page1 := file.Page(1)
	//fmt.Println(page1.Content())

	//
	//
	//var pdf1 pdft.PDFt
	//err = pdf1.Open(filename)
	//if err != nil {
	//	log.Panic("Couldn't open pdf filename: ", filename)
	//}
	//
	//file, err := os.Open(filename)
	//var outPDF *pdft.PDFData
	//outPDF = &pdft.PDFData{}
	//pdft.PDFParse(file, outPDF)
	//
	////var MassObj []int
	//MassObj := make([]int, 0)
	////MassObj, err = outPDF.GetPageObjIDs()
	//MassObj, err = pdf1.PDFdata.GetPageObjIDs()
	//for _, id1 := range MassObj {
	//	obj1 := pdf1.PDFdata.GetObjByID(id1)
	//	if obj1 == nil {
	//		continue
	//	}
	//	PropertiesData, err := obj1.ReadProperties()
	//	if err != nil {
	//		continue
	//	}
	//	fmt.Printf("PropertiesData: %#v \n", PropertiesData)
	//}
	//
	//dir1 := filepath.Dir(filename)
	//filename2 := dir1 + micro.SeparatorFile() + "test_ready.pdf"
	//
	//err = pdf1.Save(filename2)
	//if err != nil {
	//	log.Panic("Couldn't save pdf filename: ", filename2)
	//}
	return err

}

func CreatePDF1(filename, filename2 string, map1 map[string]interface{}) error {
	var err error

	// Initialize library.
	//if err := pdf.Init(); err != nil {
	//	log.Fatal(err)
	//}
	defer pdf.Destroy()

	// Create object from file.
	object, err := pdf.NewObject(filename)
	if err != nil {
		log.Fatal(err)
	}
	object.Header.ContentCenter = "[title]"
	object.Header.DisplaySeparator = true

	//// Create object from URL.
	//object2, err := pdf.NewObject("https://google.com")
	//if err != nil {
	//	log.Fatal(err)
	//}
	object.Footer.ContentLeft = "[date]"
	object.Footer.ContentCenter = "Sample footer information"
	object.Footer.ContentRight = "[page]"
	object.Footer.DisplaySeparator = true

	//// Create object from reader.
	//inFile, err := os.Open(filename)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer inFile.Close()
	//
	//object3, err := pdf.NewObjectFromReader(inFile)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//object3.Zoom = 1.5
	//object3.TOC.Title = "Table of Contents"

	// Create converter.
	converter, err := pdf.NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add created objects to the converter.
	converter.Add(object)
	//converter.Add(object2)
	//converter.Add(object3)

	// Set converter options.
	converter.Title = "Sample document"
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Landscape
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "10mm"
	converter.MarginRight = "10mm"

	// Convert objects and save the output PDF document.
	outFile, err := os.Create(filename2)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err := converter.Run(outFile); err != nil {
		log.Fatal(err)
	}

	//file, err := pdf_io.Open(filename)
	//if err != nil {
	//	panic(err)
	//}
	//page1 := file.Page(1)
	//fmt.Println(page1.Content())

	//
	//
	//var pdf1 pdft.PDFt
	//err = pdf1.Open(filename)
	//if err != nil {
	//	log.Panic("Couldn't open pdf filename: ", filename)
	//}
	//
	//file, err := os.Open(filename)
	//var outPDF *pdft.PDFData
	//outPDF = &pdft.PDFData{}
	//pdft.PDFParse(file, outPDF)
	//
	////var MassObj []int
	//MassObj := make([]int, 0)
	////MassObj, err = outPDF.GetPageObjIDs()
	//MassObj, err = pdf1.PDFdata.GetPageObjIDs()
	//for _, id1 := range MassObj {
	//	obj1 := pdf1.PDFdata.GetObjByID(id1)
	//	if obj1 == nil {
	//		continue
	//	}
	//	PropertiesData, err := obj1.ReadProperties()
	//	if err != nil {
	//		continue
	//	}
	//	fmt.Printf("PropertiesData: %#v \n", PropertiesData)
	//}
	//
	//dir1 := filepath.Dir(filename)
	//filename2 := dir1 + micro.SeparatorFile() + "test_ready.pdf"
	//
	//err = pdf1.Save(filename2)
	//if err != nil {
	//	log.Panic("Couldn't save pdf filename: ", filename2)
	//}
	return err
}
