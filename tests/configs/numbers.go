package configs

type SpecialInt int
type SpecialPtrInt *int
type SpecialInt8 int8
type SpecialPtrInt8 *int8
type SpecialInt16 int16
type SpecialPtrInt16 *int16
type SpecialInt32 int32
type SpecialPtrInt32 *int32
type SpecialInt64 int64
type SpecialPtrInt64 *int64

type SpecialUInt uint
type SpecialPtrUInt *uint
type SpecialUInt8 uint8
type SpecialPtrUInt8 *uint8
type SpecialUInt16 uint16
type SpecialPtrUInt16 *uint16
type SpecialUInt32 uint32
type SpecialPtrUInt32 *uint32
type SpecialUInt64 uint64
type SpecialPtrUInt64 *uint64

type IntegerConfiguration struct {
	Int                int        `def:"-3" validate:"min=0"`
	IntNoDef           int        `validate:"ltfield=Int"`
	SpecialInt         SpecialInt `def:"-3"`
	SpecialIntNoDef    SpecialInt
	Int8               int8 `def:"-3"`
	Int8NoDef          int8
	SpecialInt8        SpecialInt8 `def:"-3"`
	SpecialInt8NoDef   SpecialInt8
	Int16              int16 `def:"-3"`
	Int16NoDef         int16
	SpecialInt16       SpecialInt16 `def:"-3"`
	SpecialInt16NoDef  SpecialInt16
	Int32              int32 `def:"-3"`
	Int32NoDef         int32
	SpecialInt32       SpecialInt32 `def:"-3"`
	SpecialInt32NoDef  SpecialInt32
	Int64              int64 `def:"-3"`
	Int64NoDef         int64
	SpecialInt64       SpecialInt64 `def:"-3"`
	SpecialInt64NoDef  SpecialInt64
	UInt               uint `def:"3"`
	UIntNoDef          uint
	SpecialUInt        SpecialUInt `def:"3"`
	SpecialUIntNoDef   SpecialUInt
	UInt8              uint8 `def:"6"`
	UInt8NoDef         uint8
	SpecialUInt8       SpecialUInt8 `def:"3"`
	SpecialUInt8NoDef  SpecialUInt8
	UInt16             uint16 `def:"3"`
	UInt16NoDef        uint16
	SpecialUInt16      SpecialUInt16 `def:"3"`
	SpecialUInt16NoDef SpecialUInt16
	UInt32             uint32 `def:"3"`
	UInt32NoDef        uint32
	SpecialUInt32      SpecialUInt32 `def:"3"`
	SpecialUInt32NoDef SpecialUInt32
	UInt64             uint64 `def:"3"`
	UInt64NoDef        uint64
	SpecialUInt64      SpecialUInt64 `def:"3"`
	SpecialUInt64NoDef SpecialUInt64
}

type BadDefaultIntConfiguration struct {
	Int                 int           `def:"a"`
	SpecialInt          SpecialInt    `def:"a"`
	Int8                int8          `def:"a"`
	SpecialInt8         SpecialInt8   `def:"a"`
	Int16               int16         `def:"a"`
	SpecialInt16        SpecialInt16  `def:"a"`
	Int32               int32         `def:"a"`
	SpecialInt32        SpecialInt32  `def:"a"`
	Int64               int64         `def:"a"`
	SpecialInt64        SpecialInt64  `def:"a"`
	UIntLetter          uint8         `def:"a"`
	SpecialUIntLetter   SpecialUInt   `def:"a"`
	UInt8Letter         uint8         `def:"a"`
	SpecialUInt8Letter  SpecialUInt8  `def:"a"`
	UInt16Letter        uint16        `def:"a"`
	SpecialUInt16Letter SpecialUInt16 `def:"a"`
	UInt32Letter        uint16        `def:"a"`
	SpecialUInt32Letter SpecialUInt32 `def:"a"`
	UInt64Letter        uint16        `def:"a"`
	SpecialUInt64Letter SpecialUInt64 `def:"a"`
	UIntNeg             uint          `def:"-1"`
	SpecialUIntNeg      SpecialUInt   `def:"-1"`
	UInt8               uint8         `def:"-1"`
	SpecialUInt8Neg     SpecialUInt8  `def:"-1"`
	UInt16              uint16        `def:"-1"`
	SpecialUInt16Neg    SpecialUInt16 `def:"-1"`
	UInt32              uint32        `def:"-1"`
	SpecialUInt32Neg    SpecialUInt32 `def:"-1"`
	UInt64              uint64        `def:"-1"`
	SpecialUInt64Neg    SpecialUInt64 `def:"-1"`
}

type IntegerPointerConfiguration struct {
	Int                     *int `def:"-3"`
	IntNoDef                *int
	IntZeroPtr              *int        `def:"*0*"`
	SpecialIntPtr           *SpecialInt `def:"-3"`
	SpecialIntPtrNoDef      *SpecialInt
	SpecialIntPtrZeroPtr    *SpecialInt   `def:"*0*"`
	SpecialPtrInt           SpecialPtrInt `def:"-3"`
	SpecialPtrIntNoDef      SpecialPtrInt
	SpecialPtrIntZeroPtr    SpecialPtrInt `def:"*0*"`
	Int8                    *int8         `def:"-3"`
	Int8NoDef               *int8
	Int8ZeroPtr             *int8        `def:"*0*"`
	SpecialInt8Ptr          *SpecialInt8 `def:"-3"`
	SpecialInt8PtrNoDef     *SpecialInt8
	SpecialInt8PtrZeroPtr   *SpecialInt8   `def:"*0*"`
	SpecialPtrInt8          SpecialPtrInt8 `def:"-3"`
	SpecialPtrInt8NoDef     SpecialPtrInt8
	SpecialPtrInt8ZeroPtr   SpecialPtrInt8 `def:"*0*"`
	Int16                   *int16         `def:"-3"`
	Int16NoDef              *int16
	Int16ZeroPtr            *int16        `def:"*0*"`
	SpecialInt16Ptr         *SpecialInt16 `def:"-3"`
	SpecialInt16PtrNoDef    *SpecialInt16
	SpecialInt16PtrZeroPtr  *SpecialInt16   `def:"*0*"`
	SpecialPtrInt16         SpecialPtrInt16 `def:"-3"`
	SpecialPtrInt16NoDef    SpecialPtrInt16
	SpecialPtrInt16ZeroPtr  SpecialPtrInt16 `def:"*0*"`
	Int32                   *int32          `def:"-3"`
	Int32NoDef              *int32
	Int32ZeroPtr            *int32        `def:"*0*"`
	SpecialInt32Ptr         *SpecialInt32 `def:"-3"`
	SpecialInt32PtrNoDef    *SpecialInt32
	SpecialInt32PtrZeroPtr  *SpecialInt32   `def:"*0*"`
	SpecialPtrInt32         SpecialPtrInt32 `def:"-3"`
	SpecialPtrInt32NoDef    SpecialPtrInt32
	SpecialPtrInt32ZeroPtr  SpecialPtrInt32 `def:"*0*"`
	Int64                   *int64          `def:"-3"`
	Int64NoDef              *int64
	Int64ZeroPtr            *int64        `def:"*0*"`
	SpecialInt64Ptr         *SpecialInt64 `def:"-3"`
	SpecialInt64PtrNoDef    *SpecialInt64
	SpecialInt64PtrZeroPtr  *SpecialInt64   `def:"*0*"`
	SpecialPtrInt64         SpecialPtrInt64 `def:"-3"`
	SpecialPtrInt64NoDef    SpecialPtrInt64
	SpecialPtrInt64ZeroPtr  SpecialPtrInt64 `def:"*0*"`
	UInt                    *uint           `def:"3"`
	UIntNoDef               *uint
	UIntZeroPtr             *uint        `def:"*0*"`
	SpecialUIntPtr          *SpecialUInt `def:"3"`
	SpecialUIntPtrNoDef     *SpecialUInt
	SpecialUIntPtrZeroPtr   *SpecialUInt   `def:"*0*"`
	SpecialPtrUInt          SpecialPtrUInt `def:"3"`
	SpecialPtrUIntNoDef     SpecialPtrUInt
	SpecialPtrUIntZeroPtr   SpecialPtrUInt `def:"*0*"`
	UInt8                   *uint8         `def:"3"`
	UInt8NoDef              *uint8
	UInt8ZeroPtr            *uint8        `def:"*0*"`
	SpecialUInt8Ptr         *SpecialUInt8 `def:"3"`
	SpecialUInt8PtrNoDef    *SpecialUInt8
	SpecialUInt8PtrZeroPtr  *SpecialUInt8   `def:"*0*"`
	SpecialPtrUInt8         SpecialPtrUInt8 `def:"3"`
	SpecialPtrUInt8NoDef    SpecialPtrUInt8
	SpecialPtrUInt8ZeroPtr  SpecialPtrUInt8 `def:"*0*"`
	UInt16                  *uint16         `def:"3"`
	UInt16NoDef             *uint16
	UInt16ZeroPtr           *uint16        `def:"*0*"`
	SpecialUInt16Ptr        *SpecialUInt16 `def:"3"`
	SpecialUInt16PtrNoDef   *SpecialUInt16
	SpecialUInt16PtrZeroPtr *SpecialUInt16   `def:"*0*"`
	SpecialPtrUInt16        SpecialPtrUInt16 `def:"3"`
	SpecialPtrUInt16NoDef   SpecialPtrUInt16
	SpecialPtrUInt16ZeroPtr SpecialPtrUInt16 `def:"*0*"`
	UInt32                  *uint32          `def:"3"`
	UInt32NoDef             *uint32
	UInt32ZeroPtr           *uint32        `def:"*0*"`
	SpecialUInt32Ptr        *SpecialUInt32 `def:"3"`
	SpecialUInt32PtrNoDef   *SpecialUInt32
	SpecialUInt32PtrZeroPtr *SpecialUInt32   `def:"*0*"`
	SpecialPtrUInt32        SpecialPtrUInt32 `def:"3"`
	SpecialPtrUInt32NoDef   SpecialPtrUInt32
	SpecialPtrUInt32ZeroPtr SpecialPtrUInt32 `def:"*0*"`
	UInt64                  *uint64          `def:"3"`
	UInt64NoDef             *uint64
	UInt64ZeroPtr           *uint64        `def:"*0*"`
	SpecialUInt64Ptr        *SpecialUInt64 `def:"3"`
	SpecialUInt64PtrNoDef   *SpecialUInt64
	SpecialUInt64PtrZeroPtr *SpecialUInt64   `def:"*0*"`
	SpecialPtrUInt64        SpecialPtrUInt64 `def:"3"`
	SpecialPtrUInt64NoDef   SpecialPtrUInt64
	SpecialPtrUInt64ZeroPtr SpecialPtrUInt64 `def:"*0*"`
}

type BadDefaultIntPtrConfiguration struct {
	Int                    *int             `def:"a"`
	SpecialInt             *SpecialInt      `def:"a"`
	SpecialPtrInt          SpecialPtrInt    `def:"a"`
	Int8                   *int8            `def:"a"`
	SpecialInt8            *SpecialInt8     `def:"a"`
	SpecialPtrInt8         SpecialPtrInt8   `def:"a"`
	Int16                  *int16           `def:"a"`
	SpecialInt16           *SpecialInt16    `def:"a"`
	SpecialPtrInt16        SpecialPtrInt16  `def:"a"`
	Int32                  *int32           `def:"a"`
	SpecialInt32           *SpecialInt32    `def:"a"`
	SpecialPtrInt32        SpecialPtrInt32  `def:"a"`
	Int64                  *int64           `def:"a"`
	SpecialInt64           *SpecialInt64    `def:"a"`
	SpecialPtrInt64        SpecialPtrInt64  `def:"a"`
	UIntLetter             *uint8           `def:"a"`
	SpecialUIntLetter      *SpecialUInt     `def:"a"`
	SpecialPtrUIntLetter   SpecialPtrUInt   `def:"a"`
	UInt8Letter            *uint8           `def:"a"`
	SpecialUInt8Letter     *SpecialUInt8    `def:"a"`
	SpecialPtrUInt8Letter  SpecialPtrUInt8  `def:"a"`
	UInt16Letter           *uint16          `def:"a"`
	SpecialUInt16Letter    *SpecialUInt16   `def:"a"`
	SpecialPtrUInt16Letter SpecialPtrUInt16 `def:"a"`
	UInt32Letter           *uint16          `def:"a"`
	SpecialUInt32Letter    *SpecialUInt32   `def:"a"`
	SpecialPtrUInt32Letter SpecialPtrUInt32 `def:"a"`
	UInt64Letter           *uint16          `def:"a"`
	SpecialUInt64Letter    *SpecialUInt64   `def:"a"`
	SpecialPtrUInt64Letter SpecialPtrUInt64 `def:"a"`
	UIntNeg                *uint            `def:"-1"`
	SpecialUIntNeg         *SpecialUInt     `def:"-1"`
	SpecialPtrUIntNeg      SpecialPtrUInt   `def:"-1"`
	UInt8                  *uint8           `def:"-1"`
	SpecialUInt8Neg        *SpecialUInt8    `def:"-1"`
	SpecialPtrUInt8Neg     SpecialPtrUInt8  `def:"-1"`
	UInt16                 *uint16          `def:"-1"`
	SpecialUInt16Neg       *SpecialUInt16   `def:"-1"`
	SpecialPtrUInt16Neg    SpecialPtrUInt16 `def:"-1"`
	UInt32                 *uint32          `def:"-1"`
	SpecialUInt32Neg       *SpecialUInt32   `def:"-1"`
	SpecialPtrUInt32Neg    SpecialPtrUInt32 `def:"-1"`
	UInt64                 *uint64          `def:"-1"`
	SpecialUInt64Neg       *SpecialUInt64   `def:"-1"`
	SpecialPtrUInt64Neg    SpecialPtrUInt64 `def:"-1"`
}

type SpecialFloat32 float32
type SpecialPtrFloat32 *float32
type SpecialFloat64 float64
type SpecialPtrFloat64 *float64

type FloatConfiguration struct {
	Float32             float32 `def:"3.14"`
	Float32Zero         float32 `def:"0"`
	Float32NoDef        float32
	SpecialFloat32      SpecialFloat32 `def:"3.14"`
	SpecialFloat32Zero  SpecialFloat32 `def:"0"`
	SpecialFloat32NoDef SpecialFloat32
	Float64             float64 `def:"3.14"`
	Float64Zero         float64 `def:"0"`
	Float64NoDef        float64
	SpecialFloat64      SpecialFloat64 `def:"3.14"`
	SpecialFloat64Zero  SpecialFloat64 `def:"0"`
	SpecialFloat64NoDef SpecialFloat64
}

type BadDefaultFloatConfiguration struct {
	Float32        float32        `def:"a"`
	SpecialFloat32 SpecialFloat32 `def:"a"`
	Float64        float64        `def:"b"`
	SpecialFloat64 SpecialFloat64 `def:"b"`
}
type FloatPtrConfiguration struct {
	Float32                  *float32 `def:"3.14"`
	Float32Zero              *float32 `def:"0"`
	Float32NoDef             *float32
	Float32ZeroPtr           *float32        `def:"*0*"`
	SpecialFloat32           *SpecialFloat32 `def:"3.14"`
	SpecialFloat32Zero       *SpecialFloat32 `def:"0"`
	SpecialFloat32NoDef      *SpecialFloat32
	SpecialFloat32ZeroPtr    *SpecialFloat32   `def:"*0*"`
	SpecialPtrFloat32        SpecialPtrFloat32 `def:"3.14"`
	SpecialPtrFloat32Zero    SpecialPtrFloat32 `def:"0"`
	SpecialPtrFloat32NoDef   SpecialPtrFloat32
	SpecialPtrFloat32ZeroPtr SpecialPtrFloat32 `def:"*0*"`
	Float64                  *float64          `def:"3.14"`
	Float64Zero              *float64          `def:"0"`
	Float64NoDef             *float64
	Float64ZeroPtr           *float64        `def:"*0*"`
	SpecialFloat64           *SpecialFloat64 `def:"3.14"`
	SpecialFloat64Zero       *SpecialFloat64 `def:"0"`
	SpecialFloat64NoDef      *SpecialFloat64
	SpecialFloat64ZeroPtr    *SpecialFloat64   `def:"*0*"`
	SpecialPtrFloat64        SpecialPtrFloat64 `def:"3.14"`
	SpecialPtrFloat64Zero    SpecialPtrFloat64 `def:"0"`
	SpecialPtrFloat64NoDef   SpecialPtrFloat64
	SpecialPtrFloat64ZeroPtr SpecialPtrFloat64 `def:"*0*"`
}

type BadDefaultFloatPtrConfiguration struct {
	Float32           *float32          `def:"a"`
	SpecialFloat32    *SpecialFloat32   `def:"a"`
	SpecialPtrFloat32 SpecialPtrFloat32 `def:"a"`
	Float64           *float64          `def:"b"`
	SpecialFloat64    *SpecialFloat64   `def:"b"`
	SpecialPtrFloat64 SpecialPtrFloat64 `def:"b"`
}

type SpecialComplex64 complex64
type SpecialPtrComplex64 *complex64
type SpecialComplex128 complex128
type SpecialPtrComplex128 *complex128

type ComplexConfiguration struct {
	Complex64              complex64 `def:"3.14-2i"`
	Complex64Zero          complex64 `def:"0"`
	Complex64NoDef         complex64
	SpecialComplex64       SpecialComplex64 `def:"3.14-2i"`
	SpecialComplex64Zero   SpecialComplex64 `def:"0"`
	SpecialComplex64NoDef  SpecialComplex64
	Complex128             complex128 `def:"3.14-2i"`
	Complex128Zero         complex128 `def:"0"`
	Complex128NoDef        complex128
	SpecialComplex128      SpecialComplex128 `def:"3.14-2i"`
	SpecialComplex128Zero  SpecialComplex128 `def:"0"`
	SpecialComplex128NoDef SpecialComplex128
}

type BadDefaultComplexConfiguration struct {
	Complex64         complex64         `def:"a+bi"`
	SpecialComplex64  SpecialComplex64  `def:"a+bi"`
	Complex128        complex128        `def:"a+bi"`
	SpecialComplex128 SpecialComplex128 `def:"a+bi"`
}

type ComplexPtrConfiguration struct {
	Complex64                   *complex64 `def:"3.14-2i"`
	Complex64Zero               *complex64 `def:"0"`
	Complex64NoDef              *complex64
	Complex64ZeroPtr            *complex64        `def:"*0*"`
	SpecialComplex64            *SpecialComplex64 `def:"3.14-2i"`
	SpecialComplex64Zero        *SpecialComplex64 `def:"0"`
	SpecialComplex64NoDef       *SpecialComplex64
	SpecialComplex64ZeroPtr     *SpecialComplex64   `def:"*0*"`
	SpecialPtrComplex64         SpecialPtrComplex64 `def:"3.14-2i"`
	SpecialPtrComplex64Zero     SpecialPtrComplex64 `def:"0"`
	SpecialPtrComplex64NoDef    SpecialPtrComplex64
	SpecialPtrComplex64ZeroPtr  SpecialPtrComplex64 `def:"*0*"`
	Complex128                  *complex128         `def:"3.14-2i"`
	Complex128Zero              *complex128         `def:"0"`
	Complex128NoDef             *complex128
	Complex128ZeroPtr           *complex128        `def:"*0*"`
	SpecialComplex128           *SpecialComplex128 `def:"3.14-2i"`
	SpecialComplex128Zero       *SpecialComplex128 `def:"0"`
	SpecialComplex128NoDef      *SpecialComplex128
	SpecialComplex128ZeroPtr    *SpecialComplex128   `def:"*0*"`
	SpecialPtrComplex128        SpecialPtrComplex128 `def:"3.14-2i"`
	SpecialPtrComplex128Zero    SpecialPtrComplex128 `def:"0"`
	SpecialPtrComplex128NoDef   SpecialPtrComplex128
	SpecialPtrComplex128ZeroPtr SpecialPtrComplex128 `def:"*0*"`
}

type BadDefaultComplexPtrConfiguration struct {
	Complex64            *complex64           `def:"a+bi"`
	SpecialComplex64     *SpecialComplex64    `def:"a+bi"`
	SpecialPtrComplex64  SpecialPtrComplex64  `def:"a+bi"`
	Complex128           *complex128          `def:"a+bi"`
	SpecialComplex128    *SpecialComplex128   `def:"a+bi"`
	SpecialPtrComplex128 SpecialPtrComplex128 `def:"a+bi"`
}
