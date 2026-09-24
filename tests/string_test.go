package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"magiccfg"
	"magiccfg/tests/configs"
)

func TestStringDefaults(t *testing.T) {
	c, err := magiccfg.NewMagicConfig(new(configs.StringConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ApplyDefaults().
		Validate().
		Config()
	assert.Len(t, errs, 5)

	assert.ErrorContains(t, errs[0], "StringNoDef")
	assert.ErrorContains(t, errs[0], "required")

	assert.ErrorContains(t, errs[1], "StringPtrNoDef")
	assert.ErrorContains(t, errs[1], "required")

	assert.ErrorContains(t, errs[2], "SpecialStringPtr")
	assert.ErrorContains(t, errs[2], "max")

	assert.ErrorContains(t, errs[3], "SpecialPtrString")
	assert.ErrorContains(t, errs[3], "numeric")

	assert.ErrorContains(t, errs[4], "SpecialPtrStringNoDef")
	assert.ErrorContains(t, errs[4], "required")
	assert.Equal(t, cfg.String, "default")
	assert.Equal(t, cfg.StringNoDef, "")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "default")
	assert.Nil(t, cfg.StringPtrNoDef)
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("default"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("default"))
	assert.Nil(t, cfg.SpecialStringPtrNoDef)
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString(""))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "default")
	assert.Nil(t, cfg.SpecialPtrStringNoDef)
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "")
}

func TestStringEnv(t *testing.T) {
	t.Setenv("STRING", "ENV")
	t.Setenv("STRING_NO_DEF", "ENV")
	t.Setenv("STRING_PTR", "ENV")
	t.Setenv("STRING_PTR_NO_DEF", "ENV")
	t.Setenv("STRING_PTR_ZERO_PTR", "ENV")
	t.Setenv("SPECIAL_STRING", "ENV")
	t.Setenv("SPECIAL_STRING_PTR", "ENV")
	t.Setenv("SPECIAL_STRING_PTR_NO_DEF", "ENV")
	t.Setenv("SPECIAL_STRING_PTR_ZERO_PTR", "ENV")
	t.Setenv("SPECIAL_PTR_STRING", "ENV")
	t.Setenv("SPECIAL_PTR_STRING_NO_DEF", "ENV")
	t.Setenv("SPECIAL_PTR_STRING_ZERO_PTR", "ENV")

	c, err := magiccfg.NewMagicConfig(new(configs.StringConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ParseEnv().
		Config()
	assert.Empty(t, errs)

	assert.Equal(t, cfg.String, "ENV")
	assert.Equal(t, cfg.StringNoDef, "ENV")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "ENV")
	assert.NotNil(t, cfg.StringPtrNoDef)
	assert.Equal(t, *cfg.StringPtrNoDef, "ENV")
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "ENV")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("ENV"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("ENV"))
	assert.NotNil(t, cfg.SpecialStringPtrNoDef)
	assert.Equal(t, *cfg.SpecialStringPtrNoDef, configs.SpecialString("ENV"))
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString("ENV"))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "ENV")
	assert.NotNil(t, cfg.SpecialPtrStringNoDef)
	assert.Equal(t, *cfg.SpecialPtrStringNoDef, "ENV")
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "ENV")
}

func TestStringFile(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/full_string.yaml")

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringConfiguration),
		&magiccfg.Options{
			ConfigFileEnvName: "CONFIG_FILE",
		},
	)
	assert.NoError(t, err)

	cfg, errs := c.ParseFiles().
		Config()
	assert.Empty(t, errs)

	assert.Equal(t, cfg.String, "file")
	assert.Equal(t, cfg.StringNoDef, "file")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "file")
	assert.NotNil(t, cfg.StringPtrNoDef)
	assert.Equal(t, *cfg.StringPtrNoDef, "file")
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "file")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("file"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("file"))
	assert.NotNil(t, cfg.SpecialStringPtrNoDef)
	assert.Equal(t, *cfg.SpecialStringPtrNoDef, configs.SpecialString("file"))
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString("file"))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "file")
	assert.NotNil(t, cfg.SpecialPtrStringNoDef)
	assert.Equal(t, *cfg.SpecialPtrStringNoDef, "file")
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "file")
}

func TestStringFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringFlags",
		"--string=flag",
		"--string-no-def=flag",
		"--string-ptr=flag",
		"--string-ptr-no-def=flag",
		"--string-ptr-zero-ptr=flag",
		"--special-string=flag",
		"--special-string-ptr=flag",
		"--special-string-ptr-no-def=flag",
		"--special-string-ptr-zero-ptr=flag",
		"--special-ptr-string=flag",
		"--special-ptr-string-no-def=flag",
		"--special-ptr-string-zero-ptr=flag",
	}

	c, err := magiccfg.NewMagicConfig(new(configs.StringConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ParseFlags().
		Config()
	assert.Empty(t, errs)

	assert.Equal(t, cfg.String, "flag")
	assert.Equal(t, cfg.StringNoDef, "flag")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "flag")
	assert.NotNil(t, cfg.StringPtrNoDef)
	assert.Equal(t, *cfg.StringPtrNoDef, "flag")
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "flag")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialStringPtrNoDef)
	assert.Equal(t, *cfg.SpecialStringPtrNoDef, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "flag")
	assert.NotNil(t, cfg.SpecialPtrStringNoDef)
	assert.Equal(t, *cfg.SpecialPtrStringNoDef, "flag")
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "flag")
}

func TestStringDEFC(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/partial_string.yaml")
	t.Setenv("STRING_PTR", "ENV")
	t.Setenv("STRING_PTR_NO_DEF", "ENV")
	t.Setenv("SPECIAL_STRING", "ENV")
	t.Setenv("SPECIAL_STRING_PTR_ZERO_PTR", "ENV")
	t.Setenv("SPECIAL_PTR_STRING", "2")
	t.Setenv("SPECIAL_PTR_STRING_NO_DEF", "ENV")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringDEFC",
		"--string-ptr-zero-ptr=flag",
		"--special-string=flag",
		"--special-string-ptr-zero-ptr=flag",
		"--special-ptr-string=4",
	}

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringConfiguration),
		&magiccfg.Options{
			ConfigFileEnvName: "CONFIG_FILE",
		},
	)
	assert.NoError(t, err)

	cfg, errs := c.ApplyDefaults().
		ParseEnv().
		ParseFiles().
		ParseFlags().
		Validate().
		Config()
	assert.Len(t, errs, 1)

	assert.ErrorContains(t, errs[0], "StringNoDef")
	assert.ErrorContains(t, errs[0], "required")

	assert.Equal(t, cfg.String, "default")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "ENV")
	assert.NotNil(t, cfg.StringPtrNoDef)
	assert.Equal(t, *cfg.StringPtrNoDef, "file")
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "flag")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("file"))
	assert.Nil(t, cfg.SpecialStringPtrNoDef)
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "4")
	assert.NotNil(t, cfg.SpecialPtrStringNoDef)
	assert.Equal(t, *cfg.SpecialPtrStringNoDef, "file")
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "")
}

func TestStringDECF(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/partial_string.yaml")
	t.Setenv("STRING_PTR", "ENV")
	t.Setenv("STRING_PTR_NO_DEF", "ENV")
	t.Setenv("SPECIAL_STRING", "ENV")
	t.Setenv("SPECIAL_STRING_PTR_ZERO_PTR", "ENV")
	t.Setenv("SPECIAL_PTR_STRING", "2")
	t.Setenv("SPECIAL_PTR_STRING_NO_DEF", "ENV")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringDECF",
		"--string-ptr-zero-ptr=flag",
		"--special-string=flag",
		"--special-string-ptr-zero-ptr=flag",
		"--special-ptr-string=4",
	}

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringConfiguration),
		&magiccfg.Options{
			ConfigFileEnvName: "CONFIG_FILE",
		},
	)
	assert.NoError(t, err)

	cfg, errs := c.ApplyDefaults().
		ParseEnv().
		ParseFlags().
		ParseFiles().
		Validate().
		Config()
	assert.Len(t, errs, 1)

	assert.ErrorContains(t, errs[0], "StringNoDef")
	assert.ErrorContains(t, errs[0], "required")

	assert.Equal(t, cfg.String, "default")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "ENV")
	assert.NotNil(t, cfg.StringPtrNoDef)
	assert.Equal(t, *cfg.StringPtrNoDef, "file")
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "file")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("flag"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("file"))
	assert.Nil(t, cfg.SpecialStringPtrNoDef)
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString("file"))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "3")
	assert.NotNil(t, cfg.SpecialPtrStringNoDef)
	assert.Equal(t, *cfg.SpecialPtrStringNoDef, "file")
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "")
}

func TestStringDCFE(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/partial_string.yaml")
	t.Setenv("STRING_PTR", "ENV")
	t.Setenv("STRING_PTR_NO_DEF", "ENV")
	t.Setenv("SPECIAL_STRING", "ENV")
	t.Setenv("SPECIAL_STRING_PTR_ZERO_PTR", "ENV")
	t.Setenv("SPECIAL_PTR_STRING", "2")
	t.Setenv("SPECIAL_PTR_STRING_NO_DEF", "ENV")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringDCFE",
		"--string-ptr-zero-ptr=flag",
		"--special-string=flag",
		"--special-string-ptr-zero-ptr=flag",
		"--special-ptr-string=4",
	}

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringConfiguration),
		&magiccfg.Options{
			ConfigFileEnvName: "CONFIG_FILE",
		},
	)
	assert.NoError(t, err)

	cfg, errs := c.ApplyDefaults().
		ParseFlags().
		ParseFiles().
		ParseEnv().
		Validate().
		Config()
	assert.Len(t, errs, 1)

	assert.ErrorContains(t, errs[0], "StringNoDef")
	assert.ErrorContains(t, errs[0], "required")

	assert.Equal(t, cfg.String, "default")
	assert.NotNil(t, cfg.StringPtr)
	assert.Equal(t, *cfg.StringPtr, "ENV")
	assert.NotNil(t, cfg.StringPtrNoDef)
	assert.Equal(t, *cfg.StringPtrNoDef, "ENV")
	assert.NotNil(t, cfg.StringPtrZeroPtr)
	assert.Equal(t, *cfg.StringPtrZeroPtr, "file")
	assert.Equal(t, cfg.SpecialString, configs.SpecialString("ENV"))
	assert.NotNil(t, cfg.SpecialStringPtr)
	assert.Equal(t, *cfg.SpecialStringPtr, configs.SpecialString("file"))
	assert.Nil(t, cfg.SpecialStringPtrNoDef)
	assert.NotNil(t, cfg.SpecialStringPtrZeroPtr)
	assert.Equal(t, *cfg.SpecialStringPtrZeroPtr, configs.SpecialString("ENV"))
	assert.NotNil(t, cfg.SpecialPtrString)
	assert.Equal(t, *cfg.SpecialPtrString, "2")
	assert.NotNil(t, cfg.SpecialPtrStringNoDef)
	assert.Equal(t, *cfg.SpecialPtrStringNoDef, "ENV")
	assert.NotNil(t, cfg.SpecialPtrStringZeroPtr)
	assert.Equal(t, *cfg.SpecialPtrStringZeroPtr, "")
}
