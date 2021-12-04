package main

import (
	"bytes"
	"encode/proto"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gogo/protobuf/types"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

const (
	MaxVarintLen8    = 2
	VarintLenBool    = 1
	VarintLenTime    = 15
	VarintLenPointer = 1
)

func main() {
	{
		i := 12345678
		bytes, nEn := IntEncode(i)
		decoded, nDe := IntDecode(bytes)
		fmt.Println(i, nEn, decoded, nDe)
		printDiff(i, decoded)
	}
	{
		u := uint(12345678)
		bytes, nEn := UintEncode(u)
		decoded, nDe := UintDecode(bytes)
		fmt.Println(u, nEn, decoded, nDe)
		printDiff(u, decoded)
	}
	{
		str := "test_string"
		bytes, nEn := StringEncode(str)
		decoded, nDe := StringDecode(bytes)
		fmt.Println(str, nEn, decoded, nDe)
		printDiff(str, decoded)
	}
	{
		b := true
		bytes, nEn := BoolEncode(b)
		decoded, nDe := BoolDecode(bytes)
		fmt.Println(b, nEn, decoded, nDe)
		printDiff(b, decoded)
	}
	{
		b := false
		bytes, nEn := BoolEncode(b)
		decoded, nDe := BoolDecode(bytes)
		fmt.Println(b, nEn, decoded, nDe)
		printDiff(b, decoded)
	}
	{
		f := 1234.5678
		bytes, nEn := FloatEncode(f)
		decoded, nDe := FloatDecode(bytes)
		fmt.Println(f, nEn, decoded, nDe)
		printDiff(f, decoded)
	}
	{
		p := &[]int{100}[0]
		bytes, nEn := PointerEncode(p)
		decoded, nDe := PointerDecode(bytes)
		fmt.Println(*p, nEn, *decoded, nDe)
		printDiff(p, decoded)
	}
	{
		var p *int
		bytes, nEn := PointerEncode(p)
		decoded, nDe := PointerDecode(bytes)
		fmt.Println(p, nEn, decoded, nDe)
		printDiff(p, decoded)
	}
	{
		slice := []int{1, 1000000000000000000, -1000000000000000000}
		bytes, nEn := SliceEncode(slice)
		decoded, nDe := SliceDecode(bytes)
		fmt.Println(slice, nEn, decoded, nDe)
		printDiff(slice, decoded)
	}
	{
		slice := []int{}
		bytes, nEn := SliceEncode(slice)
		decoded, nDe := SliceDecode(bytes)
		fmt.Println(slice, nEn, decoded, nDe)
		printDiff(slice, decoded)
	}
	{
		var slice []int
		bytes, nEn := SliceEncode(slice)
		decoded, nDe := SliceDecode(bytes)
		fmt.Println(slice, nEn, decoded, nDe)
		printDiff(slice, decoded)
	}
	{
		data := TestSubStruct{
			Str:    "test_string",
			Bool:   true,
			Int:    1,
			Int16:  10000,
			Int64:  1000000000000000000,
			Uint:   1,
			Uint8:  100,
			Uint32: 1000000000,
			Time:   time.Now(),
		}
		bytes, _ := data.Encode()
		fmt.Println(len(bytes))

		decoded := TestSubStruct{}
		decoded.Decode(bytes)

		if diff := cmp.Diff(data, decoded); diff != "" {
			fmt.Println(diff)
		} else {
			fmt.Println("no diff")
		}
	}
	data := createTestStructs(10000)
	{
		bytes, err := json.Marshal(data)
		fataiIf(err)

		fmt.Println(len(bytes))

		decoded := TestStructs{}
		err = json.Unmarshal(bytes, &decoded)
		fataiIf(err)

		if diff := cmp.Diff(data, decoded); diff != "" {
			fmt.Println(diff)
		}
	}
	{
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(&data)
		fataiIf(err)
		byt := buf.Bytes()
		fmt.Println(len(byt))

		decoded := TestStructs{}
		buf = bytes.NewBuffer(byt)
		err = gob.NewDecoder(buf).Decode(&decoded)
		fataiIf(err)
		if diff := cmp.Diff(data, decoded); diff != "" {
			fmt.Println(diff)
		}
	}
	{
		bytes, err := data.Encode()
		fataiIf(err)

		fmt.Println(len(bytes))

		decoded := TestStructs{}
		_, err = decoded.Decode(bytes)
		fataiIf(err)

		if diff := cmp.Diff(data, decoded); diff != "" {
			fmt.Println(diff)
		}
	}
	{
		bytes, err := data.EncodeTime()
		fataiIf(err)

		fmt.Println(len(bytes))

		decoded := TestStructs{}
		_, err = decoded.Decode(bytes)
		fataiIf(err)

		if diff := cmp.Diff(data, decoded); diff != "" {
			fmt.Println(diff)
		}
	}
	{
		data := makeProtoTestStructs(data)
		bytes, err := protobuf.Marshal(data)
		fataiIf(err)

		fmt.Println(len(bytes))

		decoded := &proto.TestStructs{}
		err = protobuf.Unmarshal(bytes, decoded)
		fataiIf(err)

		if diff := cmp.Diff(data, decoded); diff != "" {
			fmt.Println(diff)
		}
	}
}

func fataiIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printDiff(x interface{}, y interface{}) {
	if diff := cmp.Diff(x, y); diff != "" {
		fmt.Println(diff)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(rand.Int63()%int64(len(letters)))]
	}
	return string(b)
}

func createTestStructs(size int) TestStructs {
	ss := make(TestStructs, size)
	for i := range ss {
		subs := make(TestSubStructs, 10)
		for j := range subs {
			subs[j] = createTestSubStruct()
		}
		sub := createTestSubStruct()
		ss[i] = TestStruct{
			Str:        randString(10),
			Bool:       rand.Int()%2 == 1,
			Int:        rand.Int(),
			Int16:      int16(rand.Int()),
			Int64:      int64(rand.Int()),
			Uint:       uint(rand.Uint64()),
			Uint8:      uint8(rand.Uint64()),
			Uint32:     uint32(rand.Uint64()),
			Time:       time.Now(),
			SubPointer: &sub,
			Subs:       subs,
		}
	}
	return ss
}

func createTestSubStruct() TestSubStruct {
	return TestSubStruct{
		Str:    randString(10),
		Bool:   rand.Int()%2 == 1,
		Int:    rand.Int(),
		Int16:  int16(rand.Int()),
		Int64:  int64(rand.Int()),
		Uint:   uint(rand.Uint64()),
		Uint8:  uint8(rand.Uint64()),
		Uint32: uint32(rand.Uint64()),
		Time:   time.Now(),
	}
}

var timeZero = time.Time{}.Unix()

func TimeMarshalBinary(t time.Time, out []byte) (int, error) {
	var offsetMin int16 // minutes east of UTC. -1 is UTC.

	if t.Location() == time.UTC {
		offsetMin = -1
	} else {
		_, offset := t.Zone()
		if offset%60 != 0 {
			return 0, errors.New("Time.MarshalBinary: zone offset has fractional minute")
		}
		offset /= 60
		if offset < -32768 || offset == -1 || offset > 32767 {
			return 0, errors.New("Time.MarshalBinary: unexpected zone offset")
		}
		offsetMin = int16(offset)
	}

	unix := t.Unix()
	sec := unix - timeZero
	nsec := t.UnixNano() - unix*1000000000
	out[0] = 1               // byte 0 : version
	out[1] = byte(sec >> 56) // bytes 1-8: seconds
	out[2] = byte(sec >> 48)
	out[3] = byte(sec >> 40)
	out[4] = byte(sec >> 32)
	out[5] = byte(sec >> 24)
	out[6] = byte(sec >> 16)
	out[7] = byte(sec >> 8)
	out[8] = byte(sec)
	out[9] = byte(nsec >> 24) // bytes 9-12: nanoseconds
	out[10] = byte(nsec >> 16)
	out[11] = byte(nsec >> 8)
	out[12] = byte(nsec)
	out[13] = byte(offsetMin >> 8) // bytes 13-14: zone offset in minutes
	out[14] = byte(offsetMin)

	return VarintLenTime, nil
}

func IntEncode(i int) ([]byte, int) {
	// intの変換に必要な最大容量を確保
	out := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(out, int64(i))
	return out[:n], n
}

func IntDecode(in []byte) (int, int) {
	intRaw, intLen := binary.Varint(in)
	return int(intRaw), intLen
}

func UintEncode(u uint) ([]byte, int) {
	// uintの変換に必要な最大容量を確保
	out := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(out, uint64(u))
	return out[:n], n
}

func UintDecode(in []byte) (uint, int) {
	uintRaw, uintLen := binary.Uvarint(in)
	return uint(uintRaw), uintLen
}

func StringEncode(str string) ([]byte, int) {
	n := 0
	strSize := len(str)
	// 文字列の長さの最大 + 文字列の長さ
	out := make([]byte, binary.MaxVarintLen64+strSize)
	// 文字列の長さをセット
	n += binary.PutUvarint(out, uint64(strSize))
	// 文字列をコピー
	copy(out[n:n+strSize], str)
	n += strSize
	return out[:n], n
}

func StringDecode(in []byte) (string, int) {
	n := 0
	// 文字列の長さを読み取る
	strLen, strLenLen := binary.Uvarint(in)
	n += strLenLen
	// 長さの分だけ文字列として読み取る
	str := string(in[n : n+int(strLen)])
	n += int(strLen)
	return str, n
}

func BoolEncode(b bool) ([]byte, int) {
	if b {
		return []byte{1}, 1
	}
	return []byte{0}, 1
}

func BoolDecode(in []byte) (bool, int) {
	return in[0] == 1, 1
}

func FloatEncode(f float64) ([]byte, int) {
	out := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(out, math.Float64bits(f))
	return out[:n], n
}

func FloatDecode(in []byte) (float64, int) {
	floatRaw, floatLen := binary.Uvarint(in)
	return math.Float64frombits(floatRaw), floatLen
}

func PointerEncode(p *int) ([]byte, int) {
	n := 1
	if p == nil {
		// nilの場合は1バイト目に0をセット
		return []byte{0}, n
	}
	// nil判定 + intの最大
	out := make([]byte, n+binary.MaxVarintLen64)
	// nilでない場合は1バイト目に1をセット
	out[0] = 1
	n += binary.PutVarint(out[n:], int64(*p))
	return out[:n], n
}

func PointerDecode(in []byte) (*int, int) {
	n := 0
	// 1バイト目が0ならnilを返す
	isNotNil, isNotNilLen := binary.Uvarint(in[n:])
	n += isNotNilLen
	if isNotNil == 0 {
		return nil, n
	}
	intRaw, intLen := binary.Varint(in[n:])
	n += intLen
	t := int(intRaw)
	return &t, n
}

func SliceEncode(ints []int) ([]byte, int) {
	n := 1
	if ints == nil {
		// nilの場合は1バイト目に0をセット
		return []byte{0}, n
	}
	intsSize := len(ints)
	// nil判定 + スライスの長さの最大 + スライスの要素数*1要素の最大
	out := make([]byte, n+binary.MaxVarintLen64+intsSize*binary.MaxVarintLen64)
	// nilでない場合は1バイト目に1をセット
	out[0] = 1
	// スライスの長さをセット
	n += binary.PutUvarint(out[n:], uint64(intsSize))
	// スライスを1要素ずつセット
	for _, i := range ints {
		n += binary.PutVarint(out[n:], int64(i))
	}
	return out[:n], n
}

func SliceDecode(in []byte) ([]int, int) {
	n := 0
	// 1バイト目が0ならnilを返す
	sliceIsNotNil, sliceIsNotNilLen := binary.Uvarint(in[n:])
	n += sliceIsNotNilLen
	if sliceIsNotNil == 0 {
		return nil, n
	}
	// スライスの長さを読み取る
	sliceLen, sliceLenLen := binary.Uvarint(in[n:])
	n += sliceLenLen
	slice := make([]int, sliceLen)
	// 長さの回数だけ読み取る
	for i := 0; i < int(sliceLen); i++ {
		intRaw, intLen := binary.Varint(in[n:])
		slice[i] = int(intRaw)
		n += intLen
	}
	return slice, n
}

func makeProtoTestStructs(ss TestStructs) *proto.TestStructs {
	ssProto := make([]*proto.TestStruct, len(ss))
	for i, s := range ss {
		subs := make([]*proto.TestSubStruct, len(s.Subs))
		for j, sub := range s.Subs {
			subs[j] = makeProtoTestSubStruct(&sub)
		}
		ts, err := types.TimestampProto(s.Time)
		fataiIf(err)
		ssProto[i] = &proto.TestStruct{
			Str:        s.Str,
			Bool:       s.Bool,
			Int:        int64(s.Int),
			Int16:      int32(s.Int16),
			Int64:      s.Int64,
			Uint:       uint64(s.Uint),
			Uint8:      uint32(s.Uint8),
			Uint32:     s.Uint32,
			Time:       ts,
			SubPointer: makeProtoTestSubStruct(s.SubPointer),
			Subs:       subs,
		}
	}
	return &proto.TestStructs{
		Ss: ssProto,
	}
}

func makeProtoTestSubStruct(sub *TestSubStruct) *proto.TestSubStruct {
	ts, err := types.TimestampProto(sub.Time)
	fataiIf(err)
	return &proto.TestSubStruct{
		Str:    sub.Str,
		Bool:   sub.Bool,
		Int:    int64(sub.Int),
		Int16:  int32(sub.Int16),
		Int64:  sub.Int64,
		Uint:   uint64(sub.Uint),
		Uint8:  uint32(sub.Uint8),
		Uint32: sub.Uint32,
		Time:   ts,
	}
}
