package entry

import (
  "testing"
)

func TestDBPFEntryTGIString(t *testing.T) {
  tgi := &DBPFEntryTGI{TypeId: 0x12345678, GroupId: 0x9abcdef0, InstanceId: 0x13579bdf}

  expected := "T: 0x12345678, G: 0x9ABCDEF0, I: 0x13579BDF"
  actual := tgi.String()

  if expected != actual {
    t.Error()
  }
}

func TestZeroDBPFEntryTGIString(t *testing.T) {
  tgi := new(DBPFEntryTGI)

  expected := "T: 0x00000000, G: 0x00000000, I: 0x00000000"
  actual := tgi.String()

  if expected != actual {
    t.Error()
  }
}

func TestDBPFEntryTGIEquals(t *testing.T) {
  tgi := &DBPFEntryTGI{TypeId: 0x12345678, GroupId: 0x9abcdef0, InstanceId: 0x13579bdf}
  same := tgi
  equivalent := &DBPFEntryTGI{TypeId: 0x12345678, GroupId: 0x9abcdef0, InstanceId: 0x13579bdf}
  different := &DBPFEntryTGI{TypeId: 0x01234567, GroupId: 0x89abcdef, InstanceId: 0x02468ace}

  if !tgi.Equals(tgi) {
    t.Error()
  }

  if !tgi.Equals(same) {
    t.Error()
  }

  if !tgi.Equals(equivalent) {
    t.Error()
  }

  if tgi.Equals(different) {
    t.Error()
  }

  if tgi.Equals(nil) {
    t.Error()
  }
}

func TestDBPFEntryBytes(t *testing.T) {
  bytes := make([]byte, 12)
  tgi := &DBPFEntryTGI{TypeId: 0x32F9A014, GroupId: 0x7B5CA68E, InstanceId: 0x2468ACE0}
  expected := []byte{ 0x14, 0xA0, 0xF9, 0x32, 0x8E, 0xA6, 0x5C, 0x7B, 0xE0, 0xAC, 0x68, 0x24}

  tgi.Bytes(bytes)

  CheckIfSlicesAreEqual(t, bytes, expected)
}

func CheckIfSlicesAreEqual(t *testing.T, actual, expected []byte) {
  if len(actual) != len(expected) {
    t.Errorf("Actual slice size %d didn't match expected size %d", len(actual), len(expected))
  }

  for i, v := range actual {
    if v != expected[i] {
      t.Errorf("Byte %d: expected %2x, but actually was %2x", i, expected[i], v)
    }
  }
}
