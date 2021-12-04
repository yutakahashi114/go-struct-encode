package main

import (
	"bytes"
	"encode/proto"
	"encoding/gob"
	"encoding/json"
	"testing"

	protobuf "github.com/golang/protobuf/proto"
)

var testStructsMap = map[int]TestStructs{
	1:     createTestStructs(1),
	10:    createTestStructs(10),
	100:   createTestStructs(100),
	1000:  createTestStructs(1000),
	10000: createTestStructs(10000),
}

var testStructsProtoMap = map[int]*proto.TestStructs{
	1:     makeProtoTestStructs(testStructsMap[1]),
	10:    makeProtoTestStructs(testStructsMap[10]),
	100:   makeProtoTestStructs(testStructsMap[100]),
	1000:  makeProtoTestStructs(testStructsMap[1000]),
	10000: makeProtoTestStructs(testStructsMap[10000]),
}

func encodeBase(b *testing.B, sliceSize int, encodeFn func(TestStructs) ([]byte, error)) {
	ss := testStructsMap[sliceSize]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encodeFn(ss)
		if err != nil {
			panic(err)
		}
	}
}

func encodeJson(ss TestStructs) ([]byte, error) {
	return json.Marshal(ss)
}

func encodeGob(ss TestStructs) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(&ss)
	bytes := buf.Bytes()
	return bytes, err
}

func encodeSelf(ss TestStructs) ([]byte, error) {
	return ss.Encode()
}

func encodeSelfTime(ss TestStructs) ([]byte, error) {
	return ss.EncodeTime()
}

func encodeProto(b *testing.B, sliceSize int) {
	ss := testStructsProtoMap[sliceSize]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := protobuf.Marshal(ss)
		if err != nil {
			panic(err)
		}
	}
}

func decodeBase(b *testing.B, sliceSize int, encodeFn func(TestStructs) ([]byte, error), decodeFn func([]byte) (TestStructs, error)) {
	ss := testStructsMap[sliceSize]
	bs, err := encodeFn(ss)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := decodeFn(bs)
		if err != nil {
			panic(err)
		}
	}
}

func decodeJson(bs []byte) (TestStructs, error) {
	decoded := TestStructs{}
	err := json.Unmarshal(bs, &decoded)
	return decoded, err
}

func decodeGob(bs []byte) (TestStructs, error) {
	decoded := TestStructs{}
	buf := bytes.NewBuffer(bs)
	err := gob.NewDecoder(buf).Decode(&decoded)
	return decoded, err
}

func decodeSelf(bs []byte) (TestStructs, error) {
	decoded := TestStructs{}
	_, err := decoded.Decode(bs)
	return decoded, err
}

func decodeSelfTime(bs []byte) (TestStructs, error) {
	decoded := TestStructs{}
	_, err := decoded.Decode(bs)
	return decoded, err
}

func decodeProto(b *testing.B, sliceSize int) {
	ss := testStructsProtoMap[sliceSize]
	bs, err := protobuf.Marshal(ss)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded := &proto.TestStructs{}
		err = protobuf.Unmarshal(bs, decoded)
		if err != nil {
			panic(err)
		}
	}
}

func Benchmark_encode_____json_____1(b *testing.B) {
	encodeBase(b, 1, encodeJson)
}

func Benchmark_encode______gob_____1(b *testing.B) {
	encodeBase(b, 1, encodeGob)
}

func Benchmark_encode_____self_____1(b *testing.B) {
	encodeBase(b, 1, encodeSelf)
}

func Benchmark_encode_selftime_____1(b *testing.B) {
	encodeBase(b, 1, encodeSelfTime)
}

func Benchmark_encode_protobuf_____1(b *testing.B) {
	encodeProto(b, 1)
}

func Benchmark_encode_____json____10(b *testing.B) {
	encodeBase(b, 10, encodeJson)
}

func Benchmark_encode______gob____10(b *testing.B) {
	encodeBase(b, 10, encodeGob)
}

func Benchmark_encode_____self____10(b *testing.B) {
	encodeBase(b, 10, encodeSelf)
}

func Benchmark_encode_selftime____10(b *testing.B) {
	encodeBase(b, 10, encodeSelfTime)
}

func Benchmark_encode_protobuf____10(b *testing.B) {
	encodeProto(b, 10)
}

func Benchmark_encode_____json___100(b *testing.B) {
	encodeBase(b, 100, encodeJson)
}

func Benchmark_encode______gob___100(b *testing.B) {
	encodeBase(b, 100, encodeGob)
}

func Benchmark_encode_____self___100(b *testing.B) {
	encodeBase(b, 100, encodeSelf)
}

func Benchmark_encode_selftime___100(b *testing.B) {
	encodeBase(b, 100, encodeSelfTime)
}

func Benchmark_encode_protobuf___100(b *testing.B) {
	encodeProto(b, 100)
}

func Benchmark_encode_____json__1000(b *testing.B) {
	encodeBase(b, 1000, encodeJson)
}

func Benchmark_encode______gob__1000(b *testing.B) {
	encodeBase(b, 1000, encodeGob)
}

func Benchmark_encode_____self__1000(b *testing.B) {
	encodeBase(b, 1000, encodeSelf)
}

func Benchmark_encode_selftime__1000(b *testing.B) {
	encodeBase(b, 1000, encodeSelfTime)
}

func Benchmark_encode_protobuf__1000(b *testing.B) {
	encodeProto(b, 1000)
}

func Benchmark_encode_____json_10000(b *testing.B) {
	encodeBase(b, 10000, encodeJson)
}

func Benchmark_encode______gob_10000(b *testing.B) {
	encodeBase(b, 10000, encodeGob)
}

func Benchmark_encode_____self_10000(b *testing.B) {
	encodeBase(b, 10000, encodeSelf)
}

func Benchmark_encode_selftime_10000(b *testing.B) {
	encodeBase(b, 10000, encodeSelfTime)
}

func Benchmark_encode_protobuf_10000(b *testing.B) {
	encodeProto(b, 10000)
}

func Benchmark_decode_____json_____1(b *testing.B) {
	decodeBase(b, 1, encodeJson, decodeJson)
}

func Benchmark_decode______gob_____1(b *testing.B) {
	decodeBase(b, 1, encodeGob, decodeGob)
}

func Benchmark_decode_____self_____1(b *testing.B) {
	decodeBase(b, 1, encodeSelf, decodeSelf)
}

func Benchmark_decode_selftime_____1(b *testing.B) {
	decodeBase(b, 1, encodeSelfTime, decodeSelfTime)
}

func Benchmark_decode_protobuf_____1(b *testing.B) {
	decodeProto(b, 1)
}

func Benchmark_decode_____json____10(b *testing.B) {
	decodeBase(b, 10, encodeJson, decodeJson)
}

func Benchmark_decode______gob____10(b *testing.B) {
	decodeBase(b, 10, encodeGob, decodeGob)
}

func Benchmark_decode_____self____10(b *testing.B) {
	decodeBase(b, 10, encodeSelf, decodeSelf)
}

func Benchmark_decode_selftime____10(b *testing.B) {
	decodeBase(b, 10, encodeSelfTime, decodeSelfTime)
}

func Benchmark_decode_protobuf____10(b *testing.B) {
	decodeProto(b, 10)
}

func Benchmark_decode_____json___100(b *testing.B) {
	decodeBase(b, 100, encodeJson, decodeJson)
}

func Benchmark_decode______gob___100(b *testing.B) {
	decodeBase(b, 100, encodeGob, decodeGob)
}

func Benchmark_decode_____self___100(b *testing.B) {
	decodeBase(b, 100, encodeSelf, decodeSelf)
}

func Benchmark_decode_selftime___100(b *testing.B) {
	decodeBase(b, 100, encodeSelfTime, decodeSelfTime)
}

func Benchmark_decode_protobuf___100(b *testing.B) {
	decodeProto(b, 100)
}

func Benchmark_decode_____json__1000(b *testing.B) {
	decodeBase(b, 1000, encodeJson, decodeJson)
}

func Benchmark_decode______gob__1000(b *testing.B) {
	decodeBase(b, 1000, encodeGob, decodeGob)
}

func Benchmark_decode_____self__1000(b *testing.B) {
	decodeBase(b, 1000, encodeSelf, decodeSelf)
}

func Benchmark_decode_selftime__1000(b *testing.B) {
	decodeBase(b, 1000, encodeSelfTime, decodeSelfTime)
}

func Benchmark_decode_protobuf__1000(b *testing.B) {
	decodeProto(b, 1000)
}

func Benchmark_decode_____json_10000(b *testing.B) {
	decodeBase(b, 10000, encodeJson, decodeJson)
}

func Benchmark_decode______gob_10000(b *testing.B) {
	decodeBase(b, 10000, encodeGob, decodeGob)
}

func Benchmark_decode_____self_10000(b *testing.B) {
	decodeBase(b, 10000, encodeSelf, decodeSelf)
}

func Benchmark_decode_selftime_10000(b *testing.B) {
	decodeBase(b, 10000, encodeSelfTime, decodeSelfTime)
}

func Benchmark_decode_protobuf_10000(b *testing.B) {
	decodeProto(b, 10000)
}
