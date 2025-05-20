package tests

type Configuration struct {
	TestInt    int  `def:"17"`
	TestIntPtr *int `def:"5" validate:"min=0,max=5"`
	//TestStringEnv *string `env:"TEST_STRING_ENV" validate:"required"`
	//TestBoolYAML  *bool   `yaml:"test-bool" def:"true"`
}
