package main

import "fmt"

// Decorator Pattern

// Component interface
type Component interface {
	Operation() string
}

// Concrete Component
type ConcreteComponent struct{}

func (cc *ConcreteComponent) Operation() string {
	return "ConcreteComponent"
}

// Base Decorator
type BaseDecorator struct {
	component Component
}

func (bd *BaseDecorator) Operation() string {
	return bd.component.Operation()
}

// Concrete Decorators
type ConcreteDecoratorA struct {
	BaseDecorator
	addedState string
}

func NewConcreteDecoratorA(component Component) *ConcreteDecoratorA {
	return &ConcreteDecoratorA{
		BaseDecorator: BaseDecorator{component: component},
		addedState:    "Added State A",
	}
}

func (cda *ConcreteDecoratorA) Operation() string {
	return fmt.Sprintf("ConcreteDecoratorA(%s) [%s]", cda.BaseDecorator.Operation(), cda.addedState)
}

type ConcreteDecoratorB struct {
	BaseDecorator
}

func NewConcreteDecoratorB(component Component) *ConcreteDecoratorB {
	return &ConcreteDecoratorB{
		BaseDecorator: BaseDecorator{component: component},
	}
}

func (cdb *ConcreteDecoratorB) Operation() string {
	return fmt.Sprintf("ConcreteDecoratorB(%s)", cdb.BaseDecorator.Operation())
}

func (cdb *ConcreteDecoratorB) AddedBehavior() string {
	return "Added Behavior B"
}

// Coffee Shop Example
type Coffee interface {
	GetCost() float64
	GetDescription() string
}

type SimpleCoffee struct{}

func (sc *SimpleCoffee) GetCost() float64 {
	return 2.0
}

func (sc *SimpleCoffee) GetDescription() string {
	return "Simple coffee"
}

type MilkDecorator struct {
	coffee Coffee
}

func NewMilkDecorator(coffee Coffee) *MilkDecorator {
	return &MilkDecorator{coffee: coffee}
}

func (md *MilkDecorator) GetCost() float64 {
	return md.coffee.GetCost() + 0.5
}

func (md *MilkDecorator) GetDescription() string {
	return md.coffee.GetDescription() + ", milk"
}

type SugarDecorator struct {
	coffee Coffee
}

func NewSugarDecorator(coffee Coffee) *SugarDecorator {
	return &SugarDecorator{coffee: coffee}
}

func (sd *SugarDecorator) GetCost() float64 {
	return sd.coffee.GetCost() + 0.2
}

func (sd *SugarDecorator) GetDescription() string {
	return sd.coffee.GetDescription() + ", sugar"
}

type VanillaDecorator struct {
	coffee Coffee
}

func NewVanillaDecorator(coffee Coffee) *VanillaDecorator {
	return &VanillaDecorator{coffee: coffee}
}

func (vd *VanillaDecorator) GetCost() float64 {
	return vd.coffee.GetCost() + 0.3
}

func (vd *VanillaDecorator) GetDescription() string {
	return vd.coffee.GetDescription() + ", vanilla"
}

type WhippedCreamDecorator struct {
	coffee Coffee
}

func NewWhippedCreamDecorator(coffee Coffee) *WhippedCreamDecorator {
	return &WhippedCreamDecorator{coffee: coffee}
}

func (wcd *WhippedCreamDecorator) GetCost() float64 {
	return wcd.coffee.GetCost() + 0.7
}

func (wcd *WhippedCreamDecorator) GetDescription() string {
	return wcd.coffee.GetDescription() + ", whipped cream"
}

// Pizza Shop Example
type Pizza interface {
	GetPrice() float64
	GetDescription() string
}

type MargheritaPizza struct{}

func (mp *MargheritaPizza) GetPrice() float64 {
	return 8.99
}

func (mp *MargheritaPizza) GetDescription() string {
	return "Margherita Pizza"
}

type PepperoniTopping struct {
	pizza Pizza
}

func NewPepperoniTopping(pizza Pizza) *PepperoniTopping {
	return &PepperoniTopping{pizza: pizza}
}

func (pt *PepperoniTopping) GetPrice() float64 {
	return pt.pizza.GetPrice() + 2.50
}

func (pt *PepperoniTopping) GetDescription() string {
	return pt.pizza.GetDescription() + ", pepperoni"
}

type MushroomTopping struct {
	pizza Pizza
}

func NewMushroomTopping(pizza Pizza) *MushroomTopping {
	return &MushroomTopping{pizza: pizza}
}

func (mt *MushroomTopping) GetPrice() float64 {
	return mt.pizza.GetPrice() + 1.50
}

func (mt *MushroomTopping) GetDescription() string {
	return mt.pizza.GetDescription() + ", mushrooms"
}

type ExtraCheeseTopping struct {
	pizza Pizza
}

func NewExtraCheeseTopping(pizza Pizza) *ExtraCheeseTopping {
	return &ExtraCheeseTopping{pizza: pizza}
}

func (ect *ExtraCheeseTopping) GetPrice() float64 {
	return ect.pizza.GetPrice() + 1.25
}

func (ect *ExtraCheeseTopping) GetDescription() string {
	return ect.pizza.GetDescription() + ", extra cheese"
}

// Notification System Example
type Notifier interface {
	Send(message string)
}

type BasicNotifier struct{}

func (bn *BasicNotifier) Send(message string) {
	fmt.Printf("Basic notification: %s\n", message)
}

type SMSNotifier struct {
	notifier Notifier
}

func NewSMSNotifier(notifier Notifier) *SMSNotifier {
	return &SMSNotifier{notifier: notifier}
}

func (sn *SMSNotifier) Send(message string) {
	sn.notifier.Send(message)
	fmt.Printf("SMS sent: %s\n", message)
}

type EmailNotifier struct {
	notifier Notifier
}

func NewEmailNotifier(notifier Notifier) *EmailNotifier {
	return &EmailNotifier{notifier: notifier}
}

func (en *EmailNotifier) Send(message string) {
	en.notifier.Send(message)
	fmt.Printf("Email sent: %s\n", message)
}

type SlackNotifier struct {
	notifier Notifier
}

func NewSlackNotifier(notifier Notifier) *SlackNotifier {
	return &SlackNotifier{notifier: notifier}
}

func (sn *SlackNotifier) Send(message string) {
	sn.notifier.Send(message)
	fmt.Printf("Slack message sent: %s\n", message)
}

type FacebookNotifier struct {
	notifier Notifier
}

func NewFacebookNotifier(notifier Notifier) *FacebookNotifier {
	return &FacebookNotifier{notifier: notifier}
}

func (fn *FacebookNotifier) Send(message string) {
	fn.notifier.Send(message)
	fmt.Printf("Facebook post created: %s\n", message)
}

// Data Stream Example
type DataSource interface {
	WriteData(data string)
	ReadData() string
}

type FileDataSource struct {
	filename string
	data     string
}

func NewFileDataSource(filename string) *FileDataSource {
	return &FileDataSource{filename: filename}
}

func (fds *FileDataSource) WriteData(data string) {
	fds.data = data
	fmt.Printf("Writing to file %s: %s\n", fds.filename, data)
}

func (fds *FileDataSource) ReadData() string {
	fmt.Printf("Reading from file %s: %s\n", fds.filename, fds.data)
	return fds.data
}

type EncryptionDecorator struct {
	wrapper DataSource
}

func NewEncryptionDecorator(wrapper DataSource) *EncryptionDecorator {
	return &EncryptionDecorator{wrapper: wrapper}
}

func (ed *EncryptionDecorator) WriteData(data string) {
	encrypted := "ENCRYPTED:" + data
	ed.wrapper.WriteData(encrypted)
}

func (ed *EncryptionDecorator) ReadData() string {
	data := ed.wrapper.ReadData()
	if len(data) > 10 && data[:10] == "ENCRYPTED:" {
		return data[10:] // Decrypt (simplified)
	}
	return data
}

type CompressionDecorator struct {
	wrapper DataSource
}

func NewCompressionDecorator(wrapper DataSource) *CompressionDecorator {
	return &CompressionDecorator{wrapper: wrapper}
}

func (cd *CompressionDecorator) WriteData(data string) {
	compressed := "COMPRESSED:" + data
	cd.wrapper.WriteData(compressed)
}

func (cd *CompressionDecorator) ReadData() string {
	data := cd.wrapper.ReadData()
	if len(data) > 11 && data[:11] == "COMPRESSED:" {
		return data[11:] // Decompress (simplified)
	}
	return data
}

type LoggingDecorator struct {
	wrapper DataSource
}

func NewLoggingDecorator(wrapper DataSource) *LoggingDecorator {
	return &LoggingDecorator{wrapper: wrapper}
}

func (ld *LoggingDecorator) WriteData(data string) {
	fmt.Printf("LOG: Writing data...\n")
	ld.wrapper.WriteData(data)
	fmt.Printf("LOG: Data written successfully\n")
}

func (ld *LoggingDecorator) ReadData() string {
	fmt.Printf("LOG: Reading data...\n")
	data := ld.wrapper.ReadData()
	fmt.Printf("LOG: Data read successfully\n")
	return data
}

// Weapon Enhancement Example
type Weapon interface {
	GetDamage() int
	GetDescription() string
}

type Sword struct{}

func (s *Sword) GetDamage() int {
	return 10
}

func (s *Sword) GetDescription() string {
	return "Sword"
}

type FireEnchantment struct {
	weapon Weapon
}

func NewFireEnchantment(weapon Weapon) *FireEnchantment {
	return &FireEnchantment{weapon: weapon}
}

func (fe *FireEnchantment) GetDamage() int {
	return fe.weapon.GetDamage() + 5
}

func (fe *FireEnchantment) GetDescription() string {
	return fe.weapon.GetDescription() + " (Fire)"
}

type IceEnchantment struct {
	weapon Weapon
}

func NewIceEnchantment(weapon Weapon) *IceEnchantment {
	return &IceEnchantment{weapon: weapon}
}

func (ie *IceEnchantment) GetDamage() int {
	return ie.weapon.GetDamage() + 4
}

func (ie *IceEnchantment) GetDescription() string {
	return ie.weapon.GetDescription() + " (Ice)"
}

type PoisonEnchantment struct {
	weapon Weapon
}

func NewPoisonEnchantment(weapon Weapon) *PoisonEnchantment {
	return &PoisonEnchantment{weapon: weapon}
}

func (pe *PoisonEnchantment) GetDamage() int {
	return pe.weapon.GetDamage() + 3
}

func (pe *PoisonEnchantment) GetDescription() string {
	return pe.weapon.GetDescription() + " (Poison)"
}

type SharpnessEnchantment struct {
	weapon Weapon
}

func NewSharpnessEnchantment(weapon Weapon) *SharpnessEnchantment {
	return &SharpnessEnchantment{weapon: weapon}
}

func (se *SharpnessEnchantment) GetDamage() int {
	return se.weapon.GetDamage() + 2
}

func (se *SharpnessEnchantment) GetDescription() string {
	return se.weapon.GetDescription() + " (Sharp)"
}

// Car Features Example
type Car interface {
	GetPrice() float64
	GetDescription() string
}

type BasicCar struct{}

func (bc *BasicCar) GetPrice() float64 {
	return 20000.0
}

func (bc *BasicCar) GetDescription() string {
	return "Basic Car"
}

type GPSDecorator struct {
	car Car
}

func NewGPSDecorator(car Car) *GPSDecorator {
	return &GPSDecorator{car: car}
}

func (gps *GPSDecorator) GetPrice() float64 {
	return gps.car.GetPrice() + 1500.0
}

func (gps *GPSDecorator) GetDescription() string {
	return gps.car.GetDescription() + ", GPS"
}

type SunroofDecorator struct {
	car Car
}

func NewSunroofDecorator(car Car) *SunroofDecorator {
	return &SunroofDecorator{car: car}
}

func (sr *SunroofDecorator) GetPrice() float64 {
	return sr.car.GetPrice() + 2000.0
}

func (sr *SunroofDecorator) GetDescription() string {
	return sr.car.GetDescription() + ", Sunroof"
}

type LeatherSeatsDecorator struct {
	car Car
}

func NewLeatherSeatsDecorator(car Car) *LeatherSeatsDecorator {
	return &LeatherSeatsDecorator{car: car}
}

func (ls *LeatherSeatsDecorator) GetPrice() float64 {
	return ls.car.GetPrice() + 3000.0
}

func (ls *LeatherSeatsDecorator) GetDescription() string {
	return ls.car.GetDescription() + ", Leather Seats"
}

type PremiumAudioDecorator struct {
	car Car
}

func NewPremiumAudioDecorator(car Car) *PremiumAudioDecorator {
	return &PremiumAudioDecorator{car: car}
}

func (pa *PremiumAudioDecorator) GetPrice() float64 {
	return pa.car.GetPrice() + 1200.0
}

func (pa *PremiumAudioDecorator) GetDescription() string {
	return pa.car.GetDescription() + ", Premium Audio"
}

func main() {
	fmt.Println("=== Decorator Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Decorator Example ---")
	
	simple := &ConcreteComponent{}
	decorator1 := NewConcreteDecoratorA(simple)
	decorator2 := NewConcreteDecoratorB(decorator1)
	
	fmt.Println("Client: I've got a simple component:")
	fmt.Println("Result: " + simple.Operation())
	
	fmt.Println("\nClient: Now I've got a decorated component:")
	fmt.Println("Result: " + decorator2.Operation())
	fmt.Println("Added behavior: " + decorator2.AddedBehavior())
	
	// Coffee Shop example
	fmt.Println("\n--- Coffee Shop Example ---")
	
	coffee := &SimpleCoffee{}
	fmt.Printf("%s: $%.2f\n", coffee.GetDescription(), coffee.GetCost())
	
	coffeeWithMilk := NewMilkDecorator(coffee)
	fmt.Printf("%s: $%.2f\n", coffeeWithMilk.GetDescription(), coffeeWithMilk.GetCost())
	
	coffeeWithMilkAndSugar := NewSugarDecorator(coffeeWithMilk)
	fmt.Printf("%s: $%.2f\n", coffeeWithMilkAndSugar.GetDescription(), coffeeWithMilkAndSugar.GetCost())
	
	luxuryCoffee := NewWhippedCreamDecorator(NewVanillaDecorator(NewSugarDecorator(NewMilkDecorator(coffee))))
	fmt.Printf("%s: $%.2f\n", luxuryCoffee.GetDescription(), luxuryCoffee.GetCost())
	
	// Pizza Shop example
	fmt.Println("\n--- Pizza Shop Example ---")
	
	pizza := &MargheritaPizza{}
	fmt.Printf("%s: $%.2f\n", pizza.GetDescription(), pizza.GetPrice())
	
	pizzaWithPepperoni := NewPepperoniTopping(pizza)
	fmt.Printf("%s: $%.2f\n", pizzaWithPepperoni.GetDescription(), pizzaWithPepperoni.GetPrice())
	
	supremePizza := NewExtraCheeseTopping(NewMushroomTopping(NewPepperoniTopping(pizza)))
	fmt.Printf("%s: $%.2f\n", supremePizza.GetDescription(), supremePizza.GetPrice())
	
	// Notification System example
	fmt.Println("\n--- Notification System Example ---")
	
	basicNotifier := &BasicNotifier{}
	basicNotifier.Send("System update available")
	
	smsNotifier := NewSMSNotifier(basicNotifier)
	smsNotifier.Send("System update available")
	
	multiNotifier := NewFacebookNotifier(NewSlackNotifier(NewEmailNotifier(basicNotifier)))
	multiNotifier.Send("System update available")
	
	// Data Stream example
	fmt.Println("\n--- Data Stream Example ---")
	
	fileDataSource := NewFileDataSource("data.txt")
	fileDataSource.WriteData("Sensitive data")
	fileDataSource.ReadData()
	
	encryptedDataSource := NewEncryptionDecorator(fileDataSource)
	encryptedDataSource.WriteData("Sensitive data")
	encryptedDataSource.ReadData()
	
	compressedEncryptedDataSource := NewCompressionDecorator(NewEncryptionDecorator(fileDataSource))
	compressedEncryptedDataSource.WriteData("Sensitive data")
	compressedEncryptedDataSource.ReadData()
	
	loggedDataSource := NewLoggingDecorator(NewEncryptionDecorator(fileDataSource))
	loggedDataSource.WriteData("Important data")
	loggedDataSource.ReadData()
	
	// Weapon Enhancement example
	fmt.Println("\n--- Weapon Enhancement Example ---")
	
	sword := &Sword{}
	fmt.Printf("%s: %d damage\n", sword.GetDescription(), sword.GetDamage())
	
	fireSword := NewFireEnchantment(sword)
	fmt.Printf("%s: %d damage\n", fireSword.GetDescription(), fireSword.GetDamage())
	
	superSword := NewSharpnessEnchantment(NewPoisonEnchantment(NewFireEnchantment(sword)))
	fmt.Printf("%s: %d damage\n", superSword.GetDescription(), superSword.GetDamage())
	
	// Car Features example
	fmt.Println("\n--- Car Features Example ---")
	
	basicCar := &BasicCar{}
	fmt.Printf("%s: $%.2f\n", basicCar.GetDescription(), basicCar.GetPrice())
	
	carWithGPS := NewGPSDecorator(basicCar)
	fmt.Printf("%s: $%.2f\n", carWithGPS.GetDescription(), carWithGPS.GetPrice())
	
	luxuryCar := NewPremiumAudioDecorator(NewLeatherSeatsDecorator(NewSunroofDecorator(NewGPSDecorator(basicCar))))
	fmt.Printf("%s: $%.2f\n", luxuryCar.GetDescription(), luxuryCar.GetPrice())
	
	fmt.Println("\nAll decorator patterns demonstrated successfully!")
}
