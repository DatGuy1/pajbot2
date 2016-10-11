package modules

type baseSetting struct {
	Valid bool
}

// IntSetting xD
type IntSetting struct {
	baseSetting

	Label        string
	DefaultValue int
	Value        *int

	// Constraints
	MinValue int
	MaxValue int
}

// Int returns the int value of an int setting
// If the Value is nil then we will just return the DefaultValue int
func (s *IntSetting) Int() int {
	if s.Value == nil {
		return s.DefaultValue
	}

	return *s.Value
}

// StringSetting xD
type StringSetting struct {
	baseSetting

	// Label is used to describe the setting
	Label string

	DefaultValue string
	Value        *string

	// Constraints
	MinLength int
	MaxLength int
}

// String returns the string value of a string setting
// If the Value is nil then we will just return the DefaultValue string
func (s *StringSetting) String() string {
	if s.Value == nil {
		return s.DefaultValue
	}

	return *s.Value
}

// BoolSetting xD
type BoolSetting struct {
	baseSetting

	// Label is used to describe the setting
	Label string

	DefaultValue bool
	Value        *bool
}

// Bool returns the string value of a string setting
// If the Value is nil then we will just return the DefaultValue string
func (s *BoolSetting) Bool() bool {
	if s.Value == nil {
		return s.DefaultValue
	}

	return *s.Value
}

// OptionSetting xD
type OptionSetting struct {
	baseSetting

	// Label is used to describe the setting
	Label string

	DefaultValue int
	Value        *int
	Options      []string
}

// Int returns the index of the currently selected index
// If the Index is nil then we will just return the DefaultValue string
func (s *OptionSetting) Int() int {
	if s.Value == nil {
		return s.DefaultValue
	}

	return *s.Value
}

// String returns the index of the currently selected index
// If the Index is nil then we will just return the DefaultValue string
func (s *OptionSetting) String() string {
	index := s.Int()

	// XXX(pajlada): Make bounds check here
	return s.Options[index]
}
