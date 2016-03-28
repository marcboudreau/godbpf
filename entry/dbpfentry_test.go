package entry

import (
  "testing"
)

func TestNewEntry(t *testing.T) {
  tgi := &DBPFEntryTGI{TypeId: 0x33333333, GroupId: 0x66666666, InstanceId: 0x99999999}

  entry := NewEntry(tgi)

  if !tgi.Equals(entry.TGI) {
    t.Error()
  }

  if 0 != entry.Size() {
    t.Error()
  }
}

func TestDBPFEntryData(t *testing.T) {
  tgi := &DBPFEntryTGI{TypeId: 0x33333333, GroupId: 0x66666666, InstanceId: 0x99999999}
  entry := NewEntry(tgi)
  expected := []byte{0x11}
  entry.SetData(make([]byte, 0))

  if 0 != entry.Size() {
    t.Error()
  }

  entry.SetData(expected)

  if 1 != entry.Size() {
    t.Error()
  }

  actual := entry.GetData()

  if len(expected) != len(actual) {
    t.Error()
  }

  if expected[0] != actual[0] {
    t.Error()
  }
}

func TestDBPFEntryString(t *testing.T) {
  tgi := &DBPFEntryTGI{TypeId: 0x33333333, GroupId: 0x66666666, InstanceId: 0x99999999}
  entry := NewEntry(tgi)

  entry.SetData([]byte{0x22})

  actual := entry.String()
  expected := "TGI: T: 0x33333333, G: 0x66666666, I: 0x99999999, data: [34]\n"

  if actual != expected {
    t.Error("Actual [" + actual + "] didn't match expected [" + expected + "]")
  }
}
