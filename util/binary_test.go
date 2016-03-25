package util

import (
  "testing"
)

func TestReadUint32FourByteSlice(t *testing.T) {
  data := []byte{ 0x28, 0x34, 0x71, 0xCA }

  expected := uint32(3396416552)
  actual := ReadUint32(data)

  if actual != expected {
    t.Error()
  }
}

func TestReadUint32ThreeByteSlice(t *testing.T) {
  data := []byte{ 0x28, 0x34, 0x71 }

  expected := uint32(7418920)
  actual := ReadUint32(data)

  if actual != expected {
    t.Error()
  }
}

func TestReadUint32ZeroByteSlice(t *testing.T) {
  data := []byte{}

  expected := uint32(0)
  actual := ReadUint32(data)

  if actual != expected {
    t.Error()
  }
}

func TestReadUint32SixByteSlice(t *testing.T) {
  data := []byte{ 0x28, 0x34, 0x71, 0xCA, 0xFB, 0x19 }

  expected := uint32(3396416552)
  actual := ReadUint32(data)

  if actual != expected {
    t.Error()
  }
}

func TestWriteZero(t *testing.T) {
  value := uint32(0)

  expected := []byte{ 0x00, 0x00, 0x00, 0x00 }
  actual := WriteUint32(value)

  if len(actual) != len(expected) {
    t.Error()
  }

  for i := 0; i < len(expected); i++ {
    if actual[i] != expected[i] {
      t.Error()
    }
  }
}

func TestWriteValue(t *testing.T) {
  value := uint32(3396416552)

  expected := []byte{ 0x28, 0x34, 0x71, 0xCA }
  actual := WriteUint32(value)

  if len(actual) != len(expected) {
    t.Error()
  }

  for i := 0; i < len(expected); i++ {
    if actual[i] != expected[i] {
      t.Error()
    }
  }
}
