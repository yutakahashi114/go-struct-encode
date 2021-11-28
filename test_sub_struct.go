package main

import (
	"encoding/binary"
	"time"
)

type TestSubStruct struct {
	Str    string
	Bool   bool
	Int    int
	Int16  int16
	Int64  int64
	Uint   uint
	Uint8  uint8
	Uint32 uint32
	Time   time.Time
}

func (s *TestSubStruct) Size() int {
	size := 0
	if s == nil {
		return 0
	}

	// Str
	size += binary.MaxVarintLen64
	size += len(s.Str)
	// Bool
	size += VarintLenBool
	// Int
	size += binary.MaxVarintLen64
	// Int16
	size += binary.MaxVarintLen16
	// Int64
	size += binary.MaxVarintLen64
	// Uint
	size += binary.MaxVarintLen64
	// Uint8
	size += MaxVarintLen8
	// Uint32
	size += binary.MaxVarintLen32
	// Time
	size += VarintLenTime
	return size
}

func (s TestSubStruct) EncodeWithBytesTime(out []byte) (int, error) {
	n := 0
	// Str
	strSize := len(s.Str)
	n += binary.PutUvarint(out[n:], uint64(strSize))
	copy(out[n:n+strSize], s.Str)
	n += strSize
	// Bool
	if s.Bool {
		n += binary.PutUvarint(out[n:], uint64(1))
	} else {
		n += binary.PutUvarint(out[n:], uint64(0))
	}
	// Int
	n += binary.PutVarint(out[n:], int64(s.Int))
	// Int16
	n += binary.PutVarint(out[n:], int64(s.Int16))
	// Int64
	n += binary.PutVarint(out[n:], s.Int64)
	// Uint
	n += binary.PutUvarint(out[n:], uint64(s.Uint))
	// Uint8
	n += binary.PutUvarint(out[n:], uint64(s.Uint8))
	// Uint32
	n += binary.PutUvarint(out[n:], uint64(s.Uint32))
	// Time
	timeLen, err := TimeMarshalBinary(s.Time, out[n:])
	if err != nil {
		return 0, err
	}
	n += timeLen

	return n, nil
}

func (s TestSubStruct) EncodeWithBytes(out []byte) (int, error) {
	n := 0
	// Str
	strSize := len(s.Str)
	n += binary.PutUvarint(out[n:], uint64(strSize))
	copy(out[n:n+strSize], s.Str)
	n += strSize
	// Bool
	if s.Bool {
		n += binary.PutUvarint(out[n:], uint64(1))
	} else {
		n += binary.PutUvarint(out[n:], uint64(0))
	}
	// Int
	n += binary.PutVarint(out[n:], int64(s.Int))
	// Int16
	n += binary.PutVarint(out[n:], int64(s.Int16))
	// Int64
	n += binary.PutVarint(out[n:], s.Int64)
	// Uint
	n += binary.PutUvarint(out[n:], uint64(s.Uint))
	// Uint8
	n += binary.PutUvarint(out[n:], uint64(s.Uint8))
	// Uint32
	n += binary.PutUvarint(out[n:], uint64(s.Uint32))
	// Time
	timeBytes, err := s.Time.MarshalBinary()
	if err != nil {
		return 0, err
	}
	copy(out[n:n+VarintLenTime], timeBytes)
	n += VarintLenTime

	return n, nil
}

func (s TestSubStruct) Encode() ([]byte, error) {
	// encode後の最大サイズを計算する
	size := 0
	// Str
	size += binary.MaxVarintLen64
	size += len(s.Str)
	// Bool
	size += VarintLenBool
	// Int
	size += binary.MaxVarintLen64
	// Int16
	size += binary.MaxVarintLen16
	// Int64
	size += binary.MaxVarintLen64
	// Uint
	size += binary.MaxVarintLen64
	// Uint8
	size += MaxVarintLen8
	// Uint32
	size += binary.MaxVarintLen32
	// Time
	size += VarintLenTime

	// 最大サイズ分の大きさを確保
	out := make([]byte, size)

	n := 0
	// Str
	strSize := len(s.Str)
	n += binary.PutUvarint(out[n:], uint64(strSize))
	copy(out[n:n+strSize], s.Str)
	n += strSize
	// Bool
	if s.Bool {
		n += binary.PutUvarint(out[n:], uint64(1))
	} else {
		n += binary.PutUvarint(out[n:], uint64(0))
	}
	// Int
	n += binary.PutVarint(out[n:], int64(s.Int))
	// Int16
	n += binary.PutVarint(out[n:], int64(s.Int16))
	// Int64
	n += binary.PutVarint(out[n:], s.Int64)
	// Uint
	n += binary.PutUvarint(out[n:], uint64(s.Uint))
	// Uint8
	n += binary.PutUvarint(out[n:], uint64(s.Uint8))
	// Uint32
	n += binary.PutUvarint(out[n:], uint64(s.Uint32))
	// Time
	timeBytes, err := s.Time.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(out[n:n+VarintLenTime], timeBytes)
	n += VarintLenTime

	return out[:n], nil
}

func (s *TestSubStruct) Decode(in []byte) (int, error) {
	n := 0

	// Str
	strLen, strLenLen := binary.Uvarint(in)
	n += strLenLen
	s.Str = string(in[n : n+int(strLen)])
	n += int(strLen)
	// Bool
	boolRaw, boolLen := binary.Uvarint(in[n:])
	s.Bool = boolRaw == 1
	n += boolLen
	// Int
	intRaw, intLen := binary.Varint(in[n:])
	s.Int = int(intRaw)
	n += intLen
	// Int16
	int16Raw, int16Len := binary.Varint(in[n:])
	s.Int16 = int16(int16Raw)
	n += int16Len
	// Int64
	int64Raw, int64Len := binary.Varint(in[n:])
	s.Int64 = int64(int64Raw)
	n += int64Len
	// Uint
	uintRaw, uintLen := binary.Uvarint(in[n:])
	s.Uint = uint(uintRaw)
	n += uintLen
	// Uint8
	uint8Raw, uint8Len := binary.Uvarint(in[n:])
	s.Uint8 = uint8(uint8Raw)
	n += uint8Len
	// Uint32
	uint32Raw, uint32Len := binary.Uvarint(in[n:])
	s.Uint32 = uint32(uint32Raw)
	n += uint32Len
	// Time
	err := s.Time.UnmarshalBinary(in[n : n+VarintLenTime])
	if err != nil {
		return 0, err
	}
	n += VarintLenTime
	return n, nil
}

type TestSubStructs []TestSubStruct

func (ss TestSubStructs) Size() int {
	size := 0

	size += VarintLenPointer
	if ss == nil {
		return size
	}

	// スライスの長さのサイズ
	size += binary.MaxVarintLen64
	// スライスの要素のサイズ
	for _, s := range ss {
		size += s.Size()
	}

	return size
}

func (ss *TestSubStructs) Decode(in []byte) (int, error) {
	n := 0
	isNotNil, isNotNilLen := binary.Uvarint(in)
	if isNotNil == 0 {
		*ss = nil
		return isNotNilLen, nil
	}
	n += isNotNilLen

	ssLen, ssLenLen := binary.Varint(in[n:])
	n += ssLenLen
	ssLenInt := int(ssLen)
	*ss = make(TestSubStructs, ssLenInt)
	for i := 0; i < ssLenInt; i++ {
		sLen, err := (*ss)[i].Decode(in[n:])
		if err != nil {
			return 0, nil
		}
		n += sLen
	}
	return n, nil
}
