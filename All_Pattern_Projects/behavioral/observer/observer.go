package main

import "fmt"

// Observer interface defines the update method
type Observer interface {
	Update(subject Subject)
	GetName() string
}

// Subject interface defines methods for attaching/detaching observers and notifying
type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	Notify()
}

// ConcreteSubject maintains state and notifies observers on changes
type WeatherStation struct {
	observers    []Observer
	temperature  float64
	humidity     float64
	pressure     float64
}

func NewWeatherStation() *WeatherStation {
	return &WeatherStation{
		observers: make([]Observer, 0),
	}
}

func (ws *WeatherStation) Attach(observer Observer) {
	ws.observers = append(ws.observers, observer)
	fmt.Printf("WeatherStation: Attached observer: %s\n", observer.GetName())
}

func (ws *WeatherStation) Detach(observer Observer) {
	for i, obs := range ws.observers {
		if obs.GetName() == observer.GetName() {
			ws.observers = append(ws.observers[:i], ws.observers[i+1:]...)
			fmt.Printf("WeatherStation: Detached observer: %s\n", observer.GetName())
			return
		}
	}
}

func (ws *WeatherStation) Notify() {
	fmt.Println("WeatherStation: Notifying observers...")
	for _, observer := range ws.observers {
		observer.Update(ws)
	}
}

func (ws *WeatherStation) SetMeasurements(temperature, humidity, pressure float64) {
	ws.temperature = temperature
	ws.humidity = humidity
	ws.pressure = pressure
	fmt.Printf("WeatherStation: Measurements changed - Temp: %.1f°C, Humidity: %.1f%%, Pressure: %.1f hPa\n", 
		temperature, humidity, pressure)
	ws.Notify()
}

func (ws *WeatherStation) GetTemperature() float64 { return ws.temperature }
func (ws *WeatherStation) GetHumidity() float64    { return ws.humidity }
func (ws *WeatherStation) GetPressure() float64    { return ws.pressure }

// ConcreteObserver1 - Display device
type CurrentConditionsDisplay struct {
	name string
}

func NewCurrentConditionsDisplay() *CurrentConditionsDisplay {
	return &CurrentConditionsDisplay{name: "Current Conditions Display"}
}

func (ccd *CurrentConditionsDisplay) GetName() string {
	return ccd.name
}

func (ccd *CurrentConditionsDisplay) Update(subject Subject) {
	if ws, ok := subject.(*WeatherStation); ok {
		fmt.Printf("Current Conditions Display: %.1f°C and %.1f%% humidity\n", 
			ws.GetTemperature(), ws.GetHumidity())
	}
}

// ConcreteObserver2 - Statistics display
type StatisticsDisplay struct {
	name       string
	maxTemp    float64
	minTemp    float64
	tempSum    float64
	numReadings int
}

func NewStatisticsDisplay() *StatisticsDisplay {
	return &StatisticsDisplay{
		name:    "Statistics Display",
		maxTemp: -999.0,
		minTemp: 999.0,
	}
}

func (sd *StatisticsDisplay) GetName() string {
	return sd.name
}

func (sd *StatisticsDisplay) Update(subject Subject) {
	if ws, ok := subject.(*WeatherStation); ok {
		temp := ws.GetTemperature()
		sd.tempSum += temp
		sd.numReadings++
		
		if temp > sd.maxTemp {
			sd.maxTemp = temp
		}
		if temp < sd.minTemp {
			sd.minTemp = temp
		}
		
		avgTemp := sd.tempSum / float64(sd.numReadings)
		fmt.Printf("Statistics Display: Avg/Max/Min temperature = %.1f/%.1f/%.1f\n", 
			avgTemp, sd.maxTemp, sd.minTemp)
	}
}

// ConcreteObserver3 - Forecast display
type ForecastDisplay struct {
	name      string
	lastPressure float64
}

func NewForecastDisplay() *ForecastDisplay {
	return &ForecastDisplay{
		name:        "Forecast Display",
		lastPressure: 0.0,
	}
}

func (fd *ForecastDisplay) GetName() string {
	return fd.name
}

func (fd *ForecastDisplay) Update(subject Subject) {
	if ws, ok := subject.(*WeatherStation); ok {
		currentPressure := ws.GetPressure()
		
		if fd.lastPressure > 0 {
			if currentPressure > fd.lastPressure {
				fmt.Println("Forecast Display: Improving weather on the way!")
			} else if currentPressure == fd.lastPressure {
				fmt.Println("Forecast Display: More of the same")
			} else {
				fmt.Println("Forecast Display: Watch out for cooler, rainy weather")
			}
		}
		
		fd.lastPressure = currentPressure
	}
}

// Stock Market example
type StockMarket struct {
	observers []Observer
	stocks    map[string]float64
}

func NewStockMarket() *StockMarket {
	return &StockMarket{
		observers: make([]Observer, 0),
		stocks:    make(map[string]float64),
	}
}

func (sm *StockMarket) Attach(observer Observer) {
	sm.observers = append(sm.observers, observer)
}

func (sm *StockMarket) Detach(observer Observer) {
	for i, obs := range sm.observers {
		if obs.GetName() == observer.GetName() {
			sm.observers = append(sm.observers[:i], sm.observers[i+1:]...)
			return
		}
	}
}

func (sm *StockMarket) Notify() {
	for _, observer := range sm.observers {
		observer.Update(sm)
	}
}

func (sm *StockMarket) UpdateStock(symbol string, price float64) {
	sm.stocks[symbol] = price
	fmt.Printf("StockMarket: %s updated to $%.2f\n", symbol, price)
	sm.Notify()
}

func (sm *StockMarket) GetStockPrice(symbol string) float64 {
	return sm.stocks[symbol]
}

type StockTrader struct {
	name string
}

func NewStockTrader(name string) *StockTrader {
	return &StockTrader{name: name}
}

func (st *StockTrader) GetName() string {
	return st.name
}

func (st *StockTrader) Update(subject Subject) {
	if sm, ok := subject.(*StockMarket); ok {
		fmt.Printf("Stock Trader %s: Checking portfolio...\n", st.name)
		for symbol, price := range sm.stocks {
			fmt.Printf("  %s: $%.2f\n", symbol, price)
		}
	}
}

func main() {
	fmt.Println("=== Observer Pattern Demo ===")
	
	// Weather Station example
	fmt.Println("\n--- Weather Station Example ---")
	weatherStation := NewWeatherStation()
	
	currentDisplay := NewCurrentConditionsDisplay()
	statisticsDisplay := NewStatisticsDisplay()
	forecastDisplay := NewForecastDisplay()
	
	weatherStation.Attach(currentDisplay)
	weatherStation.Attach(statisticsDisplay)
	weatherStation.Attach(forecastDisplay)
	
	weatherStation.SetMeasurements(25.0, 65.0, 1013.1)
	weatherStation.SetMeasurements(27.0, 70.0, 1015.2)
	weatherStation.SetMeasurements(23.0, 60.0, 1010.5)
	
	// Stock Market example
	fmt.Println("\n--- Stock Market Example ---")
	stockMarket := NewStockMarket()
	
	trader1 := NewStockTrader("Alice")
	trader2 := NewStockTrader("Bob")
	
	stockMarket.Attach(trader1)
	stockMarket.Attach(trader2)
	
	stockMarket.UpdateStock("AAPL", 150.25)
	stockMarket.UpdateStock("GOOGL", 2750.80)
	stockMarket.UpdateStock("MSFT", 305.15)
}
