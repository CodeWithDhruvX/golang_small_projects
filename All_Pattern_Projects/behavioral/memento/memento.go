package main

import "fmt"

// Memento contains the state of an object
type Memento interface {
	GetState() string
	GetName() string
	GetDate() string
}

// ConcreteMemento stores the state
type ConcreteMemento struct {
	state string
	name  string
	date  string
}

func NewConcreteMemento(state string) *ConcreteMemento {
	return &ConcreteMemento{
		state: state,
		name:  fmt.Sprintf("Snapshot_%d", len(state)),
		date:  fmt.Sprintf("%d", len(state)), // Simplified date
	}
}

func (m *ConcreteMemento) GetState() string {
	return m.state
}

func (m *ConcreteMemento) GetName() string {
	return m.name
}

func (m *ConcreteMemento) GetDate() string {
	return m.date
}

// Originator creates a memento containing a snapshot of its current internal state
type Originator struct {
	state string
}

func NewOriginator(state string) *Originator {
	return &Originator{state: state}
}

func (o *Originator) SetState(state string) {
	o.state = state
	fmt.Printf("Originator: State changed to: %s\n", state)
}

func (o *Originator) GetState() string {
	return o.state
}

func (o *Originator) Save() Memento {
	return NewConcreteMemento(o.state)
}

func (o *Originator) Restore(memento Memento) {
	o.state = memento.GetState()
	fmt.Printf("Originator: State restored to: %s\n", o.state)
}

// Caretaker manages mementos but never operates on or examines their content
type Caretaker struct {
	mementos []Memento
	originator *Originator
}

func NewCaretaker(originator *Originator) *Caretaker {
	return &Caretaker{
		originator: originator,
		mementos:   make([]Memento, 0),
	}
}

func (c *Caretaker) Backup() {
	fmt.Println("\nCaretaker: Saving Originator's state...")
	c.mementos = append(c.mementos, c.originator.Save())
}

func (c *Caretaker) Undo() {
	if len(c.mementos) == 0 {
		fmt.Println("Caretaker: No states to restore")
		return
	}
	
	memento := c.mementos[len(c.mementos)-1]
	c.mementos = c.mementos[:len(c.mementos)-1]
	
	fmt.Printf("Caretaker: Restoring state to: %s\n", memento.GetName())
	c.originator.Restore(memento)
}

func (c *Caretaker) ShowHistory() {
	fmt.Println("Caretaker: Here's the list of mementos:")
	for _, memento := range c.mementos {
		fmt.Printf(memento.GetName() + "\n")
	}
}

// Text Editor example
type TextEditor struct {
	content string
}

func NewTextEditor() *TextEditor {
	return &TextEditor{content: ""}
}

func (te *TextEditor) Write(text string) {
	te.content += text
	fmt.Printf("Text Editor: Content: '%s'\n", te.content)
}

func (te *TextEditor) Save() Memento {
	return NewConcreteMemento(te.content)
}

func (te *TextEditor) Restore(memento Memento) {
	te.content = memento.GetState()
	fmt.Printf("Text Editor: Restored content: '%s'\n", te.content)
}

type TextEditorCaretaker struct {
	history []Memento
}

func NewTextEditorCaretaker() *TextEditorCaretaker {
	return &TextEditorCaretaker{history: make([]Memento, 0)}
}

func (tec *TextEditorCaretaker) Backup(editor *TextEditor) {
	tec.history = append(tec.history, editor.Save())
	fmt.Printf("Text Editor Caretaker: Saved state (total: %d)\n", len(tec.history))
}

func (tec *TextEditorCaretaker) Undo(editor *TextEditor) {
	if len(tec.history) == 0 {
		fmt.Println("Text Editor Caretaker: No backups to restore")
		return
	}
	
	lastState := tec.history[len(tec.history)-1]
	tec.history = tec.history[:len(tec.history)-1]
	editor.Restore(lastState)
	fmt.Printf("Text Editor Caretaker: Restored state (remaining: %d)\n", len(tec.history))
}

func main() {
	fmt.Println("=== Memento Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Originator Example ---")
	originator := NewOriginator("Super-duper-super-puper-super.")
	caretaker := NewCaretaker(originator)
	
	caretaker.Backup()
	originator.SetState("State #2")
	caretaker.Backup()
	originator.SetState("State #3")
	caretaker.Backup()
	originator.SetState("State #4")
	
	fmt.Println("\nClient: Now, let's rollback!\n")
	caretaker.Undo()
	caretaker.Undo()
	caretaker.Undo()
	
	// Text Editor example
	fmt.Println("\n--- Text Editor Example ---")
	editor := NewTextEditor()
	editorCaretaker := NewTextEditorCaretaker()
	
	editor.Write("Hello ")
	editorCaretaker.Backup(editor)
	
	editor.Write("World ")
	editorCaretaker.Backup(editor)
	
	editor.Write("in Go!")
	fmt.Println("\nClient: Let's undo the last change")
	editorCaretaker.Undo(editor)
	
	fmt.Println("\nClient: Let's undo again")
	editorCaretaker.Undo(editor)
}
