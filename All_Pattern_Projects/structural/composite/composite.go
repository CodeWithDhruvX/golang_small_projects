package main

import "fmt"

// Composite Pattern

// Component interface
type Component interface {
	Operation() string
	Add(component Component)
	Remove(component Component)
	GetChild(index int) Component
}

// Leaf component
type Leaf struct {
	name string
}

func NewLeaf(name string) *Leaf {
	return &Leaf{name: name}
}

func (l *Leaf) Operation() string {
	return fmt.Sprintf("Leaf: %s", l.name)
}

func (l *Leaf) Add(component Component) {
	fmt.Println("Cannot add to leaf")
}

func (l *Leaf) Remove(component Component) {
	fmt.Println("Cannot remove from leaf")
}

func (l *Leaf) GetChild(index int) Component {
	return nil
}

// Composite component
type Composite struct {
	name      string
	children  []Component
}

func NewComposite(name string) *Composite {
	return &Composite{
		name:     name,
		children: make([]Component, 0),
	}
}

func (c *Composite) Operation() string {
	result := fmt.Sprintf("Composite: %s\n", c.name)
	for _, child := range c.children {
		result += "  " + child.Operation() + "\n"
	}
	return result
}

func (c *Composite) Add(component Component) {
	c.children = append(c.children, component)
}

func (c *Composite) Remove(component Component) {
	for i, child := range c.children {
		if child == component {
			c.children = append(c.children[:i], c.children[i+1:]...)
			break
		}
	}
}

func (c *Composite) GetChild(index int) Component {
	if index >= 0 && index < len(c.children) {
		return c.children[index]
	}
	return nil
}

// File System Example
type FileSystemComponent interface {
	Show(indent string)
	GetSize() int
}

type File struct {
	name string
	size int
}

func NewFile(name string, size int) *File {
	return &File{name: name, size: size}
}

func (f *File) Show(indent string) {
	fmt.Printf("%sFile: %s (%d bytes)\n", indent, f.name, f.size)
}

func (f *File) GetSize() int {
	return f.size
}

type Directory struct {
	name     string
	children []FileSystemComponent
}

func NewDirectory(name string) *Directory {
	return &Directory{
		name:     name,
		children: make([]FileSystemComponent, 0),
	}
}

func (d *Directory) Show(indent string) {
	fmt.Printf("%sDirectory: %s (%d bytes)\n", indent, d.name, d.GetSize())
	indent += "  "
	for _, child := range d.children {
		child.Show(indent)
	}
}

func (d *Directory) GetSize() int {
	total := 0
	for _, child := range d.children {
		total += child.GetSize()
	}
	return total
}

func (d *Directory) Add(component FileSystemComponent) {
	d.children = append(d.children, component)
}

func (d *Directory) Remove(component FileSystemComponent) {
	for i, child := range d.children {
		if child == component {
			d.children = append(d.children[:i], d.children[i+1:]...)
			break
		}
	}
}

// Organization Structure Example
type Employee interface {
	ShowDetails(indent string)
	GetTotalSalary() float64
}

type IndividualEmployee struct {
	name   string
	title  string
	salary float64
}

func NewIndividualEmployee(name, title string, salary float64) *IndividualEmployee {
	return &IndividualEmployee{name: name, title: title, salary: salary}
}

func (ie *IndividualEmployee) ShowDetails(indent string) {
	fmt.Printf("%s%s - %s ($%.2f)\n", indent, ie.name, ie.title, ie.salary)
}

func (ie *IndividualEmployee) GetTotalSalary() float64 {
	return ie.salary
}

type Manager struct {
	name     string
	title    string
	salary   float64
	team     []Employee
}

func NewManager(name, title string, salary float64) *Manager {
	return &Manager{
		name:   name,
		title:  title,
		salary: salary,
		team:   make([]Employee, 0),
	}
}

func (m *Manager) ShowDetails(indent string) {
	fmt.Printf("%s%s - %s ($%.2f) [Team Size: %d]\n", indent, m.name, m.title, m.salary, len(m.team))
	indent += "  "
	for _, member := range m.team {
		member.ShowDetails(indent)
	}
}

func (m *Manager) GetTotalSalary() float64 {
	total := m.salary
	for _, member := range m.team {
		total += member.GetTotalSalary()
	}
	return total
}

func (m *Manager) AddToTeam(employee Employee) {
	m.team = append(m.team, employee)
}

func (m *Manager) RemoveFromTeam(employee Employee) {
	for i, member := range m.team {
		if member == employee {
			m.team = append(m.team[:i], m.team[i+1:]...)
			break
		}
	}
}

// Graphic Shape Example
type Graphic interface {
	Draw(indent string)
	GetArea() float64
}

type Circle struct {
	radius float64
	color  string
}

func NewCircle(radius float64, color string) *Circle {
	return &Circle{radius: radius, color: color}
}

func (c *Circle) Draw(indent string) {
	fmt.Printf("%sCircle (r=%.1f, %s) - Area: %.2f\n", indent, c.radius, c.color, c.GetArea())
}

func (c *Circle) GetArea() float64 {
	return 3.14159 * c.radius * c.radius
}

type Rectangle struct {
	width  float64
	height float64
	color  string
}

func NewRectangle(width, height float64, color string) *Rectangle {
	return &Rectangle{width: width, height: height, color: color}
}

func (r *Rectangle) Draw(indent string) {
	fmt.Printf("%sRectangle (%.1fx%.1f, %s) - Area: %.2f\n", indent, r.width, r.height, r.color, r.GetArea())
}

func (r *Rectangle) GetArea() float64 {
	return r.width * r.height
}

type Group struct {
	name    string
	shapes  []Graphic
}

func NewGroup(name string) *Group {
	return &Group{name: name, shapes: make([]Graphic, 0)}
}

func (g *Group) Draw(indent string) {
	fmt.Printf("%sGroup: %s (Total Area: %.2f)\n", indent, g.name, g.GetArea())
	indent += "  "
	for _, shape := range g.shapes {
		shape.Draw(indent)
	}
}

func (g *Group) GetArea() float64 {
	total := 0.0
	for _, shape := range g.shapes {
		total += shape.GetArea()
	}
	return total
}

func (g *Group) Add(shape Graphic) {
	g.shapes = append(g.shapes, shape)
}

func (g *Group) Remove(shape Graphic) {
	for i, s := range g.shapes {
		if s == shape {
			g.shapes = append(g.shapes[:i], g.shapes[i+1:]...)
			break
		}
	}
}

// Menu System Example
type MenuItem interface {
	Display(indent string)
	GetPrice() float64
}

type SimpleMenuItem struct {
	name  string
	price float64
}

func NewSimpleMenuItem(name string, price float64) *SimpleMenuItem {
	return &SimpleMenuItem{name: name, price: price}
}

func (smi *SimpleMenuItem) Display(indent string) {
	fmt.Printf("%s%s - $%.2f\n", indent, smi.name, smi.price)
}

func (smi *SimpleMenuItem) GetPrice() float64 {
	return smi.price
}

type MenuCategory struct {
	name     string
	items    []MenuItem
}

func NewMenuCategory(name string) *MenuCategory {
	return &MenuCategory{name: name, items: make([]MenuItem, 0)}
}

func (mc *MenuCategory) Display(indent string) {
	fmt.Printf("%s%s (Total: $%.2f)\n", indent, mc.name, mc.GetPrice())
	indent += "  "
	for _, item := range mc.items {
		item.Display(indent)
	}
}

func (mc *MenuCategory) GetPrice() float64 {
	total := 0.0
	for _, item := range mc.items {
		total += item.GetPrice()
	}
	return total
}

func (mc *MenuCategory) Add(item MenuItem) {
	mc.items = append(mc.items, item)
}

func (mc *MenuCategory) Remove(item MenuItem) {
	for i, menuItem := range mc.items {
		if menuItem == item {
			mc.items = append(mc.items[:i], mc.items[i+1:]...)
			break
		}
	}
}

// Task Management Example
type Task interface {
	Execute(indent string)
	GetDuration() int
}

type SimpleTask struct {
	name     string
	duration int // in minutes
}

func NewSimpleTask(name string, duration int) *SimpleTask {
	return &SimpleTask{name: name, duration: duration}
}

func (st *SimpleTask) Execute(indent string) {
	fmt.Printf("%sExecuting task: %s (%d min)\n", indent, st.name, st.duration)
}

func (st *SimpleTask) GetDuration() int {
	return st.duration
}

type TaskGroup struct {
	name     string
	tasks    []Task
}

func NewTaskGroup(name string) *TaskGroup {
	return &TaskGroup{name: name, tasks: make([]Task, 0)}
}

func (tg *TaskGroup) Execute(indent string) {
	fmt.Printf("%sExecuting task group: %s (Total: %d min)\n", indent, tg.name, tg.GetDuration())
	indent += "  "
	for _, task := range tg.tasks {
		task.Execute(indent)
	}
}

func (tg *TaskGroup) GetDuration() int {
	total := 0
	for _, task := range tg.tasks {
		total += task.GetDuration()
	}
	return total
}

func (tg *TaskGroup) Add(task Task) {
	tg.tasks = append(tg.tasks, task)
}

func (tg *TaskGroup) Remove(task Task) {
	for i, t := range tg.tasks {
		if t == task {
			tg.tasks = append(tg.tasks[:i], t.tasks[i+1:]...)
			break
		}
	}
}

func main() {
	fmt.Println("=== Composite Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Composite Example ---")
	
	leaf1 := NewLeaf("Leaf 1")
	leaf2 := NewLeaf("Leaf 2")
	leaf3 := NewLeaf("Leaf 3")
	
	composite1 := NewComposite("Composite 1")
	composite1.Add(leaf1)
	composite1.Add(leaf2)
	
	composite2 := NewComposite("Composite 2")
	composite2.Add(leaf3)
	
	root := NewComposite("Root")
	root.Add(composite1)
	root.Add(composite2)
	
	fmt.Println(root.Operation())
	
	// File System example
	fmt.Println("\n--- File System Example ---")
	
	file1 := NewFile("document.txt", 1024)
	file2 := NewFile("image.jpg", 2048)
	file3 := NewFile("video.mp4", 10240)
	
	docsDir := NewDirectory("Documents")
	docsDir.Add(file1)
	docsDir.Add(file2)
	
	mediaDir := NewDirectory("Media")
	mediaDir.Add(file3)
	
	rootDir := NewDirectory("Root")
	rootDir.Add(docsDir)
	rootDir.Add(mediaDir)
	
	rootDir.Show("")
	fmt.Printf("Total size: %d bytes\n", rootDir.GetSize())
	
	// Organization Structure example
	fmt.Println("\n--- Organization Structure Example ---")
	
	developer1 := NewIndividualEmployee("Alice", "Developer", 75000)
	developer2 := NewIndividualEmployee("Bob", "Developer", 80000)
	designer := NewIndividualEmployee("Carol", "Designer", 65000)
	
	teamLead := NewManager("David", "Team Lead", 95000)
	teamLead.AddToTeam(developer1)
	teamLead.AddToTeam(developer2)
	teamLead.AddToTeam(designer)
	
	manager := NewManager("Eve", "Engineering Manager", 120000)
	manager.AddToTeam(teamLead)
	
	manager.ShowDetails("")
	fmt.Printf("Total team salary: $%.2f\n", manager.GetTotalSalary())
	
	// Graphic Shape example
	fmt.Println("\n--- Graphic Shape Example ---")
	
	circle1 := NewCircle(5.0, "red")
	circle2 := NewCircle(3.0, "blue")
	rectangle1 := NewRectangle(10.0, 5.0, "green")
	rectangle2 := NewRectangle(8.0, 6.0, "yellow")
	
	group1 := NewGroup("Basic Shapes")
	group1.Add(circle1)
	group1.Add(rectangle1)
	
	group2 := NewGroup("Advanced Shapes")
	group2.Add(circle2)
	group2.Add(rectangle2)
	
	mainGroup := NewGroup("All Shapes")
	mainGroup.Add(group1)
	mainGroup.Add(group2)
	
	mainGroup.Draw("")
	
	// Menu System example
	fmt.Println("\n--- Menu System Example ---")
	
	burger := NewSimpleMenuItem("Burger", 8.99)
	fries := NewSimpleMenuItem("Fries", 3.99)
	shake := NewSimpleMenuItem("Shake", 4.99)
	
	salad := NewSimpleMenuItem("Salad", 6.99)
	soup := NewSimpleMenuItem("Soup", 4.49)
	
	mainCourse := NewMenuCategory("Main Course")
	mainCourse.Add(burger)
	mainCourse.Add(salad)
	
	sides := NewMenuCategory("Sides")
	sides.Add(fries)
	sides.Add(soup)
	
	beverages := NewMenuCategory("Beverages")
	beverages.Add(shake)
	
	restaurantMenu := NewMenuCategory("Restaurant Menu")
	restaurantMenu.Add(mainCourse)
	restaurantMenu.Add(sides)
	restaurantMenu.Add(beverages)
	
	restaurantMenu.Display("")
	fmt.Printf("Total menu price: $%.2f\n", restaurantMenu.GetPrice())
	
	// Task Management example
	fmt.Println("\n--- Task Management Example ---")
	
	task1 := NewSimpleTask("Code Review", 30)
	task2 := NewSimpleTask("Write Tests", 45)
	task3 := NewSimpleTask("Documentation", 20)
	
	sprintTasks := NewTaskGroup("Sprint Tasks")
	sprintTasks.Add(task1)
	sprintTasks.Add(task2)
	sprintTasks.Add(task3)
	
	task4 := NewSimpleTask("Meeting", 60)
	task5 := NewSimpleTask("Planning", 90)
	
	planningTasks := NewTaskGroup("Planning Tasks")
	planningTasks.Add(task4)
	planningTasks.Add(task5)
	
	allTasks := NewTaskGroup("Project Tasks")
	allTasks.Add(sprintTasks)
	allTasks.Add(planningTasks)
	
	allTasks.Execute("")
	fmt.Printf("Total duration: %d minutes\n", allTasks.GetDuration())
	
	fmt.Println("\nAll composite patterns demonstrated successfully!")
}
