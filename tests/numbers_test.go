package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magiccfg"
	"magiccfg/tests/configs"
)

func TestIntegerDefaults(t *testing.T) {
	c, err := magiccfg.NewMagicConfig(new(configs.IntegerConfiguration), nil)
	assert.NoError(t, err)

	cfg, errs := c.ApplyDefaults().
		Validate().
		Config()
	assert.Len(t, errs, 2)

	assert.ErrorContains(t, errs[0], "Int")
	assert.ErrorContains(t, errs[0], "min")

	assert.ErrorContains(t, errs[1], "IntNoDef")
	assert.ErrorContains(t, errs[1], "ltfield")

	assert.Equal(t, cfg.Int, -3)
	assert.Equal(t, cfg.IntNoDef, 0)
	assert.Equal(t, cfg.SpecialInt, configs.SpecialInt(-3))
	assert.Equal(t, cfg.SpecialIntNoDef, configs.SpecialInt(0))
	assert.Equal(t, cfg.Int8, int8(-3))
	assert.Equal(t, cfg.Int8NoDef, int8(0))
	assert.Equal(t, cfg.SpecialInt8, configs.SpecialInt8(-3))
	assert.Equal(t, cfg.SpecialInt8NoDef, configs.SpecialInt8(0))
	assert.Equal(t, cfg.Int16, int16(-3))
	assert.Equal(t, cfg.Int16NoDef, int16(0))
	assert.Equal(t, cfg.SpecialInt16, configs.SpecialInt16(-3))
	assert.Equal(t, cfg.SpecialInt16NoDef, configs.SpecialInt16(0))
	assert.Equal(t, cfg.Int32, int32(-3))
	assert.Equal(t, cfg.Int32NoDef, int32(0))
	assert.Equal(t, cfg.SpecialInt32, configs.SpecialInt32(-3))
	assert.Equal(t, cfg.SpecialInt32NoDef, configs.SpecialInt32(0))
	assert.Equal(t, cfg.Int64, int64(-3))
	assert.Equal(t, cfg.Int64NoDef, int64(0))
	assert.Equal(t, cfg.SpecialInt64, configs.SpecialInt64(-3))
	assert.Equal(t, cfg.SpecialInt64NoDef, configs.SpecialInt64(0))
	assert.Equal(t, cfg.UInt, uint(3))
	assert.Equal(t, cfg.UIntNoDef, uint(0))
	assert.Equal(t, cfg.SpecialUInt, configs.SpecialUInt(3))
	assert.Equal(t, cfg.SpecialUIntNoDef, configs.SpecialUInt(0))
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
	//assert.Equal(t, cfg.Int, -3)
}
