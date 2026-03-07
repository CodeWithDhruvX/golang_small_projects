package main

import "fmt"

// Abstract Factory Pattern

// AbstractFactory interface declares a set of methods that return abstract products
type AbstractFactory interface {
	CreateProductA() ProductA
	CreateProductB() ProductB
}

// AbstractProductA interface
type ProductA interface {
	UseA()
}

// AbstractProductB interface
type ProductB interface {
	UseB()
	InteractWith(productA ProductA)
}

// ConcreteFactory1 produces products of family 1
type ConcreteFactory1 struct{}

func (cf *ConcreteFactory1) CreateProductA() ProductA {
	return &ConcreteProductA1{}
}

func (cf *ConcreteFactory1) CreateProductB() ProductB {
	return &ConcreteProductB1{}
}

// ConcreteFactory2 produces products of family 2
type ConcreteFactory2 struct{}

func (cf *ConcreteFactory2) CreateProductA() ProductA {
	return &ConcreteProductA2{}
}

func (cf *ConcreteFactory2) CreateProductB() ProductB {
	return &ConcreteProductB2{}
}

// ConcreteProductA1
type ConcreteProductA1 struct{}

func (cpa *ConcreteProductA1) UseA() {
	fmt.Println("ConcreteProductA1: Used")
}

// ConcreteProductB1
type ConcreteProductB1 struct{}

func (cpb *ConcreteProductB1) UseB() {
	fmt.Println("ConcreteProductB1: Used")
}

func (cpb *ConcreteProductB1) InteractWith(productA ProductA) {
	fmt.Printf("ConcreteProductB1 interacts with %T\n", productA)
	productA.UseA()
}

// ConcreteProductA2
type ConcreteProductA2 struct{}

func (cpa *ConcreteProductA2) UseA() {
	fmt.Println("ConcreteProductA2: Used")
}

// ConcreteProductB2
type ConcreteProductB2 struct{}

func (cpb *ConcreteProductB2) UseB() {
	fmt.Println("ConcreteProductB2: Used")
}

func (cpb *ConcreteProductB2) InteractWith(productA ProductA) {
	fmt.Printf("ConcreteProductB2 interacts with %T\n", productA)
	productA.UseA()
}

// Client code works with factories and products only through abstract types
func clientCode(factory AbstractFactory) {
	productA := factory.CreateProductA()
	productB := factory.CreateProductB()
	
	fmt.Println("Client: Using products from the same factory:")
	productB.UseB()
	productB.InteractWith(productA)
}

// GUI Widget Example
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

type Button interface {
	Render()
	OnClick()
}

type Checkbox interface {
	Render()
	Toggle()
}

// Windows GUI Factory
type WindowsFactory struct{}

func (wf *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (wf *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

type WindowsButton struct{}

func (wb *WindowsButton) Render() {
	fmt.Println("Windows Button: Rendered with Windows style")
}

func (wb *WindowsButton) OnClick() {
	fmt.Println("Windows Button: Clicked with Windows behavior")
}

type WindowsCheckbox struct{}

func (wc *WindowsCheckbox) Render() {
	fmt.Println("Windows Checkbox: Rendered with Windows style")
}

func (wc *WindowsCheckbox) Toggle() {
	fmt.Println("Windows Checkbox: Toggled with Windows behavior")
}

// macOS GUI Factory
class MacOSFactory struct{}

func (mf *MacOSFactory) CreateButton() Button {
	return &MacOSButton{}
}

func (mf *MacOSFactory) CreateCheckbox() Checkbox {
	return &MacOSCheckbox{}
}

type MacOSButton struct{}

func (mb *MacOSButton) Render() {
	fmt.Println("macOS Button: Rendered with macOS style")
}

func (mb *MacOSButton) OnClick() {
	fmt.Println("macOS Button: Clicked with macOS behavior")
}

type MacOSCheckbox struct{}

func (mc *MacOSCheckbox) Render() {
	fmt.Println("macOS Checkbox: Rendered with macOS style")
}

func (mc *MacOSCheckbox) Toggle() {
	fmt.Println("macOS Checkbox: Toggled with macOS behavior")
}

// GUI Application
type GUIApplication struct {
	button   Button
	checkbox Checkbox
}

func NewGUIApplication(factory GUIFactory) *GUIApplication {
	return &GUIApplication{
		button:   factory.CreateButton(),
		checkbox: factory.CreateCheckbox(),
	}
}

func (app *GUIApplication) RenderUI() {
	fmt.Println("Rendering UI components:")
	app.button.Render()
	app.checkbox.Render()
}

func (app *GUIApplication) SimulateInteraction() {
	fmt.Println("\nSimulating user interactions:")
	app.button.OnClick()
	app.checkbox.Toggle()
}

// Database Connection Example
type DatabaseFactory interface {
	CreateConnection() Connection
	CreateQuery() Query
}

type Connection interface {
	Connect()
	Disconnect()
	Execute(query string)
}

type Query interface {
	SetSQL(sql string)
	GetSQL() string
	Execute()
}

// MySQL Factory
type MySQLFactory struct{}

func (mf *MySQLFactory) CreateConnection() Connection {
	return &MySQLConnection{}
}

func (mf *MySQLFactory) CreateQuery() Query {
	return &MySQLQuery{}
}

type MySQLConnection struct {
	connected bool
}

func (mc *MySQLConnection) Connect() {
	mc.connected = true
	fmt.Println("MySQL Connection: Connected to MySQL database")
}

func (mc *MySQLConnection) Disconnect() {
	mc.connected = false
	fmt.Println("MySQL Connection: Disconnected from MySQL database")
}

func (mc *MySQLConnection) Execute(query string) {
	if mc.connected {
		fmt.Printf("MySQL Connection: Executing query: %s\n", query)
	} else {
		fmt.Println("MySQL Connection: Not connected")
	}
}

type MySQLQuery struct {
	sql string
}

func (mq *MySQLQuery) SetSQL(sql string) {
	mq.sql = sql
}

func (mq *MySQLQuery) GetSQL() string {
	return mq.sql
}

func (mq *MySQLQuery) Execute() {
	fmt.Printf("MySQL Query: Executing SQL: %s\n", mq.sql)
}

// PostgreSQL Factory
type PostgreSQLFactory struct{}

func (pf *PostgreSQLFactory) CreateConnection() Connection {
	return &PostgreSQLConnection{}
}

func (pf *PostgreSQLFactory) CreateQuery() Query {
	return &PostgreSQLQuery{}
}

type PostgreSQLConnection struct {
	connected bool
}

func (pc *PostgreSQLConnection) Connect() {
	pc.connected = true
	fmt.Println("PostgreSQL Connection: Connected to PostgreSQL database")
}

func (pc *PostgreSQLConnection) Disconnect() {
	pc.connected = false
	fmt.Println("PostgreSQL Connection: Disconnected from PostgreSQL database")
}

func (pc *PostgreSQLConnection) Execute(query string) {
	if pc.connected {
		fmt.Printf("PostgreSQL Connection: Executing query: %s\n", query)
	} else {
		fmt.Println("PostgreSQL Connection: Not connected")
	}
}

type PostgreSQLQuery struct {
	sql string
}

func (pq *PostgreSQLQuery) SetSQL(sql string) {
	pq.sql = sql
}

func (pq *PostgreSQLQuery) GetSQL() string {
	return pq.sql
}

func (pq *PostgreSQLQuery) Execute() {
	fmt.Printf("PostgreSQL Query: Executing SQL: %s\n", pq.sql)
}

// Database Application
type DatabaseApplication struct {
	connection Connection
	query      Query
}

func NewDatabaseApplication(factory DatabaseFactory) *DatabaseApplication {
	return &DatabaseApplication{
		connection: factory.CreateConnection(),
		query:      factory.CreateQuery(),
	}
}

func (da *DatabaseApplication) RunDatabaseOperations() {
	fmt.Println("Starting database operations:")
	da.connection.Connect()
	
	da.query.SetSQL("SELECT * FROM users")
	da.query.Execute()
	da.connection.Execute(da.query.GetSQL())
	
	da.connection.Disconnect()
}

func main() {
	fmt.Println("=== Abstract Factory Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Abstract Factory Example ---")
	fmt.Println("Client: Testing with ConcreteFactory1:")
	clientCode(&ConcreteFactory1{})
	
	fmt.Println("\nClient: Testing with ConcreteFactory2:")
	clientCode(&ConcreteFactory2{})
	
	// GUI Widget example
	fmt.Println("\n--- GUI Widget Example ---")
	fmt.Println("Creating Windows UI:")
	windowsApp := NewGUIApplication(&WindowsFactory{})
	windowsApp.RenderUI()
	windowsApp.SimulateInteraction()
	
	fmt.Println("\nCreating macOS UI:")
	macOSApp := NewGUIApplication(&MacOSFactory{})
	macOSApp.RenderUI()
	macOSApp.SimulateInteraction()
	
	// Database Connection example
	fmt.Println("\n--- Database Connection Example ---")
	fmt.Println("Using MySQL:")
	mysqlApp := NewDatabaseApplication(&MySQLFactory{})
	mysqlApp.RunDatabaseOperations()
	
	fmt.Println("\nUsing PostgreSQL:")
	postgreSQLApp := NewDatabaseApplication(&PostgreSQLFactory{})
	postgreSQLApp.RunDatabaseOperations()
	
	fmt.Println("\nAll abstract factory patterns demonstrated successfully!")
}
