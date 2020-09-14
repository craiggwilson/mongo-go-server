package tree

// New makes a new Tree.
func New() *Tree {
	return &Tree{}
}

type Tree struct {
	AttributeContainer

	Filename string

	Commands []*Command
	Services []*Service
	Structs  []*Struct
}

func (t *Tree) AddCommand(cmd *Command) {
	t.Commands = append(t.Commands, cmd)
}

func (t *Tree) AddService(s *Service) {
	t.Services = append(t.Services, s)
}

func (t *Tree) AddStruct(s *Struct) {
	t.Structs = append(t.Structs, s)
}

func NewCommand(name string) *Command {
	return &Command{
		Name: name,
	}
}

type Command struct {
	AttributeContainer

	Name     string
	Request  *Struct
	Response *Struct
}

func NewService(name string) *Service {
	return &Service{
		Name: name,
	}
}

type Service struct {
	AttributeContainer

	Name        string
	CommandRefs []string
}

func (s *Service) AddCommandRef(name string) {
	s.CommandRefs = append(s.CommandRefs, name)
}

func NewStruct(name string) *Struct {
	return &Struct{
		Name: name,
	}
}

type Struct struct {
	AttributeContainer

	Name   string
	Fields []*Field
}

func (s *Struct) AddField(f *Field) {
	s.Fields = append(s.Fields, f)
}

func NewField(name string) *Field {
	return &Field{
		Name: name,
	}
}

type Field struct {
	AttributeContainer

	Name    string
	TypeRef string
}

type AttributeContainer []*Attribute

func (ac *AttributeContainer) AddAttribute(attr *Attribute) {
	*ac = append(*ac, attr)
}

func (ac AttributeContainer) Attribute(lang, name string) string {
	for _, attr := range ac {
		if attr.matches(lang, name) {
			return attr.Value
		}
	}

	return ""
}

func (ac AttributeContainer) AttributeOrDefault(lang, name string, def string) string {
	for _, attr := range ac {
		if attr.matches(lang, name) {
			return attr.Value
		}
	}

	return def
}

func (ac AttributeContainer) Attributes(lang, name string) []string {
	var values []string
	for _, attr := range ac {
		if attr.matches(lang, name) {
			values = append(values, attr.Value)
		}
	}

	return values
}

func (ac AttributeContainer) AttributesOrDefault(lang, name string, def string) []string {
	var values []string
	for _, attr := range ac {
		if attr.matches(lang, name) {
			values = append(values, attr.Value)
		}
	}

	if len(values) == 0 {
		values = append(values, def)
	}

	return values
}

func NewAttribute(lang, name, value string) *Attribute {
	return &Attribute{
		Lang:  lang,
		Name:  name,
		Value: value,
	}
}

type Attribute struct {
	Lang  string
	Name  string
	Value string
}

func (a *Attribute) matches(lang, name string) bool {
	if a.Name != name {
		return false
	}

	return lang == "" || a.Lang == lang
}
