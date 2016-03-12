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
