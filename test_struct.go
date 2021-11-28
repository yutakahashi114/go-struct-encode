package main

import (
	"encoding/binary"
	"time"
)

type TestStruct struct {
	Str        string
	Bool       bool
	Int        int
	Int16      int16
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint32     uint32
	Time       time.Time
	SubPointer *TestSubStruct
	Subs       TestSubStructs
}

func (s *TestStruct) Size() int {
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
	// SubPointer
	size += VarintLenPointer
	size += s.SubPointer.Size()
	// Subs
	size += s.Subs.Size()
	return size
}

func (s TestStruct) EncodeWithBytes(out []byte) (int, error) {
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
	// SubPointer
	if s.SubPointer == nil {
		out[n] = 0
		n += VarintLenPointer
	} else {
		out[n] = 1
		n += VarintLenPointer
		subLen, err := s.SubPointer.EncodeWithBytes(out[n:])
		if err != nil {
			return 0, err
		}
		n += subLen
	}
	// Subs
	if s.Subs == nil {
		out[n] = 0
		n += VarintLenPointer
	} else {
		out[n] = 1
		n += VarintLenPointer
		// スライスの長さ
		n += binary.PutVarint(out[n:], int64(len(s.Subs)))
		for _, s := range s.Subs {
			subLen, err := s.EncodeWithBytes(out[n:])
			if err != nil {
				return 0, err
			}
			n += subLen
		}
	}

	return n, nil
}

func (s TestStruct) EncodeWithBytesTime(out []byte) (int, error) {
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
	// SubPointer
	if s.SubPointer == nil {
		out[n] = 0
		n += VarintLenPointer
	} else {
		out[n] = 1
		n += VarintLenPointer
		subLen, err := s.SubPointer.EncodeWithBytesTime(out[n:])
		if err != nil {
			return 0, err
		}
		n += subLen
	}
	// Subs
	if s.Subs == nil {
		out[n] = 0
		n += VarintLenPointer
	} else {
		out[n] = 1
		n += VarintLenPointer
		// スライスの長さ
		n += binary.PutVarint(out[n:], int64(len(s.Subs)))
		for _, s := range s.Subs {
			subLen, err := s.EncodeWithBytesTime(out[n:])
			if err != nil {
				return 0, err
			}
			n += subLen
		}
	}

	return n, nil
}

func (s *TestStruct) Decode(in []byte) (int, error) {
	*s = TestStruct{}
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
	// SubPointer
	subPointerIsNotNil, subPointerIsNotNilLen := binary.Uvarint(in[n:])
	n += subPointerIsNotNilLen
	if subPointerIsNotNil == 1 {
		s.SubPointer = &TestSubStruct{}
		subPointerLen, err := s.SubPointer.Decode(in[n:])
		if err != nil {
			return 0, err
		}
		n += subPointerLen
	}
	// Subs
	subsLen, err := s.Subs.Decode(in[n:])
	if err != nil {
		return 0, err
	}
	n += subsLen

	return n, nil
}

type TestStructs []TestStruct

func (ss TestStructs) Encode() ([]byte, error) {
	size := VarintLenPointer

	if ss == nil {
		// nil
		return []byte{0}, nil
	}

	size += binary.MaxVarintLen64
	ssLen := len(ss)
	for _, s := range ss {
		size += s.Size()
	}

	out := make([]byte, size)
	n := 0
	// nilでない
	n += binary.PutUvarint(out[n:], uint64(1))
	// スライスの長さ
	n += binary.PutVarint(out[n:], int64(ssLen))
	for _, s := range ss {
		bytesLen, err := s.EncodeWithBytes(out[n:])
		if err != nil {
			return nil, err
		}
		n += bytesLen
	}

	return out[:n], nil
}

func (ss TestStructs) EncodeTime() ([]byte, error) {
	size := VarintLenPointer

	if ss == nil {
		// nil
		return []byte{0}, nil
	}

	size += binary.MaxVarintLen64
	ssLen := len(ss)
	for _, s := range ss {
		size += s.Size()
	}

	out := make([]byte, size)
	n := 0
	// nilでない
	n += binary.PutUvarint(out[n:], uint64(1))
	// スライスの長さ
	n += binary.PutVarint(out[n:], int64(ssLen))
	for _, s := range ss {
		bytesLen, err := s.EncodeWithBytesTime(out[n:])
		if err != nil {
			return nil, err
		}
		n += bytesLen
	}

	return out[:n], nil
}

func (ss *TestStructs) Decode(in []byte) (int, error) {
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
	*ss = make(TestStructs, ssLenInt)
	for i := 0; i < ssLenInt; i++ {
		sLen, err := (*ss)[i].Decode(in[n:])
		if err != nil {
			return 0, nil
		}
		n += sLen
	}
	return n, nil
}
