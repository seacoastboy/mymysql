package mymy

import (
    "testing"
    "bytes"
    "math"
    "reflect"
)

var (
    Bytes  = []byte("Ala ma Kota!")
    String = "ssss" //"A kot ma Alę!"
    blob   = Blob{1, 2, 3}
    date   = Datetime{Year: 2010, Month: 12, Day: 30, Hour: 17, Minute:21}
    tstamp = Timestamp{Year: 2001, Month: 2, Day: 3, Hour: 7, Minute:2}

    pBytes  *[]byte
    pString *string
    pBlob   *Blob
    pDate   *Datetime
    pTstamp *Timestamp

    raw    = Raw{&[]byte{3, 2, 1}, MYSQL_TYPE_INT24}

    Int8   = int8(1)
    Uint8  = uint8(2)
    Int16  = int16(3)
    Uint16 = uint16(4)
    Int32  = int32(5)
    Uint32 = uint32(6)
    Int64  = int64(0x7000100020003001)
    Uint64 = uint64(0xffff0000ffff0000)
    Int    = int(7)
    Uint   = uint(8)

    Float   = float(3.14159e3)
    Float32 = float32(1e10)
    Float64 = float64(256e256)

    pInt8    *int8
    pUint8   *uint8
    pInt16   *int16
    pUint16  *uint16
    pInt32   *int32
    pUint32  *uint32
    pInt64   *int64
    pUint64  *uint64
    pInt     *int
    pUint    *uint
    pFloat   *float
    pFloat32 *float
    pFloat64 *float
)

type BindTest struct {
    val    interface{}
    typ    uint16
    is_ptr bool
    length int
}

var bindTests = []BindTest {
    BindTest{nil,     MYSQL_TYPE_NULL,       false,  0},

    BindTest{Bytes,   MYSQL_TYPE_VAR_STRING, false, -1},
    BindTest{String,  MYSQL_TYPE_STRING,     false, -1},
    BindTest{blob,    MYSQL_TYPE_BLOB,       false, -1},
    BindTest{date,    MYSQL_TYPE_DATETIME,   false, -1},
    BindTest{tstamp,  MYSQL_TYPE_TIMESTAMP,  false, -1},

    BindTest{&Bytes,  MYSQL_TYPE_VAR_STRING, true,  -1},
    BindTest{&String, MYSQL_TYPE_STRING,     true,  -1},
    BindTest{&blob,   MYSQL_TYPE_BLOB,       true,  -1},
    BindTest{&date,   MYSQL_TYPE_DATETIME,   true,  -1},
    BindTest{&tstamp, MYSQL_TYPE_TIMESTAMP,  true,  -1},

    BindTest{pBytes,  MYSQL_TYPE_VAR_STRING, true,  -1},
    BindTest{pString, MYSQL_TYPE_STRING,     true,  -1},
    BindTest{pBlob,   MYSQL_TYPE_BLOB,       true,  -1},
    BindTest{pDate,   MYSQL_TYPE_DATETIME,   true,  -1},
    BindTest{pTstamp, MYSQL_TYPE_TIMESTAMP,  true,  -1},

    BindTest{raw,     MYSQL_TYPE_INT24,    true,  -1},

    BindTest{Int8,    MYSQL_TYPE_TINY,     false,  1},
    BindTest{Int16,   MYSQL_TYPE_SHORT,    false,  2},
    BindTest{Int32,   MYSQL_TYPE_LONG,     false,  4},
    BindTest{Int64,   MYSQL_TYPE_LONGLONG, false,  8},
    BindTest{Int,     MYSQL_TYPE_LONG,     false,  4}, // Hack

    BindTest{&Int8,   MYSQL_TYPE_TINY,     true,   1},
    BindTest{&Int16,  MYSQL_TYPE_SHORT,    true,   2},
    BindTest{&Int32,  MYSQL_TYPE_LONG,     true,   4},
    BindTest{&Int64,  MYSQL_TYPE_LONGLONG, true,   8},
    BindTest{&Int,    MYSQL_TYPE_LONG,     true,   4}, // Hack

    BindTest{pInt8,   MYSQL_TYPE_TINY,     true,   1},
    BindTest{pInt16,  MYSQL_TYPE_SHORT,    true,   2},
    BindTest{pInt32,  MYSQL_TYPE_LONG,     true,   4},
    BindTest{pInt64,  MYSQL_TYPE_LONGLONG, true,   8},
    BindTest{pInt,    MYSQL_TYPE_LONG,     true,   4}, // Hack

    BindTest{Uint8,   MYSQL_TYPE_TINY     | MYSQL_UNSIGNED_MASK, false, 1},
    BindTest{Uint16,  MYSQL_TYPE_SHORT    | MYSQL_UNSIGNED_MASK, false, 2},
    BindTest{Uint32,  MYSQL_TYPE_LONG     | MYSQL_UNSIGNED_MASK, false, 4},
    BindTest{Uint64,  MYSQL_TYPE_LONGLONG | MYSQL_UNSIGNED_MASK, false, 8},
    BindTest{Uint,    MYSQL_TYPE_LONG     | MYSQL_UNSIGNED_MASK, false,4},//Hack

    BindTest{&Uint8,  MYSQL_TYPE_TINY     | MYSQL_UNSIGNED_MASK, true, 1},
    BindTest{&Uint16, MYSQL_TYPE_SHORT    | MYSQL_UNSIGNED_MASK, true, 2},
    BindTest{&Uint32, MYSQL_TYPE_LONG     | MYSQL_UNSIGNED_MASK, true, 4},
    BindTest{&Uint64, MYSQL_TYPE_LONGLONG | MYSQL_UNSIGNED_MASK, true, 8},
    BindTest{&Uint,   MYSQL_TYPE_LONG     | MYSQL_UNSIGNED_MASK, true, 4},//Hack

    BindTest{pUint8,  MYSQL_TYPE_TINY     | MYSQL_UNSIGNED_MASK, true, 1},
    BindTest{pUint16, MYSQL_TYPE_SHORT    | MYSQL_UNSIGNED_MASK, true, 2},
    BindTest{pUint32, MYSQL_TYPE_LONG     | MYSQL_UNSIGNED_MASK, true, 4},
    BindTest{pUint64, MYSQL_TYPE_LONGLONG | MYSQL_UNSIGNED_MASK, true, 8},
    BindTest{pUint,   MYSQL_TYPE_LONG     | MYSQL_UNSIGNED_MASK, true, 4},//Hack

    BindTest{Float32, MYSQL_TYPE_FLOAT,   false, 4},
    BindTest{Float64, MYSQL_TYPE_DOUBLE,  false, 8},
    BindTest{Float,   MYSQL_TYPE_FLOAT,   false, 4}, // Hack

    BindTest{&Float32, MYSQL_TYPE_FLOAT,  true,  4},
    BindTest{&Float64, MYSQL_TYPE_DOUBLE, true,  8},
    BindTest{&Float,   MYSQL_TYPE_FLOAT,  true,  4}, // Hack
}

func TestBind(t *testing.T) {
    for _, test := range bindTests {
        val := bindValue(reflect.NewValue(test.val))
        if val.typ != test.typ || val.is_ptr != test.is_ptr ||
                val.length != test.length {
            t.Errorf(
                "Type: %s exp=0x%x res=0x%x IsPtr: exp=%t res=%t " +
                "Len: exp=%d res=%d", reflect.Typeof(test.val), test.typ,
                val.typ, test.is_ptr, val.is_ptr, test.length, val.length,
            )
        }
    }
}

type WriteTest struct {
    val interface{}
    exp []byte
}

var writeTest []WriteTest
func init() {
    b := make([]byte, 64 * 1024)
    for ii := range b {
        b[ii] = byte(ii)
    }
    blob = Blob(b)

    writeTest = []WriteTest{
        WriteTest{Bytes,  append([]byte{byte(len(Bytes))}, Bytes...)},
        WriteTest{String, append([]byte{byte(len(String))}, []byte(String)...)},
        WriteTest{pBytes,  nil},
        WriteTest{pString,  nil},
        WriteTest {
            blob,
            append(
                append(
                    []byte{253},
                    *EncodeU24(uint32(len(blob)))...
                ),
                []byte(blob)...,
            ),
        },
        WriteTest {
            date,
            []byte{
                7, byte(date.Year), byte(date.Year >> 8), byte(date.Month),
                byte(date.Day), byte(date.Hour), byte(date.Minute),
                byte(date.Second),
            },
        },
        WriteTest {
            &date,
            []byte{
                7, byte(date.Year), byte(date.Year >> 8), byte(date.Month),
                byte(date.Day), byte(date.Hour), byte(date.Minute),
                byte(date.Second),
            },
        },
        WriteTest{date,  *EncodeDatetime(&date)},
        WriteTest{&date, *EncodeDatetime(&date)},
        WriteTest{pDate, nil},

        WriteTest{tstamp,  *EncodeDatetime((*Datetime)(&tstamp))},
        WriteTest{&tstamp, *EncodeDatetime((*Datetime)(&tstamp))},
        WriteTest{pTstamp, nil},

        WriteTest{Int,     *EncodeU32(uint32(Int))}, // Hack
        WriteTest{Int16,   *EncodeU16(uint16(Int16))},
        WriteTest{Int32,   *EncodeU32(uint32(Int32))},
        WriteTest{Int64,   *EncodeU64(uint64(Int64))},

        WriteTest{Int   ,  *EncodeU32(uint32(Int))}, // Hack
        WriteTest{Uint16,  *EncodeU16(Uint16)},
        WriteTest{Uint32,  *EncodeU32(Uint32)},
        WriteTest{Uint64,  *EncodeU64(Uint64)},

        WriteTest{&Int,    *EncodeU32(uint32(Int))}, // Hack
        WriteTest{&Int16,  *EncodeU16(uint16(Int16))},
        WriteTest{&Int32,  *EncodeU32(uint32(Int32))},
        WriteTest{&Int64,  *EncodeU64(uint64(Int64))},

        WriteTest{&Uint,   *EncodeU32(uint32(Uint))}, // Hack
        WriteTest{&Uint16, *EncodeU16(Uint16)},
        WriteTest{&Uint32, *EncodeU32(Uint32)},
        WriteTest{&Uint64, *EncodeU64(Uint64)},

        WriteTest{pInt,    nil},
        WriteTest{pInt16,  nil},
        WriteTest{pInt32,  nil},
        WriteTest{pInt64,  nil},

        WriteTest{Float,   *EncodeU32(math.Float32bits(float32(Float)))},
        WriteTest{Float32, *EncodeU32(math.Float32bits(Float32))},
        WriteTest{Float64, *EncodeU64(math.Float64bits(Float64))},

        WriteTest{&Float,   *EncodeU32(math.Float32bits(float32(Float)))},
        WriteTest{&Float32, *EncodeU32(math.Float32bits(Float32))},
        WriteTest{&Float64, *EncodeU64(math.Float64bits(Float64))},

        WriteTest{pFloat,   nil},
        WriteTest{pFloat32, nil},
        WriteTest{pFloat64, nil},
    }
}

func TestWrite(t *testing.T) {
    buf := new(bytes.Buffer)
    for _, test := range writeTest {
        buf.Reset()
        val := bindValue(reflect.NewValue(test.val))
        writeValue(buf, val)
        if !bytes.Equal(buf.Bytes(), test.exp) || val.Len() != len(test.exp) {
            t.Errorf("%s - exp_len=%d res_len=%d exp: %v res: %v",
                reflect.Typeof(test.val), len(test.exp), val.Len(),
                test.exp, buf.Bytes(),
            )
        }
    }
}