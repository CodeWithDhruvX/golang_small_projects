package main

import "fmt"

// Factory Method Pattern

// Product interface
type Product interface {
	Use()
}

// ConcreteProductA
type ConcreteProductA struct{}

func (cpa *ConcreteProductA) Use() {
	fmt.Println("ConcreteProductA: Used")
}

// ConcreteProductB
type ConcreteProductB struct{}

func (cpb *ConcreteProductB) Use() {
	fmt.Println("ConcreteProductB: Used")
}

// Creator interface (Factory Method)
type Creator interface {
	FactoryMethod() Product
	SomeOperation()
}

// ConcreteCreatorA
type ConcreteCreatorA struct{}

func (cca *ConcreteCreatorA) FactoryMethod() Product {
	return &ConcreteProductA{}
}

func (cca *ConcreteCreatorA) SomeOperation() {
	product := cca.FactoryMethod()
	fmt.Println("ConcreteCreatorA: Working with", product)
	product.Use()
}

// ConcreteCreatorB
type ConcreteCreatorB struct{}

func (ccb *ConcreteCreatorB) FactoryMethod() Product {
	return &ConcreteProductB{}
}

func (ccb *ConcreteCreatorB) SomeOperation() {
	product := ccb.FactoryMethod()
	fmt.Println("ConcreteCreatorB: Working with", product)
	product.Use()
}

// Client code
func clientCode(creator Creator) {
	creator.SomeOperation()
}

// Animal Factory Example
type Animal interface {
	Sound()
}

type Dog struct{}

func (d *Dog) Sound() {
	fmt.Println("Dog: Woof!")
}

type Cat struct{}

func (c *Cat) Sound() {
	fmt.Println("Cat: Meow!")
}

type Cow struct{}

func (c *Cow) Sound() {
	fmt.Println("Cow: Moo!")
}

// Animal Factory
type AnimalFactory interface {
	CreateAnimal() Animal
}

type DogFactory struct{}

func (df *DogFactory) CreateAnimal() Animal {
	return &Dog{}
}

type CatFactory struct{}

func (cf *CatFactory) CreateAnimal() Animal {
	return &Cat{}
}

type CowFactory struct{}

func (cf *CowFactory) CreateAnimal() Animal {
	return &Cow{}
}

func playWithAnimal(factory AnimalFactory) {
	animal := factory.CreateAnimal()
	fmt.Printf("Playing with %T: ", animal)
	animal.Sound()
}

// Vehicle Factory Example
type Vehicle interface {
	Drive()
	Stop()
}

type Car struct{}

func (c *Car) Drive() {
	fmt.Println("Car: Driving smoothly")
}

func (c *Car) Stop() {
	fmt.Println("Car: Stopping at traffic light")
}

type Motorcycle struct{}

func (m *Motorcycle) Drive() {
	fmt.Println("Motorcycle: Riding fast")
}

func (m *Motorcycle) Stop() {
	fmt.Println("Motorcycle: Stopping quickly")
}

type Truck struct{}

func (t *Truck) Drive() {
	fmt.Println("Truck: Driving slowly with heavy load")
}

func (t *Truck) Stop() {
	fmt.Println("Truck: Stopping with air brakes")
}

// Vehicle Factory
type VehicleFactory interface {
	CreateVehicle() Vehicle
}

type CarFactory struct{}

func (cf *CarFactory) CreateVehicle() Vehicle {
	return &Car{}
}

type MotorcycleFactory struct{}

func (mf *MotorcycleFactory) CreateVehicle() Vehicle {
	return &Motorcycle{}
}

type TruckFactory struct{}

func (tf *TruckFactory) CreateVehicle() Vehicle {
	return &Truck{}
}

func testDrive(factory VehicleFactory) {
	vehicle := factory.CreateVehicle()
	fmt.Printf("Test driving %T:\n", vehicle)
	vehicle.Drive()
	vehicle.Stop()
}

// Document Factory Example
type Document interface {
	Open()
	Save()
	Print()
}

type PDFDocument struct{}

func (pdf *PDFDocument) Open() {
	fmt.Println("PDF Document: Opening with PDF reader")
}

func (pdf *PDFDocument) Save() {
	fmt.Println("PDF Document: Saving as PDF format")
}

func (pdf *PDFDocument) Print() {
	fmt.Println("PDF Document: Printing to PDF printer")
}

type WordDocument struct{}

func (wd *WordDocument) Open() {
	fmt.Println("Word Document: Opening with Microsoft Word")
}

func (wd *WordDocument) Save() {
	fmt.Println("Word Document: Saving as DOCX format")
}

func (wd *WordDocument) Print() {
	fmt.Println("Word Document: Printing to default printer")
}

type ExcelDocument struct{}

func (ed *ExcelDocument) Open() {
	fmt.Println("Excel Document: Opening with Microsoft Excel")
}

func (ed *ExcelDocument) Save() {
	fmt.Println("Excel Document: Saving as XLSX format")
}

func (ed *ExcelDocument) Print() {
	fmt.Println("Excel Document: Printing to spreadsheet printer")
}

// Document Factory
type DocumentFactory interface {
	CreateDocument() Document
}

type PDFDocumentFactory struct{}

func (pdf *PDFDocumentFactory) CreateDocument() Document {
	return &PDFDocument{}
}

type WordDocumentFactory struct{}

func (wd *WordDocumentFactory) CreateDocument() Document {
	return &WordDocument{}
}

type ExcelDocumentFactory struct{}

func (ed *ExcelDocumentFactory) CreateDocument() Document {
	return &ExcelDocument{}
}

func workWithDocument(factory DocumentFactory) {
	document := factory.CreateDocument()
	fmt.Printf("Working with %T:\n", document)
	document.Open()
	document.Save()
	document.Print()
}

// Shape Factory Example
type Shape interface {
	Draw()
	GetArea() float64
}

type Circle struct {
	radius float64
}

func NewCircle(radius float64) *Circle {
	return &Circle{radius: radius}
}

func (c *Circle) Draw() {
	fmt.Printf("Drawing Circle with radius %.2f\n", c.radius)
}

func (c *Circle) GetArea() float64 {
	return 3.14159 * c.radius * c.radius
}

type Rectangle struct {
	width  float64
	height float64
}

func NewRectangle(width, height float64) *Rectangle {
	return &Rectangle{width: width, height: height}
}

func (r *Rectangle) Draw() {
	fmt.Printf("Drawing Rectangle %.2f x %.2f\n", r.width, r.height)
}

func (r *Rectangle) GetArea() float64 {
	return r.width * r.height
}

type Triangle struct {
	base   float64
	height float64
}

func NewTriangle(base, height float64) *Triangle {
	return &Triangle{base: base, height: height}
}

func (t *Triangle) Draw() {
	fmt.Printf("Drawing Triangle with base %.2f and height %.2f\n", t.base, t.height)
}

func (t *Triangle) GetArea() float64 {
	return 0.5 * t.base * t.height
}

// Shape Factory
type ShapeFactory interface {
	CreateShape() Shape
}

type CircleFactory struct {
	radius float64
}

func NewCircleFactory(radius float64) *CircleFactory {
	return &CircleFactory{radius: radius}
}

func (cf *CircleFactory) CreateShape() Shape {
	return NewCircle(cf.radius)
}

type RectangleFactory struct {
	width  float64
	height float64
}

func NewRectangleFactory(width, height float64) *RectangleFactory {
	return &RectangleFactory{width: width, height: height}
}

func (rf *RectangleFactory) CreateShape() Shape {
	return NewRectangle(rf.width, rf.height)
}

type TriangleFactory struct {
	base   float64
	height float64
}

func NewTriangleFactory(base, height float64) *TriangleFactory {
	return &TriangleFactory{base: base, height: height}
}

func (tf *TriangleFactory) CreateShape() Shape {
	return NewTriangle(tf.base, tf.height)
}

func drawShape(factory ShapeFactory) {
	shape := factory.CreateShape()
	shape.Draw()
	fmt.Printf("Area: %.2f\n", shape.GetArea())
}

// Payment Method Factory Example
type PaymentMethod interface {
	Pay(amount float64)
	Validate() bool
}

type CreditCard struct {
	cardNumber string
	name       string
}

func (cc *CreditCard) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Credit Card ending in %s\n", 
		amount, cc.cardNumber[len(cc.cardNumber)-4:])
}

func (cc *CreditCard) Validate() bool {
	fmt.Println("Validating Credit Card...")
	return len(cc.cardNumber) == 16
}

type PayPal struct {
	email string
}

func (pp *PayPal) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using PayPal (%s)\n", amount, pp.email)
}

func (pp *PayPal) Validate() bool {
	fmt.Println("Validating PayPal account...")
	return len(pp.email) > 0
}

type BankTransfer struct {
	accountNumber string
	bankName      string
}

func (bt *BankTransfer) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Bank Transfer (%s)\n", amount, bt.bankName)
}

func (bt *BankTransfer) Validate() bool {
	fmt.Println("Validating bank account...")
	return len(bt.accountNumber) > 0
}

// Payment Factory
type PaymentFactory interface {
	CreatePaymentMethod() PaymentMethod
}

type CreditCardFactory struct {
	cardNumber string
	name       string
}

func NewCreditCardFactory(cardNumber, name string) *CreditCardFactory {
	return &CreditCardFactory{cardNumber: cardNumber, name: name}
}

func (ccf *CreditCardFactory) CreatePaymentMethod() PaymentMethod {
	return &CreditCard{cardNumber: ccf.cardNumber, name: ccf.name}
}

type PayPalFactory struct {
	email string
}

func NewPayPalFactory(email string) *PayPalFactory {
	return &PayPalFactory{email: email}
}

func (ppf *PayPalFactory) CreatePaymentMethod() PaymentMethod {
	return &PayPal{email: ppf.email}
}

type BankTransferFactory struct {
	accountNumber string
	bankName      string
}

func NewBankTransferFactory(accountNumber, bankName string) *BankTransferFactory {
	return &BankTransferFactory{accountNumber: accountNumber, bankName: bankName}
}

func (btf *BankTransferFactory) CreatePaymentMethod() PaymentMethod {
	return &BankTransfer{accountNumber: btf.accountNumber, bankName: btf.bankName}
}

func processPayment(factory PaymentFactory, amount float64) {
	payment := factory.CreatePaymentMethod()
	if payment.Validate() {
		payment.Pay(amount)
	} else {
		fmt.Println("Payment validation failed!")
	}
}

func main() {
	fmt.Println("=== Factory Method Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Factory Method Example ---")
	fmt.Println("Applying ConcreteCreatorA:")
	clientCode(&ConcreteCreatorA{})
	
	fmt.Println("\nApplying ConcreteCreatorB:")
	clientCode(&ConcreteCreatorB{})
	
	// Animal Factory example
	fmt.Println("\n--- Animal Factory Example ---")
	playWithAnimal(&DogFactory{})
	playWithAnimal(&CatFactory{})
	playWithAnimal(&CowFactory{})
	
	// Vehicle Factory example
	fmt.Println("\n--- Vehicle Factory Example ---")
	testDrive(&CarFactory{})
	testDrive(&MotorcycleFactory{})
	testDrive(&TruckFactory{})
	
	// Document Factory example
	fmt.Println("\n--- Document Factory Example ---")
	workWithDocument(&PDFDocumentFactory{})
	workWithDocument(&WordDocumentFactory{})
	workWithDocument(&ExcelDocumentFactory{})
	
	// Shape Factory example
	fmt.Println("\n--- Shape Factory Example ---")
	drawShape(NewCircleFactory(5.0))
	drawShape(NewRectangleFactory(4.0, 6.0))
	drawShape(NewTriangleFactory(3.0, 4.0))
	
	// Payment Method Factory example
	fmt.Println("\n--- Payment Method Factory Example ---")
	processPayment(NewCreditCardFactory("1234567890123456", "John Doe"), 100.50)
	processPayment(NewPayPalFactory("john.doe@example.com"), 75.25)
	processPayment(NewBankTransferFactory("9876543210", "National Bank"), 200.00)
	
	fmt.Println("\nAll factory method patterns demonstrated successfully!")
}
