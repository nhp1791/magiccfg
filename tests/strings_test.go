package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"magiccfg"
	"magiccfg/tests/configs"
)

func TestStringsDefaultsValidate(t *testing.T) {
	c, err := magiccfg.NewMagicConfig(new(configs.StringsConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ApplyDefaults().
		Validate().
		Config()
	assert.Len(t, errs, 3)

	assert.ErrorContains(t, errs[0], "StringsNoDef")
	assert.ErrorContains(t, errs[0], "min")

	assert.ErrorContains(t, errs[1], "SpecialStrings")
	assert.ErrorContains(t, errs[1], "max")

	assert.ErrorContains(t, errs[2], "SpecialStringsNoDef")
	assert.ErrorContains(t, errs[2], "required")

	d := "default"
	ssd := configs.SpecialString(d)
	sse := configs.SpecialPtrString(&d)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, ""), "default")
	assert.Nil(t, cfg.StringsNoDef)
	assert.NotNil(t, cfg.StringsZeroDef)
	assert.Len(t, cfg.StringsZeroDef, 0)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, ""), ssd)
	assert.Nil(t, cfg.SpecialStringsNoDef)
	assert.NotNil(t, cfg.SpecialStringsZeroDef)
	assert.Len(t, cfg.SpecialStringsZeroDef, 0)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, ""), sse)
	assert.Nil(t, cfg.SpecialPtrStringsNoDef)
	assert.NotNil(t, cfg.SpecialPtrStringsZeroDef)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 0)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, ""), &d)
	assert.Nil(t, cfg.StringsPtrNoDef)
	assert.NotNil(t, cfg.StringsPtrZeroDef)
	assert.Len(t, cfg.StringsPtrZeroDef, 0)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, ""), &ssd)
	assert.Nil(t, cfg.SpecialStringsPtrNoDef)
	assert.NotNil(t, cfg.SpecialStringsPtrZeroDef)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 0)
}

func TestStringsEnv(t *testing.T) {
	t.Setenv("STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")

	c, err := magiccfg.NewMagicConfig(new(configs.StringsConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ParseEnv().Config()
	assert.Empty(t, errs)

	e := "ENV_STRINGS_VAL"
	sse := configs.SpecialString(e)
	spe := configs.SpecialPtrString(&e)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, "_"), e)
	assert.Len(t, cfg.StringsNoDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsNoDef, "_"), e)
	assert.Len(t, cfg.StringsZeroDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsZeroDef, "_"), e)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, "_"), sse)
	assert.Len(t, cfg.SpecialStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsNoDef, "_"), sse)
	assert.Len(t, cfg.SpecialStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsZeroDef, "_"), sse)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, "_"), spe)
	assert.Len(t, cfg.SpecialPtrStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsNoDef, "_"), spe)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsZeroDef, "_"), spe)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, "_"), &e)
	assert.Len(t, cfg.StringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrNoDef, "_"), &e)
	assert.Len(t, cfg.StringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrZeroDef, "_"), &e)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, "_"), &sse)
	assert.Len(t, cfg.SpecialStringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrNoDef, "_"), &sse)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrZeroDef, "_"), &sse)
}

func TestStringsFile(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/full_string.yaml")

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringsConfiguration),
		&magiccfg.Options{
			ConfigFileEnvName: "CONFIG_FILE",
		},
	)
	assert.NoError(t, err)

	cfg, errs := c.ParseFiles().Config()
	assert.Empty(t, errs)

	f := "file_strings_val"
	ssf := configs.SpecialString(f)
	spf := configs.SpecialPtrString(&f)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, "_"), f)
	assert.Len(t, cfg.StringsNoDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsNoDef, "_"), f)
	assert.Len(t, cfg.StringsZeroDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsZeroDef, "_"), f)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, "_"), ssf)
	assert.Len(t, cfg.SpecialStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsNoDef, "_"), ssf)
	assert.Len(t, cfg.SpecialStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsZeroDef, "_"), ssf)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, "_"), spf)
	assert.Len(t, cfg.SpecialPtrStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsNoDef, "_"), spf)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsZeroDef, "_"), spf)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, "_"), &f)
	assert.Len(t, cfg.StringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrNoDef, "_"), &f)
	assert.Len(t, cfg.StringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrZeroDef, "_"), &f)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrNoDef, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrZeroDef, "_"), &ssf)
}

func TestStringsFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringsFlags",
		"--strings=flag,strings,val",
		"--strings-no-def=flag,strings,val",
		"--strings-zero-def=flag,strings,val",
		"--special-strings=flag,strings,val",
		"--special-strings-no-def=flag,strings,val",
		"--special-strings-zero-def=flag,strings,val",
		"--special-ptr-strings=flag,strings,val",
		"--special-ptr-strings-no-def=flag,strings,val",
		"--special-ptr-strings-zero-def=flag,strings,val",
		"--strings-ptr=flag,strings,val",
		"--strings-ptr-no-def=flag,strings,val",
		"--strings-ptr-zero-def=flag,strings,val",
		"--special-strings-ptr=flag,strings,val",
		"--special-strings-ptr-no-def=flag,strings,val",
		"--special-strings-ptr-zero-def=flag,strings,val",
	}
	c, err := magiccfg.NewMagicConfig(new(configs.StringsConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ParseFlags().Config()
	assert.Empty(t, errs)

	f := "flag_strings_val"
	ssf := configs.SpecialString(f)
	spf := configs.SpecialPtrString(&f)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, "_"), f)
	assert.Len(t, cfg.StringsNoDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsNoDef, "_"), f)
	assert.Len(t, cfg.StringsZeroDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsZeroDef, "_"), f)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, "_"), ssf)
	assert.Len(t, cfg.SpecialStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsNoDef, "_"), ssf)
	assert.Len(t, cfg.SpecialStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsZeroDef, "_"), ssf)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, "_"), spf)
	assert.Len(t, cfg.SpecialPtrStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsNoDef, "_"), spf)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsZeroDef, "_"), spf)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, "_"), &f)
	assert.Len(t, cfg.StringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrNoDef, "_"), &f)
	assert.Len(t, cfg.StringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrZeroDef, "_"), &f)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrNoDef, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrZeroDef, "_"), &ssf)
}

func TestStringsDEFC(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/partial_string.yaml")

	t.Setenv("STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringsDEFC",
		"--strings-no-def=flag,strings,val",
		"--strings-zero-def=flag,strings,val",
		"--special-strings-zero-def=flag,strings,val",
		"--special-ptr-strings=flag,strings,val",
		"--special-ptr-strings-no-def=flag,strings,val",
		"--strings-ptr=flag,strings,val",
		"--strings-ptr-zero-def=flag,strings,val",
		"--special-strings-ptr=flag,strings,val",
		"--special-strings-ptr-zero-def=flag,strings,val",
	}

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringsConfiguration),
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
	assert.Len(t, errs, 2)

	assert.ErrorContains(t, errs[0], "SpecialStrings")
	assert.ErrorContains(t, errs[0], "max")

	assert.ErrorContains(t, errs[1], "SpecialStringsNoDef")
	assert.ErrorContains(t, errs[1], "required")

	g := "flag_strings_val"
	ssg := configs.SpecialString(g)
	spg := configs.SpecialPtrString(&g)
	e := "ENV_STRINGS_VAL"
	sse := configs.SpecialString(e)
	f := "file_strings_val"
	ssf := configs.SpecialString(f)
	spe := configs.SpecialPtrString(&f)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, ""), "default")
	assert.Len(t, cfg.StringsNoDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsNoDef, "_"), g)
	assert.Len(t, cfg.StringsZeroDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsZeroDef, "_"), g)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, "_"), sse)
	assert.Nil(t, cfg.SpecialStringsNoDef)
	assert.Len(t, cfg.SpecialStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsZeroDef, "_"), ssg)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, "_"), spg)
	assert.Len(t, cfg.SpecialPtrStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsNoDef, "_"), spg)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsZeroDef, "_"), spe)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, "_"), &g)
	assert.Len(t, cfg.StringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrNoDef, "_"), &f)
	assert.Len(t, cfg.StringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrZeroDef, "_"), &g)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, "_"), &ssg)
	assert.Len(t, cfg.SpecialStringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrNoDef, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrZeroDef, "_"), &ssg)
}

func TestStringsDECF(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/partial_string.yaml")

	t.Setenv("STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringsDECF",
		"--strings-no-def=flag,strings,val",
		"--strings-zero-def=flag,strings,val",
		"--special-strings-zero-def=flag,strings,val",
		"--special-ptr-strings=flag,strings,val",
		"--special-ptr-strings-no-def=flag,strings,val",
		"--strings-ptr=flag,strings,val",
		"--strings-ptr-zero-def=flag,strings,val",
		"--special-strings-ptr=flag,strings,val",
		"--special-strings-ptr-zero-def=flag,strings,val",
	}

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringsConfiguration),
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
	assert.Len(t, errs, 2)

	assert.ErrorContains(t, errs[0], "SpecialStrings")
	assert.ErrorContains(t, errs[0], "max")

	assert.ErrorContains(t, errs[1], "SpecialStringsNoDef")
	assert.ErrorContains(t, errs[1], "required")

	g := "flag_strings_val"
	ssg := configs.SpecialString(g)
	e := "ENV_STRINGS_VAL"
	sse := configs.SpecialString(e)
	f := "file_strings_val"
	ssf := configs.SpecialString(f)
	spf := configs.SpecialPtrString(&f)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, ""), "default")
	assert.Len(t, cfg.StringsNoDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsNoDef, "_"), f)
	assert.Len(t, cfg.StringsZeroDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsZeroDef, "_"), f)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, "_"), sse)
	assert.Nil(t, cfg.SpecialStringsNoDef)
	assert.Len(t, cfg.SpecialStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsZeroDef, "_"), ssf)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, "_"), spf)
	assert.Len(t, cfg.SpecialPtrStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsNoDef, "_"), spf)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsZeroDef, "_"), spf)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, "_"), &f)
	assert.Len(t, cfg.StringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrNoDef, "_"), &f)
	assert.Len(t, cfg.StringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrZeroDef, "_"), &f)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, "_"), &ssg)
	assert.Len(t, cfg.SpecialStringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrNoDef, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrZeroDef, "_"), &ssg)
}

func TestStringsDCFE(t *testing.T) {
	t.Setenv("CONFIG_FILE", "./configfiles/partial_string.yaml")

	t.Setenv("STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_PTR_STRINGS_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_NO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR", "ENV,STRINGS,VAL")
	t.Setenv("SPECIAL_STRINGS_PTR_ZERO_DEF", "ENV,STRINGS,VAL")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original arguments after test ends

	os.Args = []string{
		"TestStringsDCFE",
		"--strings-no-def=flag,strings,val",
		"--strings-zero-def=flag,strings,val",
		"--special-strings-zero-def=flag,strings,val",
		"--special-ptr-strings=flag,strings,val",
		"--special-ptr-strings-no-def=flag,strings,val",
		"--strings-ptr=flag,strings,val",
		"--strings-ptr-zero-def=flag,strings,val",
		"--special-strings-ptr=flag,strings,val",
		"--special-strings-ptr-zero-def=flag,strings,val",
	}

	c, err := magiccfg.NewMagicConfig(
		new(configs.StringsConfiguration),
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
	assert.Len(t, errs, 2)

	assert.ErrorContains(t, errs[0], "SpecialStrings")
	assert.ErrorContains(t, errs[0], "max")

	assert.ErrorContains(t, errs[1], "SpecialStringsNoDef")
	assert.ErrorContains(t, errs[1], "required")

	e := "ENV_STRINGS_VAL"
	sse := configs.SpecialString(e)
	f := "file_strings_val"
	ssf := configs.SpecialString(f)
	spe := configs.SpecialPtrString(&e)

	assert.Len(t, cfg.Strings, 3)
	assert.Equal(t, strings.Join(cfg.Strings, ""), "default")
	assert.Len(t, cfg.StringsNoDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsNoDef, "_"), f)
	assert.Len(t, cfg.StringsZeroDef, 3)
	assert.Equal(t, strings.Join(cfg.StringsZeroDef, "_"), e)
	assert.Len(t, cfg.SpecialStrings, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStrings, "_"), sse)
	assert.Nil(t, cfg.SpecialStringsNoDef)
	assert.Len(t, cfg.SpecialStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringTypeJoin(cfg.SpecialStringsZeroDef, "_"), ssf)
	assert.Len(t, cfg.SpecialPtrStrings, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStrings, "_"), spe)
	assert.Len(t, cfg.SpecialPtrStringsNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsNoDef, "_"), spe)
	assert.Len(t, cfg.SpecialPtrStringsZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.SpecialPtrStringsZeroDef, "_"), spe)
	assert.Len(t, cfg.StringsPtr, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtr, "_"), &e)
	assert.Len(t, cfg.StringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrNoDef, "_"), &e)
	assert.Len(t, cfg.StringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.StringPtrTypeJoin(cfg.StringsPtrZeroDef, "_"), &e)
	assert.Len(t, cfg.SpecialStringsPtr, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtr, "_"), &sse)
	assert.Len(t, cfg.SpecialStringsPtrNoDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrNoDef, "_"), &ssf)
	assert.Len(t, cfg.SpecialStringsPtrZeroDef, 3)
	assert.Equal(t, magiccfg.PtrStringTypeJoin(cfg.SpecialStringsPtrZeroDef, "_"), &sse)
}
