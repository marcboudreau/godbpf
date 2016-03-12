package qfs

import (
  "testing"
  "bytes"
)

func TestEncodeWithZeroBytes(t *testing.T) {
  buffer := new(bytes.Buffer)
  data := make([]byte, 0)

  if e := Encode(buffer, data); e != nil {
    t.Error()
  }

  expected := []byte{ 0x10, 0xFB, 0x0, 0x0, 0x0, 0xFC }

  CheckIfSlicesAreEqual(t, buffer.Bytes(), expected)
}

func TestEncodeWithSingleByte(t *testing.T) {
  buffer := new(bytes.Buffer)
  data := []byte{ 0xA5 }

  if e := Encode(buffer, data); e != nil {
    t.Error()
  }

  expected := []byte{ 0x10, 0xFB, 0x0, 0x0, 0x1, 0xFD, 0xA5 }

  CheckIfSlicesAreEqual(t, buffer.Bytes(), expected)
}

func TestEncodeWithA3ByteChain(t *testing.T) {
  buffer := new(bytes.Buffer)
  data := []byte{ 0xA5, 0x24, 0xA5, 0xA5, 0xA5 }

  if e := Encode(buffer, data); e != nil {
    t.Error()
  }

  expected := []byte{ 0x10, 0xFB, 0x0, 0x0, 0x5, 0x2, 0x1, 0xA5, 0x24, 0xFC }

  CheckIfSlicesAreEqual(t, buffer.Bytes(), expected)
}

func TestEncodeWithMultipleRepeatingChains(t *testing.T) {
  buffer := new(bytes.Buffer)
  data := []byte{ 0xA5, 0x24, 0x5C, 0x71, 0xA5, 0xA5, 0xA5, 0x2E, 0x6A, 0x71, 0x71, 0x71, 0x71, 0x88, 0x04 }

  if e := Encode(buffer, data); e != nil {
    t.Error()
  }

  expected := []byte{ 0x10, 0xFB, 0x0, 0x0, 0xF, 0xE0, 0xA5, 0x24, 0x5C, 0x71, 0x0, 0x3, 0x6, 0x5, 0x2E, 0x6A, 0xFE, 0x88, 0x04 }

  CheckIfSlicesAreEqual(t, buffer.Bytes(), expected)
}

func TestDecodeZeroBytes(t *testing.T) {
  buffer := bytes.NewBuffer([]byte{ 0x10, 0xFB, 0x0, 0x0, 0x0, 0xFC })

  data, e := Decode(buffer)
  if e != nil {
    t.Error()
  }

  expected := make([]byte, 0)

  CheckIfSlicesAreEqual(t, data, expected)
}

func TestDecodeSingleByte(t *testing.T) {
  buffer := bytes.NewBuffer([]byte{ 0x10, 0xFB, 0x0, 0x0, 0x1, 0xFD, 0x47 })

  data, e := Decode(buffer)
  if e != nil {
    t.Error()
  }

  expected := []byte{ 0x47 }

  CheckIfSlicesAreEqual(t, data, expected)
}

func TestDecodeSingle3ByteChain(t *testing.T) {
  buffer := bytes.NewBuffer([]byte{ 0x10, 0xFB, 0x0, 0x0, 0x7, 0x03, 0x02, 0x47, 0x69, 0x22, 0xFD, 0x3D })

  data, e := Decode(buffer)
  if e != nil {
    t.Error()
  }

  expected := []byte{ 0x47, 0x69, 0x22, 0x47, 0x47, 0x47, 0x3D }

  CheckIfSlicesAreEqual(t, data, expected)
}

func TestDecodeMultipleByteChains(t *testing.T) {
  buffer := bytes.NewBuffer([]byte{ 0x10, 0xFB, 0x0, 0x0, 0xF, 0xE0, 0xA5, 0x24, 0x5C, 0x71, 0x0, 0x3, 0x6, 0x5, 0x2E, 0x6A, 0xFE, 0x88, 0x04 })

  data, e := Decode(buffer)
  if e != nil {
    t.Error()
  }

  expected := []byte{ 0xA5, 0x24, 0x5C, 0x71, 0xA5, 0xA5, 0xA5, 0x2E, 0x6A, 0x71, 0x71, 0x71, 0x71, 0x88, 0x04 }

  CheckIfSlicesAreEqual(t, data, expected)
}

func CheckIfSlicesAreEqual(t *testing.T, actual []byte, expected []byte) {
  if len(actual) != len(expected) {
    t.Errorf("Actual slice size %d didn't match expected size %d", len(actual), len(expected))
  }

  for i, v := range actual {
    if v != expected[i] {
      t.Errorf("Byte %d: expected %2x, but actually was %2x", i, expected[i], v)
    }
  }
}
