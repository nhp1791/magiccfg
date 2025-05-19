package tests

import "slices"

type TestEnumOne string

const (
	TestOne   TestEnumOne = "one"
	TestTwo   TestEnumOne = "two"
	TestThree TestEnumOne = "three"
)

var TestEnumOnes = []TestEnumOne{
	TestOne,
	TestTwo,
	TestThree,
}

func (t TestEnumOne) ValidateEnum() bool {
	return slices.Contains(TestEnumOnes, t)
}

type FileType string

const (
	FileTypeYAML FileType = "yaml"
	FileTypeJSON FileType = "json"
	FileTypeTOML FileType = "toml"
)
