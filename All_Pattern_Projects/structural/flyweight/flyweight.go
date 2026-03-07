package main

import "fmt"

// Flyweight Pattern

// Flyweight interface
type Flyweight interface {
	Operation(extrinsicState string)
}

// Concrete Flyweight
type ConcreteFlyweight struct {
	intrinsicState string
}

func NewConcreteFlyweight(intrinsicState string) *ConcreteFlyweight {
	return &ConcreteFlyweight{intrinsicState: intrinsicState}
}

func (cf *ConcreteFlyweight) Operation(extrinsicState string) {
	fmt.Printf("ConcreteFlyweight: Intrinsic = %s, Extrinsic = %s\n", cf.intrinsicState, extrinsicState)
}

// Unshared Concrete Flyweight
type UnsharedConcreteFlyweight struct {
	allState string
}

func NewUnsharedConcreteFlyweight(allState string) *UnsharedConcreteFlyweight {
	return &UnsharedConcreteFlyweight{allState: allState}
}

func (ucf *UnsharedConcreteFlyweight) Operation(extrinsicState string) {
	fmt.Printf("UnsharedConcreteFlyweight: All state = %s\n", ucf.allState)
}

// Flyweight Factory
type FlyweightFactory struct {
	flyweights map[string]Flyweight
}

func NewFlyweightFactory() *FlyweightFactory {
	return &FlyweightFactory{
		flyweights: make(map[string]Flyweight),
	}
}

func (ff *FlyweightFactory) GetFlyweight(key string) Flyweight {
	if flyweight, exists := ff.flyweights[key]; exists {
		return flyweight
	}
	
	// Create new flyweight
	flyweight := NewConcreteFlyweight(key)
	ff.flyweights[key] = flyweight
	fmt.Printf("Creating new flyweight with key: %s\n", key)
	return flyweight
}

func (ff *FlyweightFactory) ListFlyweights() {
	fmt.Printf("Flyweight Factory has %d flyweights:\n", len(ff.flyweights))
	for key, flyweight := range ff.flyweights {
		fmt.Printf("  %s: %T\n", key, flyweight)
	}
}

// Text Editor Example
type CharacterFlyweight struct {
	character rune
	font      string
	size      int
	color     string
}

func NewCharacterFlyweight(character rune, font string, size int, color string) *CharacterFlyweight {
	return &CharacterFlyweight{
		character: character,
		font:      font,
		size:      size,
		color:     color,
	}
}

func (cf *CharacterFlyweight) Display(position int) {
	fmt.Printf("'%c' at position %d (Font: %s, Size: %d, Color: %s)\n", 
		cf.character, position, cf.font, cf.size, cf.color)
}

type CharacterFactory struct {
	flyweights map[string]*CharacterFlyweight
}

func NewCharacterFactory() *CharacterFactory {
	return &CharacterFactory{
		flyweights: make(map[string]*CharacterFlyweight),
	}
}

func (cf *CharacterFactory) GetCharacter(character rune, font string, size int, color string) *CharacterFlyweight {
	key := fmt.Sprintf("%c-%s-%d-%s", character, font, size, color)
	
	if flyweight, exists := cf.flyweights[key]; exists {
		return flyweight
	}
	
	flyweight := NewCharacterFlyweight(character, font, size, color)
	cf.flyweights[key] = flyweight
	fmt.Printf("Creating new character flyweight: %s\n", key)
	return flyweight
}

func (cf *CharacterFactory) GetFlyweightCount() int {
	return len(cf.flyweights)
}

type CharacterContext struct {
	character *CharacterFlyweight
	position  int
}

func NewCharacterContext(character *CharacterFlyweight, position int) *CharacterContext {
	return &CharacterContext{
		character: character,
		position:  position,
	}
}

func (cc *CharacterContext) Display() {
	cc.character.Display(cc.position)
}

// Tree Rendering Example
type TreeType struct {
	name  string
	color string
	texture string
}

func NewTreeType(name, color, texture string) *TreeType {
	return &TreeType{
		name:    name,
		color:   color,
		texture: texture,
	}
}

func (tt *TreeType) Draw(x, y int) {
	fmt.Printf("Drawing %s tree at (%d, %d) with color %s and texture %s\n", 
		tt.name, x, y, tt.color, tt.texture)
}

type TreeFactory struct {
	treeTypes map[string]*TreeType
}

func NewTreeFactory() *TreeFactory {
	return &TreeFactory{
		treeTypes: make(map[string]*TreeType),
	}
}

func (tf *TreeFactory) GetTreeType(name, color, texture string) *TreeType {
	key := fmt.Sprintf("%s-%s-%s", name, color, texture)
	
	if treeType, exists := tf.treeTypes[key]; exists {
		return treeType
	}
	
	treeType := NewTreeType(name, color, texture)
	tf.treeTypes[key] = treeType
	fmt.Printf("Creating new tree type: %s\n", key)
	return treeType
}

func (tf *TreeFactory) GetTreeTypeCount() int {
	return len(tf.treeTypes)
}

type Tree struct {
	x, y   int
	treeType *TreeType
}

func NewTree(x, y int, treeType *TreeType) *Tree {
	return &Tree{x: x, y: y, treeType: treeType}
}

func (t *Tree) Draw() {
	t.treeType.Draw(t.x, t.y)
}

type Forest struct {
	trees []*Tree
	factory *TreeFactory
}

func NewForest() *Forest {
	return &Forest{
		trees:   make([]*Tree, 0),
		factory: NewTreeFactory(),
	}
}

func (f *Forest) PlantTree(x, y int, name, color, texture string) {
	treeType := f.factory.GetTreeType(name, color, texture)
	tree := NewTree(x, y, treeType)
	f.trees = append(f.trees, tree)
}

func (f *Forest) Draw() {
	for _, tree := range f.trees {
		tree.Draw()
	}
}

func (f *Forest) GetTreeCount() int {
	return len(f.trees)
}

func (f *Forest) GetTreeTypeCount() int {
	return f.factory.GetTreeTypeCount()
}

// Game Objects Example
type Mesh struct {
	vertices []string
}

func NewMesh(vertices []string) *Mesh {
	return &Mesh{vertices: vertices}
}

type Texture struct {
	image string
}

func NewTexture(image string) *Texture {
	return &Texture{image: image}
}

type GameObjectFlyweight struct {
	mesh    *Mesh
	texture *Texture
	shader  string
}

func NewGameObjectFlyweight(mesh *Mesh, texture *Texture, shader string) *GameObjectFlyweight {
	return &GameObjectFlyweight{
		mesh:    mesh,
		texture: texture,
		shader:  shader,
	}
}

func (gof *GameObjectFlyweight) Render(position string, rotation float64, scale float64) {
	fmt.Printf("Rendering object at %s (rotation: %.1f, scale: %.1f)\n", position, rotation, scale)
	fmt.Printf("  Mesh: %v, Texture: %s, Shader: %s\n", gof.mesh.vertices, gof.texture.image, gof.shader)
}

type GameObjectFactory struct {
	flyweights map[string]*GameObjectFlyweight
}

func NewGameObjectFactory() *GameObjectFactory {
	return &GameObjectFactory{
		flyweights: make(map[string]*GameObjectFlyweight),
	}
}

func (gof *GameObjectFactory) GetGameObject(meshType, textureType, shader string) *GameObjectFlyweight {
	key := fmt.Sprintf("%s-%s-%s", meshType, textureType, shader)
	
	if flyweight, exists := gof.flyweights[key]; exists {
		return flyweight
	}
	
	// Create mesh and texture (in real app, these would be loaded from files)
	var mesh *Mesh
	var texture *Texture
	
	switch meshType {
	case "cube":
		mesh = NewMesh([]string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8"})
	case "sphere":
		mesh = NewMesh([]string{"sphere_vertices"})
	case "plane":
		mesh = NewMesh([]string{"plane_vertices"})
	default:
		mesh = NewMesh([]string{"default_vertices"})
	}
	
	texture = NewTexture(textureType)
	
	flyweight := NewGameObjectFlyweight(mesh, texture, shader)
	gof.flyweights[key] = flyweight
	fmt.Printf("Creating new game object flyweight: %s\n", key)
	return flyweight
}

func (gof *GameObjectFactory) GetFlyweightCount() int {
	return len(gof.flyweights)
}

type GameObject struct {
	position string
	rotation float64
	scale    float64
	flyweight *GameObjectFlyweight
}

func NewGameObject(position string, rotation, scale float64, flyweight *GameObjectFlyweight) *GameObject {
	return &GameObject{
		position:  position,
		rotation:  rotation,
		scale:     scale,
		flyweight: flyweight,
	}
}

func (go *GameObject) Render() {
	go.flyweight.Render(go.position, go.rotation, go.scale)
}

type GameWorld struct {
	objects []*GameObject
	factory *GameObjectFactory
}

func NewGameWorld() *GameWorld {
	return &GameWorld{
		objects: make([]*GameObject, 0),
		factory: NewGameObjectFactory(),
	}
}

func (gw *GameWorld) AddGameObject(position string, rotation, scale float64, meshType, textureType, shader string) {
	flyweight := gw.factory.GetGameObject(meshType, textureType, shader)
	gameObject := NewGameObject(position, rotation, scale, flyweight)
	gw.objects = append(gw.objects, gameObject)
}

func (gw *GameWorld) Render() {
	for _, obj := range gw.objects {
		obj.Render()
	}
}

func (gw *GameWorld) GetObjectCount() int {
	return len(gw.objects)
}

func (gw *GameWorld) GetFlyweightCount() int {
	return gw.factory.GetFlyweightCount()
}

// Music Notes Example
type MusicNote struct {
	pitch    string
	duration string
	volume   int
}

func NewMusicNote(pitch, duration string, volume int) *MusicNote {
	return &MusicNote{
		pitch:    pitch,
		duration: duration,
		volume:   volume,
	}
}

func (mn *MusicNote) Play(timestamp float64) {
	fmt.Printf("Playing %s note for %s at volume %d (time: %.1f)\n", 
		mn.pitch, mn.duration, mn.volume, timestamp)
}

type MusicNoteFactory struct {
	notes map[string]*MusicNote
}

func NewMusicNoteFactory() *MusicNoteFactory {
	return &MusicNoteFactory{
		notes: make(map[string]*MusicNote),
	}
}

func (mnf *MusicNoteFactory) GetNote(pitch, duration string, volume int) *MusicNote {
	key := fmt.Sprintf("%s-%s-%d", pitch, duration, volume)
	
	if note, exists := mnf.notes[key]; exists {
		return note
	}
	
	note := NewMusicNote(pitch, duration, volume)
	mnf.notes[key] = note
	fmt.Printf("Creating new music note: %s\n", key)
	return note
}

func (mnf *MusicNoteFactory) GetNoteCount() int {
	return len(mnf.notes)
}

type MusicScore struct {
	notes []*MusicNote
	timestamps []float64
	factory *MusicNoteFactory
}

func NewMusicScore() *MusicScore {
	return &MusicScore{
		notes:      make([]*MusicNote, 0),
		timestamps: make([]float64, 0),
		factory:    NewMusicNoteFactory(),
	}
}

func (ms *MusicScore) AddNote(pitch, duration string, volume int, timestamp float64) {
	note := ms.factory.GetNote(pitch, duration, volume)
	ms.notes = append(ms.notes, note)
	ms.timestamps = append(ms.timestamps, timestamp)
}

func (ms *MusicScore) Play() {
	for i, note := range ms.notes {
		note.Play(ms.timestamps[i])
	}
}

func (ms *MusicScore) GetNoteCount() int {
	return len(ms.notes)
}

func (ms *MusicScore) GetUniqueNoteCount() int {
	return ms.factory.GetNoteCount()
}

func main() {
	fmt.Println("=== Flyweight Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Flyweight Example ---")
	
	factory := NewFlyweightFactory()
	
	flyweight1 := factory.GetFlyweight("A")
	flyweight1.Operation("X")
	
	flyweight2 := factory.GetFlyweight("B")
	flyweight2.Operation("Y")
	
	flyweight3 := factory.GetFlyweight("A") // Should return existing flyweight
	flyweight3.Operation("Z")
	
	unshared := NewUnsharedConcreteFlyweight("Unshared")
	unshared.Operation("W")
	
	factory.ListFlyweights()
	
	// Text Editor example
	fmt.Println("\n--- Text Editor Example ---")
	
	characterFactory := NewCharacterFactory()
	
	// Create characters with shared properties
	char1 := characterFactory.GetCharacter('H', "Arial", 12, "black")
	char2 := characterFactory.GetCharacter('e', "Arial", 12, "black")
	char3 := characterFactory.GetCharacter('l', "Arial", 12, "black")
	char4 := characterFactory.GetCharacter('l', "Arial", 12, "black")
	char5 := characterFactory.GetCharacter('o', "Arial", 12, "black")
	
	// Some characters with different properties
	char6 := characterFactory.GetCharacter('W', "Arial", 16, "blue")
	char7 := characterFactory.GetCharacter('o', "Arial", 16, "blue")
	char8 := characterFactory.GetCharacter('r', "Arial", 16, "blue")
	char9 := characterFactory.GetCharacter('l', "Arial", 16, "blue")
	char10 := characterFactory.GetCharacter('d', "Arial", 16, "blue")
	
	// Create character contexts with positions
	contexts := []*CharacterContext{
		NewCharacterContext(char1, 0),
		NewCharacterContext(char2, 1),
		NewCharacterContext(char3, 2),
		NewCharacterContext(char4, 3),
		NewCharacterContext(char5, 4),
		NewCharacterContext(char6, 5),
		NewCharacterContext(char7, 6),
		NewCharacterContext(char8, 7),
		NewCharacterContext(char9, 8),
		NewCharacterContext(char10, 9),
	}
	
	fmt.Println("\nDisplaying text:")
	for _, context := range contexts {
		context.Display()
	}
	
	fmt.Printf("\nTotal characters: %d, Unique flyweights: %d\n", len(contexts), characterFactory.GetFlyweightCount())
	
	// Tree Rendering example
	fmt.Println("\n--- Tree Rendering Example ---")
	
	forest := NewForest()
	
	// Plant many trees (only a few unique tree types)
	forest.PlantTree(10, 20, "Oak", "Green", "Rough")
	forest.PlantTree(15, 30, "Oak", "Green", "Rough")
	forest.PlantTree(25, 40, "Oak", "Green", "Rough")
	forest.PlantTree(35, 50, "Oak", "Green", "Rough")
	
	forest.PlantTree(20, 25, "Pine", "Dark Green", "Smooth")
	forest.PlantTree(30, 35, "Pine", "Dark Green", "Smooth")
	forest.PlantTree(40, 45, "Pine", "Dark Green", "Smooth")
	
	forest.PlantTree(45, 55, "Maple", "Red", "Medium")
	forest.PlantTree(50, 60, "Maple", "Red", "Medium")
	
	fmt.Println("\nDrawing forest:")
	forest.Draw()
	
	fmt.Printf("\nTotal trees: %d, Unique tree types: %d\n", forest.GetTreeCount(), forest.GetTreeTypeCount())
	
	// Game Objects example
	fmt.Println("\n--- Game Objects Example ---")
	
	gameWorld := NewGameWorld()
	
	// Add many game objects (only a few unique flyweights)
	gameWorld.AddGameObject("0,0,0", 0.0, 1.0, "cube", "metal", "phong")
	gameWorld.AddGameObject("10,0,0", 45.0, 1.0, "cube", "metal", "phong")
	gameWorld.AddGameObject("20,0,0", 90.0, 1.0, "cube", "metal", "phong")
	
	gameWorld.AddGameObject("0,10,0", 0.0, 1.5, "sphere", "wood", "lambert")
	gameWorld.AddGameObject("10,10,0", 180.0, 1.5, "sphere", "wood", "lambert")
	
	gameWorld.AddGameObject("0,0,10", 0.0, 2.0, "cube", "stone", "phong")
	gameWorld.AddGameObject("10,0,10", 270.0, 2.0, "cube", "stone", "phong")
	
	gameWorld.AddGameObject("0,10,10", 0.0, 0.5, "sphere", "metal", "phong")
	gameWorld.AddGameObject("10,10,10", 0.0, 0.5, "sphere", "metal", "phong")
	
	fmt.Println("\nRendering game world:")
	gameWorld.Render()
	
	fmt.Printf("\nTotal objects: %d, Unique flyweights: %d\n", gameWorld.GetObjectCount(), gameWorld.GetFlyweightCount())
	
	// Music Notes example
	fmt.Println("\n--- Music Notes Example ---")
	
	score := NewMusicScore()
	
	// Add notes to the score (many notes, but few unique types)
	score.AddNote("C", "quarter", 80, 0.0)
	score.AddNote("D", "quarter", 80, 0.5)
	score.AddNote("E", "quarter", 80, 1.0)
	score.AddNote("F", "quarter", 80, 1.5)
	score.AddNote("G", "quarter", 80, 2.0)
	
	score.AddNote("C", "half", 100, 2.5)
	score.AddNote("G", "half", 100, 3.5)
	
	score.AddNote("C", "quarter", 80, 4.5)
	score.AddNote("D", "quarter", 80, 5.0)
	score.AddNote("E", "quarter", 80, 5.5)
	score.AddNote("F", "quarter", 80, 6.0)
	score.AddNote("G", "quarter", 80, 6.5)
	
	fmt.Println("\nPlaying music score:")
	score.Play()
	
	fmt.Printf("\nTotal notes: %d, Unique note types: %d\n", score.GetNoteCount(), score.GetUniqueNoteCount())
	
	fmt.Println("\nAll flyweight patterns demonstrated successfully!")
}
