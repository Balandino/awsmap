package engine

type BorderChild struct {
	Position string `yaml:"Position,omitempty"`
	Resource string `yaml:"Resource,omitempty"`
}

type IconFill struct {
	Type  string `yaml:"Type,omitempty"`
	Color string `yaml:"Color,omitempty"`
}

type Resource struct {
	Type           string        `yaml:"Type"`
	Direction      string        `yaml:"Direction,omitempty"`
	Preset         string        `yaml:"Preset,omitempty"`
	Title          string        `yaml:"Title,omitempty"`
	Children       []string      `yaml:"Children,omitempty"`
	BorderChildren []BorderChild `yaml:"BorderChildren,omitempty"`
	FillColor      string        `yaml:"FillColor,omitempty"`
	IconFill       *IconFill     `yaml:"IconFill,omitempty"`
}

type DefinitionFile struct {
	Type      string `yaml:"Type"`
	LocalFile string `yaml:"LocalFile,omitempty"`
}

type DiagramRoot struct {
	Diagram DiagramResources `yaml:"Diagram"`
}

type Link struct {
	Source         string  `yaml:"Source"`
	SourcePosition string  `yaml:"SourcePosition"`
	Target         string  `yaml:"Target"`
	TargetPosition string  `yaml:"TargetPosition"`
	LineColor      string  `yaml:"LineColor"`
	Labels         *Labels `yaml:"Labels,omitempty"`
}

// Labels holds label details for various position attachments
type Labels struct {
	SourceLeft  *LabelDetail `yaml:"SourceLeft,omitempty"`
	SourceRight *LabelDetail `yaml:"SourceRight,omitempty"`
	TargetLeft  *LabelDetail `yaml:"TargetLeft,omitempty"`
	TargetRight *LabelDetail `yaml:"TargetRight,omitempty"`
}

// LabelDetail holds the specific text properties of the label
type LabelDetail struct {
	Title string `yaml:"Title"`
}

type DiagramResources struct {
	DefinitionFiles []DefinitionFile     `yaml:"DefinitionFiles"`
	Resources       map[string]*Resource `yaml:"Resources"`
	Links           []Link               `yaml:"Links,omitempty"`
}

type YamlResults struct {
	Account  string
	YamlTree DiagramRoot `yaml:"Diagram"`
}
