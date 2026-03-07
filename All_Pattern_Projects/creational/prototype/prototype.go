package main

import "fmt"

// Prototype Pattern

// Prototype interface
type Prototype interface {
	Clone() Prototype
	GetDetails() string
}

// ConcretePrototype1
type ConcretePrototype1 struct {
	name  string
	value int
}

func NewConcretePrototype1(name string, value int) *ConcretePrototype1 {
	return &ConcretePrototype1{name: name, value: value}
}

func (cp *ConcretePrototype1) Clone() Prototype {
	return &ConcretePrototype1{
		name:  cp.name + " (cloned)",
		value: cp.value,
	}
}

func (cp *ConcretePrototype1) GetDetails() string {
	return fmt.Sprintf("ConcretePrototype1: %s, Value: %d", cp.name, cp.value)
}

// ConcretePrototype2
type ConcretePrototype2 struct {
	description string
	data        []string
}

func NewConcretePrototype2(description string, data []string) *ConcretePrototype2 {
	// Deep copy of data slice
	dataCopy := make([]string, len(data))
	copy(dataCopy, data)
	return &ConcretePrototype2{description: description, data: dataCopy}
}

func (cp *ConcretePrototype2) Clone() Prototype {
	// Deep copy of data slice
	dataCopy := make([]string, len(cp.data))
	copy(dataCopy, cp.data)
	return &ConcretePrototype2{
		description: cp.description + " (cloned)",
		data:        dataCopy,
	}
}

func (cp *ConcretePrototype2) GetDetails() string {
	return fmt.Sprintf("ConcretePrototype2: %s, Data: %v", cp.description, cp.data)
}

// Prototype Manager
type PrototypeManager struct {
	prototypes map[string]Prototype
}

func NewPrototypeManager() *PrototypeManager {
	return &PrototypeManager{
		prototypes: make(map[string]Prototype),
	}
}

func (pm *PrototypeManager) AddPrototype(key string, prototype Prototype) {
	pm.prototypes[key] = prototype
}

func (pm *PrototypeManager) GetPrototype(key string) Prototype {
	if prototype, exists := pm.prototypes[key]; exists {
		return prototype.Clone()
	}
	return nil
}

func (pm *PrototypeManager) ListPrototypes() {
	fmt.Println("Available prototypes:")
	for key, prototype := range pm.prototypes {
		fmt.Printf("  %s: %s\n", key, prototype.GetDetails())
	}
}

// Document Prototype Example
type Document struct {
	title    string
	content  string
	author   string
	version  int
	metadata map[string]string
}

func NewDocument(title, content, author string) *Document {
	return &Document{
		title:    title,
		content:  content,
		author:   author,
		version:  1,
		metadata: make(map[string]string),
	}
}

func (d *Document) Clone() Prototype {
	// Deep copy of metadata
	metadataCopy := make(map[string]string)
	for k, v := range d.metadata {
		metadataCopy[k] = v
	}
	
	cloned := &Document{
		title:    d.title,
		content:  d.content,
		author:   d.author,
		version:  d.version + 1,
		metadata: metadataCopy,
	}
	
	return cloned
}

func (d *Document) GetDetails() string {
	return fmt.Sprintf("Document: %s by %s (v%d)", d.title, d.author, d.version)
}

func (d *Document) AddMetadata(key, value string) {
	d.metadata[key] = value
}

func (d *Document) ShowContent() {
	fmt.Printf("Title: %s\n", d.title)
	fmt.Printf("Author: %s\n", d.author)
	fmt.Printf("Version: %d\n", d.version)
	fmt.Printf("Content: %s\n", d.content)
	fmt.Printf("Metadata: %v\n", d.metadata)
}

// Shape Prototype Example
type Shape interface {
	Prototype
	Draw()
	Move(x, y int)
}

type Circle struct {
	x      int
	y      int
	radius int
	color  string
}

func NewCircle(x, y, radius int, color string) *Circle {
	return &Circle{x: x, y: y, radius: radius, color: color}
}

func (c *Circle) Clone() Prototype {
	return &Circle{
		x:      c.x,
		y:      c.y,
		radius: c.radius,
		color:  c.color,
	}
}

func (c *Circle) GetDetails() string {
	return fmt.Sprintf("Circle at (%d,%d) with radius %d and color %s", 
		c.x, c.y, c.radius, c.color)
}

func (c *Circle) Draw() {
	fmt.Printf("Drawing %s\n", c.GetDetails())
}

func (c *Circle) Move(x, y int) {
	c.x = x
	c.y = y
}

type Rectangle struct {
	x      int
	y      int
	width  int
	height int
	color  string
}

func NewRectangle(x, y, width, height int, color string) *Rectangle {
	return &Rectangle{x: x, y: y, width: width, height: height, color: color}
}

func (r *Rectangle) Clone() Prototype {
	return &Rectangle{
		x:      r.x,
		y:      r.y,
		width:  r.width,
		height: r.height,
		color:  r.color,
	}
}

func (r *Rectangle) GetDetails() string {
	return fmt.Sprintf("Rectangle at (%d,%d) %dx%d with color %s", 
		r.x, r.y, r.width, r.height, r.color)
}

func (r *Rectangle) Draw() {
	fmt.Printf("Drawing %s\n", r.GetDetails())
}

func (r *Rectangle) Move(x, y int) {
	r.x = x
	r.y = y
}

// Shape Registry
type ShapeRegistry struct {
	shapes map[string]Shape
}

func NewShapeRegistry() *ShapeRegistry {
	return &ShapeRegistry{
		shapes: make(map[string]Shape),
	}
}

func (sr *ShapeRegistry) AddShape(key string, shape Shape) {
	sr.shapes[key] = shape
}

func (sr *ShapeRegistry) GetShape(key string) Shape {
	if shape, exists := sr.shapes[key]; exists {
		return shape.Clone().(Shape)
	}
	return nil
}

func (sr *ShapeRegistry) ListShapes() {
	fmt.Println("Available shapes:")
	for key, shape := range sr.shapes {
		fmt.Printf("  %s: %s\n", key, shape.GetDetails())
	}
}

// Employee Prototype Example
type Employee struct {
	name     string
	position string
	salary   float64
	address  *Address
}

type Address struct {
	street string
	city   string
	zip    string
}

func NewEmployee(name, position string, salary float64, address *Address) *Employee {
	return &Employee{
		name:     name,
		position: position,
		salary:   salary,
		address:  address,
	}
}

func (e *Employee) Clone() Prototype {
	// Deep copy of address
	addressCopy := &Address{
		street: e.address.street,
		city:   e.address.city,
		zip:    e.address.zip,
	}
	
	return &Employee{
		name:     e.name,
		position: e.position,
		salary:   e.salary,
		address:  addressCopy,
	}
}

func (e *Employee) GetDetails() string {
	return fmt.Sprintf("Employee: %s, %s, Salary: $%.2f", e.name, e.position, e.salary)
}

func (e *Employee) ShowAddress() {
	fmt.Printf("Address: %s, %s, %s\n", e.address.street, e.address.city, e.address.zip)
}

func (e *Employee) SetPosition(position string) {
	e.position = position
}

func (e *Employee) SetSalary(salary float64) {
	e.salary = salary
}

func (e *Employee) SetAddress(address *Address) {
	e.address = address
}

// Employee Registry
type EmployeeRegistry struct {
	employees map[string]*Employee
}

func NewEmployeeRegistry() *EmployeeRegistry {
	return &EmployeeRegistry{
		employees: make(map[string]*Employee),
	}
}

func (er *EmployeeRegistry) AddEmployee(key string, employee *Employee) {
	er.employees[key] = employee
}

func (er *EmployeeRegistry) GetEmployee(key string) *Employee {
	if employee, exists := er.employees[key]; exists {
		return employee.Clone().(*Employee)
	}
	return nil
}

func (er *EmployeeRegistry) ListEmployees() {
	fmt.Println("Available employee templates:")
	for key, employee := range er.employees {
		fmt.Printf("  %s: %s\n", key, employee.GetDetails())
	}
}

// Game Character Prototype Example
type GameCharacter struct {
	name      string
	class     string
	level     int
	health    int
	mana      int
	strength  int
	agility   int
	equipment []string
}

func NewGameCharacter(name, class string, level int) *GameCharacter {
	return &GameCharacter{
		name:      name,
		class:     class,
		level:     level,
		health:    100,
		mana:      50,
		strength:  10,
		agility:   10,
		equipment: []string{"Sword", "Shield"},
	}
}

func (gc *GameCharacter) Clone() Prototype {
	// Deep copy of equipment
	equipmentCopy := make([]string, len(gc.equipment))
	copy(equipmentCopy, gc.equipment)
	
	return &GameCharacter{
		name:      gc.name,
		class:     gc.class,
		level:     gc.level,
		health:    gc.health,
		mana:      gc.mana,
		strength:  gc.strength,
		agility:   gc.agility,
		equipment: equipmentCopy,
	}
}

func (gc *GameCharacter) GetDetails() string {
	return fmt.Sprintf("Character: %s (%s), Level %d, HP: %d, MP: %d", 
		gc.name, gc.class, gc.level, gc.health, gc.mana)
}

func (gc *GameCharacter) ShowStats() {
	fmt.Printf("%s\n", gc.GetDetails())
	fmt.Printf("  STR: %d, AGI: %d\n", gc.strength, gc.agility)
	fmt.Printf("  Equipment: %v\n", gc.equipment)
}

func (gc *GameCharacter) LevelUp() {
	gc.level++
	gc.health += 20
	gc.mana += 10
	gc.strength += 2
	gc.agility += 2
}

func (gc *GameCharacter) AddEquipment(item string) {
	gc.equipment = append(gc.equipment, item)
}

// Character Registry
type CharacterRegistry struct {
	characters map[string]*GameCharacter
}

func NewCharacterRegistry() *CharacterRegistry {
	return &CharacterRegistry{
		characters: make(map[string]*GameCharacter),
	}
}

func (cr *CharacterRegistry) AddCharacter(key string, character *GameCharacter) {
	cr.characters[key] = character
}

func (cr *CharacterRegistry) GetCharacter(key string) *GameCharacter {
	if character, exists := cr.characters[key]; exists {
		return character.Clone().(*GameCharacter)
	}
	return nil
}

func (cr *CharacterRegistry) ListCharacters() {
	fmt.Println("Available character templates:")
	for key, character := range cr.characters {
		fmt.Printf("  %s: %s\n", key, character.GetDetails())
	}
}

func main() {
	fmt.Println("=== Prototype Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Prototype Example ---")
	manager := NewPrototypeManager()
	
	prototype1 := NewConcretePrototype1("Original", 100)
	prototype2 := NewConcretePrototype2("Original Data", []string{"A", "B", "C"})
	
	manager.AddPrototype("type1", prototype1)
	manager.AddPrototype("type2", prototype2)
	
	manager.ListPrototypes()
	
	clone1 := manager.GetPrototype("type1")
	clone2 := manager.GetPrototype("type2")
	
	fmt.Printf("Clone 1: %s\n", clone1.GetDetails())
	fmt.Printf("Clone 2: %s\n", clone2.GetDetails())
	
	// Document Prototype example
	fmt.Println("\n--- Document Prototype Example ---")
	originalDoc := NewDocument("Project Proposal", "This is a project proposal document...", "John Smith")
	originalDoc.AddMetadata("department", "Engineering")
	originalDoc.AddMetadata("priority", "High")
	
	fmt.Println("Original document:")
	originalDoc.ShowContent()
	
	clonedDoc := originalDoc.Clone().(*Document)
	clonedDoc.title = "Project Proposal (Copy)"
	clonedDoc.content = "This is a copy of the project proposal document..."
	
	fmt.Println("\nCloned document:")
	clonedDoc.ShowContent()
	
	// Shape Prototype example
	fmt.Println("\n--- Shape Prototype Example ---")
	shapeRegistry := NewShapeRegistry()
	
	circle := NewCircle(10, 20, 15, "Red")
	rectangle := NewRectangle(30, 40, 50, 60, "Blue")
	
	shapeRegistry.AddShape("circle", circle)
	shapeRegistry.AddShape("rectangle", rectangle)
	
	shapeRegistry.ListShapes()
	
	clonedCircle := shapeRegistry.GetShape("circle")
	clonedCircle.Move(50, 60)
	clonedCircle.Draw()
	
	clonedRectangle := shapeRegistry.GetShape("rectangle")
	clonedRectangle.Draw()
	
	// Employee Prototype example
	fmt.Println("\n--- Employee Prototype Example ---")
	employeeRegistry := NewEmployeeRegistry()
	
	address := &Address{street: "123 Main St", city: "New York", zip: "10001"}
	developer := NewEmployee("Alice Johnson", "Senior Developer", 95000.0, address)
	
	employeeRegistry.AddEmployee("developer", developer)
	employeeRegistry.ListEmployees()
	
	cloneDeveloper := employeeRegistry.GetEmployee("developer")
	cloneDeveloper.name = "Bob Smith"
	cloneDeveloper.SetPosition("Junior Developer")
	cloneDeveloper.SetSalary(75000.0)
	
	fmt.Printf("Original employee: %s\n", developer.GetDetails())
	fmt.Printf("Cloned employee: %s\n", cloneDeveloper.GetDetails())
	
	// Show that address is independent
	cloneDeveloper.SetAddress(&Address{street: "456 Oak Ave", city: "Boston", zip: "02101"})
	fmt.Printf("Original employee address: ")
	developer.ShowAddress()
	fmt.Printf("Cloned employee address: ")
	cloneDeveloper.ShowAddress()
	
	// Game Character Prototype example
	fmt.Println("\n--- Game Character Prototype Example ---")
	characterRegistry := NewCharacterRegistry()
	
	warrior := NewGameCharacter("Aragorn", "Warrior", 1)
	warrior.strength = 15
	warrior.agility = 8
	warrior.equipment = []string{"Longsword", "Heavy Armor", "Shield"}
	
	mage := NewGameCharacter("Gandalf", "Mage", 1)
	mage.health = 60
	mage.mana = 100
	mage.strength = 5
	mage.agility = 6
	mage.equipment = []string{"Magic Staff", "Robes", "Spellbook"}
	
	characterRegistry.AddCharacter("warrior", warrior)
	characterRegistry.AddCharacter("mage", mage)
	
	characterRegistry.ListCharacters()
	
	cloneWarrior := characterRegistry.GetCharacter("warrior")
	cloneWarrior.name = "Boromir"
	cloneWarrior.LevelUp()
	cloneWarrior.AddEquipment("Horn of Gondor")
	
	cloneMage := characterRegistry.GetCharacter("mage")
	cloneMage.name = "Saruman"
	cloneMage.LevelUp()
	cloneMage.AddEquipment("Palantir")
	
	fmt.Println("\nCloned characters:")
	cloneWarrior.ShowStats()
	cloneMage.ShowStats()
	
	fmt.Println("\nAll prototype patterns demonstrated successfully!")
}
