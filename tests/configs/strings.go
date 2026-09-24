package configs

type SpecialString string

type SpecialPtrString *string
type StringConfiguration struct {
	String                  string         `def:"default" validate:"min=1"`
	StringNoDef             string         `validate:"required"`
	StringPtr               *string        `def:"default" validate:"alpha"`
	StringPtrNoDef          *string        `validate:"required"`
	StringPtrZeroPtr        *string        `def:"*0*"`
	SpecialString           SpecialString  `def:"default"  validate:"alpha"`
	SpecialStringPtr        *SpecialString `def:"default" validate:"max=4"`
	SpecialStringPtrNoDef   *SpecialString
	SpecialStringPtrZeroPtr *SpecialString   `def:"*0*"`
	SpecialPtrString        SpecialPtrString `def:"default" validate:"numeric"`
	SpecialPtrStringNoDef   SpecialPtrString `validate:"required"`
	SpecialPtrStringZeroPtr SpecialPtrString `def:"*0*"`
}

type StringsConfiguration struct {
	Strings                  []string           `def:"de,fa,ult"`
	StringsNoDef             []string           `validate:"min=1"`
	StringsZeroDef           []string           `def:"*0*"`
	SpecialStrings           []SpecialString    `def:"de,fa,ult" validate:"max=2"`
	SpecialStringsNoDef      []SpecialString    `validate:"required"`
	SpecialStringsZeroDef    []SpecialString    `def:"*0*"`
	SpecialPtrStrings        []SpecialPtrString `def:"de,fa,ult"`
	SpecialPtrStringsNoDef   []SpecialPtrString
	SpecialPtrStringsZeroDef []SpecialPtrString `def:"*0*"`
	StringsPtr               []*string          `def:"de,fa,ult" validate:"unique"`
	StringsPtrNoDef          []*string
	StringsPtrZeroDef        []*string        `def:"*0*"`
	SpecialStringsPtr        []*SpecialString `def:"de,fa,ult"`
	SpecialStringsPtrNoDef   []*SpecialString
	SpecialStringsPtrZeroDef []*SpecialString `def:"*0*"`
}
