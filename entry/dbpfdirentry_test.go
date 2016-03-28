package entry

import (
  "testing"
)

func TestCreateDirEntry(t *testing.T) {
  dirEntry := CreateDirEntry()

  if !dirEntry.TGI.Equals(DIR_ENTRY_TGI) {
    t.Error()
  }
}

func TestAddEntryWithInvalidEntry(t *testing.T) {
  invalidTgi := &DBPFEntryTGI{TypeId: 0x11112222, GroupId: 0x33334444, InstanceId: 0x55556666}
  entry := NewEntry(invalidTgi)
  someTgi := &DBPFEntryTGI{TypeId: 0xFFFF0000, GroupId: 0xEEEE0000, InstanceId: 0xDDDD0000}

  defer func() {
    if r := recover(); r == nil {
      t.Error()
    }
  }()

  entry.AddEntry(someTgi, 100)
}

func TestAddEntryWhenNoneExists(t *testing.T) {
  dirEntry := CreateDirEntry()
  someTgi := &DBPFEntryTGI{TypeId: 0xFFFF0000, GroupId: 0xEEEE0000, InstanceId: 0xDDDD0000}
  expected := []byte{0x0, 0x0, 0xFF, 0xFF, 0x0, 0x0, 0xEE, 0xEE, 0x0, 0x0, 0xDD, 0xDD, 0x64, 0x0, 0x0, 0x0 }

  dirEntry.AddEntry(someTgi, 100)
  actual := dirEntry.GetData()

  CheckIfSlicesAreEqual(t, actual, expected)
}

func TestAddMoreEntries(t *testing.T) {
  dirEntry := CreateDirEntry()
  someTgi := &DBPFEntryTGI{TypeId: 0xFFFF0000, GroupId: 0xEEEE0000, InstanceId: 0xDDDD0000}
  someTgi2 := &DBPFEntryTGI{TypeId: 0x12345678, GroupId: 0x87654321, InstanceId: 0xFACDDBBE}
  expected := []byte{0x0, 0x0, 0xFF, 0xFF, 0x0, 0x0, 0xEE, 0xEE, 0x0, 0x0, 0xDD, 0xDD, 0x64, 0x0, 0x0, 0x0, 0x78, 0x56, 0x34, 0x12, 0x21, 0x43, 0x65, 0x87, 0xBE, 0xDB, 0xCD, 0xFA, 0x63, 0x0, 0x0, 0x0 }

  dirEntry.AddEntry(someTgi, 100)
  dirEntry.AddEntry(someTgi2, 99)
  actual := dirEntry.GetData()

  CheckIfSlicesAreEqual(t, actual, expected)
}
