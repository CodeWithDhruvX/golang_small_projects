package main

import "fmt"

// AbstractClass defines the template method
type DataProcessor interface {
	ProcessData()
	ReadData()
	ValidateData()
	TransformData()
	SaveData()
}

// BaseClass provides the template method implementation
type BaseDataProcessor struct{}

func (bdp *BaseDataProcessor) ProcessData() {
	fmt.Println("Starting data processing...")
	bdp.ReadData()
	if bdp.ValidateData() {
		bdp.TransformData()
		bdp.SaveData()
		fmt.Println("Data processing completed successfully!")
	} else {
		fmt.Println("Data processing failed due to validation errors!")
	}
}

// ConcreteClass1 - CSV Data Processor
type CSVDataProcessor struct {
	BaseDataProcessor
	data []string
}

func NewCSVDataProcessor() *CSVDataProcessor {
	return &CSVDataProcessor{
		data: []string{},
	}
}

func (cdp *CSVDataProcessor) ReadData() {
	fmt.Println("Reading data from CSV file...")
	cdp.data = []string{"John,30,Engineer", "Jane,25,Designer", "Bob,35,Manager"}
	fmt.Printf("Read %d records from CSV\n", len(cdp.data))
}

func (cdp *CSVDataProcessor) ValidateData() bool {
	fmt.Println("Validating CSV data...")
	for _, record := range cdp.data {
		if record == "" {
			fmt.Println("Validation failed: Empty record found")
			return false
		}
	}
	fmt.Println("CSV data validation passed")
	return true
}

func (cdp *CSVDataProcessor) TransformData() {
	fmt.Println("Transforming CSV data...")
	for i, record := range cdp.data {
		cdp.data[i] = "Transformed: " + record
	}
	fmt.Println("CSV data transformation completed")
}

func (cdp *CSVDataProcessor) SaveData() {
	fmt.Println("Saving transformed CSV data...")
	for _, record := range cdp.data {
		fmt.Printf("Saved: %s\n", record)
	}
}

// ConcreteClass2 - JSON Data Processor
type JSONDataProcessor struct {
	BaseDataProcessor
	data map[string]interface{}
}

func NewJSONDataProcessor() *JSONDataProcessor {
	return &JSONDataProcessor{
		data: make(map[string]interface{}),
	}
}

func (jdp *JSONDataProcessor) ReadData() {
	fmt.Println("Reading data from JSON file...")
	jdp.data = map[string]interface{}{
		"user": "Alice",
		"age":  28,
		"role": "Developer",
		"active": true,
	}
	fmt.Println("Read JSON data successfully")
}

func (jdp *JSONDataProcessor) ValidateData() bool {
	fmt.Println("Validating JSON data...")
	if _, exists := jdp.data["user"]; !exists {
		fmt.Println("Validation failed: Missing 'user' field")
		return false
	}
	if _, exists := jdp.data["age"]; !exists {
		fmt.Println("Validation failed: Missing 'age' field")
		return false
	}
	fmt.Println("JSON data validation passed")
	return true
}

func (jdp *JSONDataProcessor) TransformData() {
	fmt.Println("Transforming JSON data...")
	for key, value := range jdp.data {
		jdp.data[key] = fmt.Sprintf("processed_%v", value)
	}
	fmt.Println("JSON data transformation completed")
}

func (jdp *JSONDataProcessor) SaveData() {
	fmt.Println("Saving transformed JSON data...")
	for key, value := range jdp.data {
		fmt.Printf("Saved %s: %v\n", key, value)
	}
}

// Another example - Game Character Creation
type CharacterCreator interface {
	CreateCharacter()
	SelectRace()
	SelectClass()
	AssignAttributes()
	EquipItems()
}

type BaseCharacterCreator struct{}

func (bcc *BaseCharacterCreator) CreateCharacter() {
	fmt.Println("Creating new character...")
	bcc.SelectRace()
	bcc.SelectClass()
	bcc.AssignAttributes()
	bcc.EquipItems()
	fmt.Println("Character creation completed!")
}

type WarriorCreator struct {
	BaseCharacterCreator
	race string
	class string
}

func NewWarriorCreator() *WarriorCreator {
	return &WarriorCreator{}
}

func (wc *WarriorCreator) SelectRace() {
	wc.race = "Human"
	fmt.Printf("Selected race: %s\n", wc.race)
}

func (wc *WarriorCreator) SelectClass() {
	wc.class = "Warrior"
	fmt.Printf("Selected class: %s\n", wc.class)
}

func (wc *WarriorCreator) AssignAttributes() {
	fmt.Println("Assigning warrior attributes...")
	fmt.Println("- Strength: 18")
	fmt.Println("- Constitution: 16")
	fmt.Println("- Dexterity: 12")
	fmt.Println("- Intelligence: 8")
	fmt.Println("- Wisdom: 10")
	fmt.Println("- Charisma: 10")
}

func (wc *WarriorCreator) EquipItems() {
	fmt.Println("Equipping warrior items...")
	fmt.Println("- Longsword")
	fmt.Println("- Shield")
	fmt.Println("- Chain Mail Armor")
	fmt.Println("- Health Potion")
}

type MageCreator struct {
	BaseCharacterCreator
	race string
	class string
}

func NewMageCreator() *MageCreator {
	return &MageCreator{}
}

func (mc *MageCreator) SelectRace() {
	mc.race = "Elf"
	fmt.Printf("Selected race: %s\n", mc.race)
}

func (mc *MageCreator) SelectClass() {
	mc.class = "Mage"
	fmt.Printf("Selected class: %s\n", mc.class)
}

func (mc *MageCreator) AssignAttributes() {
	fmt.Println("Assigning mage attributes...")
	fmt.Println("- Strength: 8")
	fmt.Println("- Constitution: 10")
	fmt.Println("- Dexterity: 14")
	fmt.Println("- Intelligence: 18")
	fmt.Println("- Wisdom: 16")
	fmt.Println("- Charisma: 12")
}

func (mc *MageCreator) EquipItems() {
	fmt.Println("Equipping mage items...")
	fmt.Println("- Magic Staff")
	fmt.Println("- Robes")
	fmt.Println("- Spellbook")
	fmt.Println("- Mana Potion")
}

func main() {
	fmt.Println("=== Template Method Pattern Demo ===")
	
	// Data Processing example
	fmt.Println("\n--- Data Processing Example ---")
	
	fmt.Println("\nProcessing CSV Data:")
	csvProcessor := NewCSVDataProcessor()
	csvProcessor.ProcessData()
	
	fmt.Println("\nProcessing JSON Data:")
	jsonProcessor := NewJSONDataProcessor()
	jsonProcessor.ProcessData()
	
	// Game Character Creation example
	fmt.Println("\n--- Game Character Creation Example ---")
	
	fmt.Println("\nCreating Warrior Character:")
	warriorCreator := NewWarriorCreator()
	warriorCreator.CreateCharacter()
	
	fmt.Println("\nCreating Mage Character:")
	mageCreator := NewMageCreator()
	mageCreator.CreateCharacter()
}
