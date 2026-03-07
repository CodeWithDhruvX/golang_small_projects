package main

import "fmt"

// Command interface defines the execute method
type Command interface {
	Execute()
}

// Receiver contains the business logic
type Light struct {
	isOn bool
}

func (l *Light) TurnOn() {
	l.isOn = true
	fmt.Println("Light is ON")
}

func (l *Light) TurnOff() {
	l.isOn = false
	fmt.Println("Light is OFF")
}

// ConcreteCommand1 for turning on the light
type LightOnCommand struct {
	light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
	return &LightOnCommand{light: light}
}

func (c *LightOnCommand) Execute() {
	c.light.TurnOn()
}

// ConcreteCommand2 for turning off the light
type LightOffCommand struct {
	light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
	return &LightOffCommand{light: light}
}

func (c *LightOffCommand) Execute() {
	c.light.TurnOff()
}

// Invoker triggers the command
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(command Command) {
	r.command = command
}

func (r *RemoteControl) PressButton() {
	if r.command != nil {
		r.command.Execute()
	}
}

// Another receiver example
type Stereo struct {
	isOn bool
	volume int
}

func (s *Stereo) TurnOn() {
	s.isOn = true
	fmt.Println("Stereo is ON")
}

func (s *Stereo) TurnOff() {
	s.isOn = false
	fmt.Println("Stereo is OFF")
}

func (s *Stereo) SetVolume(volume int) {
	s.volume = volume
	fmt.Printf("Stereo volume set to %d\n", volume)
}

// Stereo commands
type StereoOnWithCDCommand struct {
	stereo *Stereo
}

func NewStereoOnWithCDCommand(stereo *Stereo) *StereoOnWithCDCommand {
	return &StereoOnWithCDCommand{stereo: stereo}
}

func (c *StereoOnWithCDCommand) Execute() {
	c.stereo.TurnOn()
	c.stereo.SetVolume(11)
}

type StereoOffCommand struct {
	stereo *Stereo
}

func NewStereoOffCommand(stereo *Stereo) *StereoOffCommand {
	return &StereoOffCommand{stereo: stereo}
}

func (c *StereoOffCommand) Execute() {
	c.stereo.TurnOff()
}

func main() {
	fmt.Println("=== Command Pattern Demo ===")
	
	// Create receivers
	light := &Light{}
	stereo := &Stereo{}
	
	// Create commands
	lightOn := NewLightOnCommand(light)
	lightOff := NewLightOffCommand(light)
	stereoOn := NewStereoOnWithCDCommand(stereo)
	stereoOff := NewStereoOffCommand(stereo)
	
	// Create invoker
	remote := &RemoteControl{}
	
	// Test light commands
	fmt.Println("\n--- Light Control ---")
	remote.SetCommand(lightOn)
	remote.PressButton()
	
	remote.SetCommand(lightOff)
	remote.PressButton()
	
	// Test stereo commands
	fmt.Println("\n--- Stereo Control ---")
	remote.SetCommand(stereoOn)
	remote.PressButton()
	
	remote.SetCommand(stereoOff)
	remote.PressButton()
}
