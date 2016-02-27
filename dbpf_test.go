package godbpf

import (
  "testing"
  "time"
  "bytes"
)

func TestSchemaVersionOnNewDBPF(t *testing.T) {
  dbpf := New()

  // Make sure the schema version is set to 7.0
  if dbpf.MajorVersion != 7 {
    t.Error()
  }

  if dbpf.MinorVersion != 0 {
    t.Error()
  }
}

func TestCreatedDateOnNewDBPF(t *testing.T) {
  dbpf := New()

  // Make sure the created time is within a second from now.
  now := time.Now().Unix()
  created := dbpf.CreatedDate.Unix()

  if now - created > 1 {
    t.Error()
  }
}

func TestModifiedDateOnNewDBPF(t *testing.T) {
  dbpf := New()

  // Make sure the modified time is within a second from now.
  now := time.Now().Unix()
  modified := dbpf.ModifiedDate.Unix()

  if now - modified > 1 {
    t.Error()
  }
}

func TestEntryListEmptyOnNewDBPF(t *testing.T) {
  dbpf := New()

  // Make sure the Entries list is not nil but is empty.
  if dbpf.Entries == nil {
    t.Error()
  }

  if dbpf.Entries.Len() != 0 {
    t.Error()
  }
}

func TestParseDBPF(t *testing.T) {
  data := []byte{'D', 'B', 'P', 'F',
                 7, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 224, 125, 223, 86,
                 224, 125, 223, 86 }
  buf := bytes.NewBuffer(data)

  dbpf, err := Parse(buf)
  if err != nil {
    t.Error()
  }

  if dbpf.MajorVersion != 7 {
    t.Error()
  }

  if dbpf.MinorVersion != 0 {
    t.Error()
  }
}

func TestParseDBPFWithDifferentModifiedDate(t *testing.T) {
  data := []byte{'D', 'B', 'P', 'F',
                 7, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 224, 125, 223, 86,
                 102, 139, 223, 86 }
  buf := bytes.NewBuffer(data)

  dbpf, err := Parse(buf)
  if err != nil {
    t.Error()
  }

  if dbpf.MajorVersion != 7 {
    t.Error()
  }

  if dbpf.MinorVersion != 0 {
    t.Error()
  }
}
