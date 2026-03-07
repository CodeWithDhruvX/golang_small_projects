package main

import "fmt"

// Builder Pattern

// Builder interface
type Builder interface {
	SetPart1(string)
	SetPart2(string)
	SetPart3(string)
	GetResult() Product
}

// Product represents the complex object being built
type Product struct {
	part1 string
	part2 string
	part3 string
}

func (p *Product) Show() {
	fmt.Printf("Product with parts: %s, %s, %s\n", p.part1, p.part2, p.part3)
}

// ConcreteBuilder implements the Builder interface
type ConcreteBuilder struct {
	product *Product
}

func NewConcreteBuilder() *ConcreteBuilder {
	return &ConcreteBuilder{
		product: &Product{},
	}
}

func (cb *ConcreteBuilder) SetPart1(part string) {
	cb.product.part1 = part
	fmt.Printf("Part1 set to: %s\n", part)
}

func (cb *ConcreteBuilder) SetPart2(part string) {
	cb.product.part2 = part
	fmt.Printf("Part2 set to: %s\n", part)
}

func (cb *ConcreteBuilder) SetPart3(part string) {
	cb.product.part3 = part
	fmt.Printf("Part3 set to: %s\n", part)
}

func (cb *ConcreteBuilder) GetResult() Product {
	return *cb.product
}

// Director constructs the object using the builder interface
type Director struct {
	builder Builder
}

func NewDirector(builder Builder) *Director {
	return &Director{builder: builder}
}

func (d *Director) SetBuilder(builder Builder) {
	d.builder = builder
}

func (d *Director) ConstructMinimalProduct() {
	fmt.Println("Building minimal product:")
	d.builder.SetPart1("Minimal Part1")
	d.builder.SetPart2("Minimal Part2")
}

func (d *Director) ConstructFullProduct() {
	fmt.Println("Building full product:")
	d.builder.SetPart1("Full Part1")
	d.builder.SetPart2("Full Part2")
	d.builder.SetPart3("Full Part3")
}

func (d *Director) ConstructCustomProduct(part1, part2, part3 string) {
	fmt.Println("Building custom product:")
	d.builder.SetPart1(part1)
	d.builder.SetPart2(part2)
	d.builder.SetPart3(part3)
}

// House Building Example
type House struct {
	foundation string
	walls      string
	roof       string
	doors      string
	windows    string
	garage     string
	garden     string
}

func (h *House) Display() {
	fmt.Printf("House Details:\n")
	fmt.Printf("  Foundation: %s\n", h.foundation)
	fmt.Printf("  Walls: %s\n", h.walls)
	fmt.Printf("  Roof: %s\n", h.roof)
	fmt.Printf("  Doors: %s\n", h.doors)
	fmt.Printf("  Windows: %s\n", h.windows)
	fmt.Printf("  Garage: %s\n", h.garage)
	fmt.Printf("  Garden: %s\n", h.garden)
}

type HouseBuilder interface {
	SetFoundation(string)
	SetWalls(string)
	SetRoof(string)
	SetDoors(string)
	SetWindows(string)
	SetGarage(string)
	SetGarden(string)
	GetHouse() House
}

type ConcreteHouseBuilder struct {
	house *House
}

func NewConcreteHouseBuilder() *ConcreteHouseBuilder {
	return &ConcreteHouseBuilder{
		house: &House{},
	}
}

func (chb *ConcreteHouseBuilder) SetFoundation(foundation string) {
	chb.house.foundation = foundation
}

func (chb *ConcreteHouseBuilder) SetWalls(walls string) {
	chb.house.walls = walls
}

func (chb *ConcreteHouseBuilder) SetRoof(roof string) {
	chb.house.roof = roof
}

func (chb *ConcreteHouseBuilder) SetDoors(doors string) {
	chb.house.doors = doors
}

func (chb *ConcreteHouseBuilder) SetWindows(windows string) {
	chb.house.windows = windows
}

func (chb *ConcreteHouseBuilder) SetGarage(garage string) {
	chb.house.garage = garage
}

func (chb *ConcreteHouseBuilder) SetGarden(garden string) {
	chb.house.garden = garden
}

func (chb *ConcreteHouseBuilder) GetHouse() House {
	return *chb.house
}

type HouseDirector struct {
	builder HouseBuilder
}

func NewHouseDirector(builder HouseBuilder) *HouseDirector {
	return &HouseDirector{builder: builder}
}

func (hd *HouseDirector) SetBuilder(builder HouseBuilder) {
	hd.builder = builder
}

func (hd *HouseDirector) BuildBasicHouse() {
	hd.builder.SetFoundation("Concrete Foundation")
	hd.builder.SetWalls("Brick Walls")
	hd.builder.SetRoof("Tile Roof")
	hd.builder.SetDoors("Wooden Doors")
	hd.builder.SetWindows("Standard Windows")
}

func (hd *HouseDirector) BuildLuxuryHouse() {
	hd.builder.SetFoundation("Reinforced Concrete Foundation")
	hd.builder.SetWalls("Stone Walls with Insulation")
	hd.builder.SetRoof("Slate Roof with Solar Panels")
	hd.builder.SetDoors("Mahogany Doors")
	hd.builder.SetWindows("Double-glazed Windows")
	hd.builder.SetGarage("Two-car Garage")
	hd.builder.SetGarden("Landscaped Garden")
}

func (hd *HouseDirector) BuildApartment() {
	hd.builder.SetFoundation("Steel Frame Foundation")
	hd.builder.SetWalls("Drywall Walls")
	hd.builder.SetRoof("Flat Roof")
	hd.builder.SetDoors("Steel Doors")
	hd.builder.SetWindows("Large Windows")
}

// Car Building Example
type Car struct {
	engine    string
	body      string
	wheels    string
	interior  string
	paint     string
	electronics string
}

func (c *Car) Show() {
	fmt.Printf("Car Specifications:\n")
	fmt.Printf("  Engine: %s\n", c.engine)
	fmt.Printf("  Body: %s\n", c.body)
	fmt.Printf("  Wheels: %s\n", c.wheels)
	fmt.Printf("  Interior: %s\n", c.interior)
	fmt.Printf("  Paint: %s\n", c.paint)
	fmt.Printf("  Electronics: %s\n", c.electronics)
}

type CarBuilder interface {
	SetEngine(string)
	SetBody(string)
	SetWheels(string)
	SetInterior(string)
	SetPaint(string)
	SetElectronics(string)
	GetCar() Car
}

type SportsCarBuilder struct {
	car *Car
}

func NewSportsCarBuilder() *SportsCarBuilder {
	return &SportsCarBuilder{
		car: &Car{},
	}
}

func (scb *SportsCarBuilder) SetEngine(engine string) {
	scb.car.engine = engine
}

func (scb *SportsCarBuilder) SetBody(body string) {
	scb.car.body = body
}

func (scb *SportsCarBuilder) SetWheels(wheels string) {
	scb.car.wheels = wheels
}

func (scb *SportsCarBuilder) SetInterior(interior string) {
	scb.car.interior = interior
}

func (scb *SportsCarBuilder) SetPaint(paint string) {
	scb.car.paint = paint
}

func (scb *SportsCarBuilder) SetElectronics(electronics string) {
	scb.car.electronics = electronics
}

func (scb *SportsCarBuilder) GetCar() Car {
	return *scb.car
}

type SUVBuilder struct {
	car *Car
}

func NewSUVBuilder() *SUVBuilder {
	return &SUVBuilder{
		car: &Car{},
	}
}

func (sub *SUVBuilder) SetEngine(engine string) {
	sub.car.engine = engine
}

func (sub *SUVBuilder) SetBody(body string) {
	sub.car.body = body
}

func (sub *SUVBuilder) SetWheels(wheels string) {
	sub.car.wheels = wheels
}

func (sub *SUVBuilder) SetInterior(interior string) {
	sub.car.interior = interior
}

func (sub *SUVBuilder) SetPaint(paint string) {
	sub.car.paint = paint
}

func (sub *SUVBuilder) SetElectronics(electronics string) {
	sub.car.electronics = electronics
}

func (sub *SUVBuilder) GetCar() Car {
	return *sub.car
}

type CarDirector struct {
	builder CarBuilder
}

func NewCarDirector(builder CarBuilder) *CarDirector {
	return &CarDirector{builder: builder}
}

func (cd *CarDirector) SetBuilder(builder CarBuilder) {
	cd.builder = builder
}

func (cd *CarDirector) BuildSportsCar() {
	cd.builder.SetEngine("V8 Turbo Engine")
	cd.builder.SetBody("Carbon Fiber Body")
	cd.builder.SetWheels("Performance Wheels")
	cd.builder.SetInterior("Leather Racing Seats")
	cd.builder.SetPaint("Red Metallic Paint")
	cd.builder.SetElectronics("Advanced Navigation System")
}

func (cd *CarDirector) BuildSUV() {
	cd.builder.SetEngine("V6 Engine")
	cd.builder.SetBody("Steel Frame Body")
	cd.builder.SetWheels("All-terrain Wheels")
	cd.builder.SetInterior("Fabric Seats")
	cd.builder.SetPaint("Silver Paint")
	cd.builder.SetElectronics("Basic GPS System")
}

// Report Builder Example
type Report struct {
	title    string
	header   string
	content  string
	footer   string
	format   string
}

func (r *Report) Display() {
	fmt.Printf("Report (%s format):\n", r.format)
	fmt.Printf("Title: %s\n", r.title)
	fmt.Printf("Header: %s\n", r.header)
	fmt.Printf("Content: %s\n", r.content)
	fmt.Printf("Footer: %s\n", r.footer)
}

type ReportBuilder interface {
	SetTitle(string)
	SetHeader(string)
	SetContent(string)
	SetFooter(string)
	SetFormat(string)
	GetReport() Report
}

type PDFReportBuilder struct {
	report *Report
}

func NewPDFReportBuilder() *PDFReportBuilder {
	return &PDFReportBuilder{
		report: &Report{format: "PDF"},
	}
}

func (prb *PDFReportBuilder) SetTitle(title string) {
	prb.report.title = title + " [PDF]"
}

func (prb *PDFReportBuilder) SetHeader(header string) {
	prb.report.header = header + " [PDF Header]"
}

func (prb *PDFReportBuilder) SetContent(content string) {
	prb.report.content = content + " [PDF Content]"
}

func (prb *PDFReportBuilder) SetFooter(footer string) {
	prb.report.footer = footer + " [PDF Footer]"
}

func (prb *PDFReportBuilder) SetFormat(format string) {
	prb.report.format = format
}

func (prb *PDFReportBuilder) GetReport() Report {
	return *prb.report
}

type HTMLReportBuilder struct {
	report *Report
}

func NewHTMLReportBuilder() *HTMLReportBuilder {
	return &HTMLReportBuilder{
		report: &Report{format: "HTML"},
	}
}

func (hrb *HTMLReportBuilder) SetTitle(title string) {
	hrb.report.title = "<h1>" + title + "</h1>"
}

func (hrb *HTMLReportBuilder) SetHeader(header string) {
	hrb.report.header = "<header>" + header + "</header>"
}

func (hrb *HTMLReportBuilder) SetContent(content string) {
	hrb.report.content = "<p>" + content + "</p>"
}

func (hrb *HTMLReportBuilder) SetFooter(footer string) {
	hrb.report.footer = "<footer>" + footer + "</footer>"
}

func (hrb *HTMLReportBuilder) SetFormat(format string) {
	hrb.report.format = format
}

func (hrb *HTMLReportBuilder) GetReport() Report {
	return *hrb.report
}

func main() {
	fmt.Println("=== Builder Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Builder Example ---")
	builder := NewConcreteBuilder()
	director := NewDirector(builder)
	
	director.ConstructMinimalProduct()
	minimalProduct := builder.GetResult()
	minimalProduct.Show()
	
	director.ConstructFullProduct()
	fullProduct := builder.GetResult()
	fullProduct.Show()
	
	// House Building example
	fmt.Println("\n--- House Building Example ---")
	houseBuilder := NewConcreteHouseBuilder()
	houseDirector := NewHouseDirector(houseBuilder)
	
	fmt.Println("\nBuilding Basic House:")
	houseDirector.BuildBasicHouse()
	basicHouse := houseBuilder.GetHouse()
	basicHouse.Display()
	
	fmt.Println("\nBuilding Luxury House:")
	houseDirector.BuildLuxuryHouse()
	luxuryHouse := houseBuilder.GetHouse()
	luxuryHouse.Display()
	
	// Car Building example
	fmt.Println("\n--- Car Building Example ---")
	sportsCarBuilder := NewSportsCarBuilder()
	carDirector := NewCarDirector(sportsCarBuilder)
	
	fmt.Println("\nBuilding Sports Car:")
	carDirector.BuildSportsCar()
	sportsCar := sportsCarBuilder.GetCar()
	sportsCar.Show()
	
	suvBuilder := NewSUVBuilder()
	carDirector.SetBuilder(suvBuilder)
	
	fmt.Println("\nBuilding SUV:")
	carDirector.BuildSUV()
	suv := suvBuilder.GetCar()
	suv.Show()
	
	// Report Builder example
	fmt.Println("\n--- Report Builder Example ---")
	
	fmt.Println("\nCreating PDF Report:")
	pdfBuilder := NewPDFReportBuilder()
	pdfBuilder.SetTitle("Annual Report")
	pdfBuilder.SetHeader("Company Financials")
	pdfBuilder.SetContent("This is the financial report content...")
	pdfBuilder.SetFooter("Confidential Document")
	pdfReport := pdfBuilder.GetReport()
	pdfReport.Display()
	
	fmt.Println("\nCreating HTML Report:")
	htmlBuilder := NewHTMLReportBuilder()
	htmlBuilder.SetTitle("Annual Report")
	htmlBuilder.SetHeader("Company Financials")
	htmlBuilder.SetContent("This is the financial report content...")
	htmlBuilder.SetFooter("Confidential Document")
	htmlReport := htmlBuilder.GetReport()
	htmlReport.Display()
	
	fmt.Println("\nAll builder patterns demonstrated successfully!")
}
