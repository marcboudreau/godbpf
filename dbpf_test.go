package godbpf

import (
  "testing"
  "time"
  "bytes"

  "github.com/marcboudreau/godbpf/entry"
)

func TestSchemaVersionOnNewDBPF(t *testing.T) {
  dbpf := &DBPF{MajorVersion: 1, MinorVersion: 0}

  // Make sure the schema version is set to 7.0
  if dbpf.MajorVersion != 1 {
    t.Error()
  }

  if dbpf.MinorVersion != 0 {
    t.Error()
  }
}

func TestCreatedAndModifiedDateOnNewDBPF(t *testing.T) {
  dbpf := &DBPF{CreatedDate: time.Now(), ModifiedDate: time.Now()}
  now := time.Now().Unix()
  created := dbpf.CreatedDate.Unix()
  modified := dbpf.ModifiedDate.Unix()

  // The created date should no more than 1 second earlier than now
  if now - created > 1 {
    t.Error()
  }

  // The modified date should no more than 1 second earlier than now
  if now - modified > 1 {
    t.Error()
  }
}

func TestLenOnNewDBPF(t *testing.T) {
  dbpf := New()

  // Make sure the index is empty
  if dbpf.Len() != 0 {
    t.Error()
  }
}

func TestParseDBPF(t *testing.T) {
  data := []byte{'D', 'B', 'P', 'F',
                 1, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 224, 125, 223, 86,
                 224, 125, 223, 86,
                 7, 0, 0, 0,
                 0, 0, 0, 0,
                 96, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0 }
  buf := bytes.NewBuffer(data)

  dbpf, err := Parse(buf)
  if err != nil {
    t.Error()
  }

  if dbpf.MajorVersion != 1 {
    t.Error()
  }

  if dbpf.MinorVersion != 0 {
    t.Error()
  }
}

func TestParseDBPFWithDifferentModifiedDate(t *testing.T) {
  data := []byte{'D', 'B', 'P', 'F',
                 1, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 224, 125, 223, 86,
                 102, 139, 223, 86,
                 7, 0, 0, 0,
                 0, 0, 0, 0,
                 96, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0,
                 0, 0, 0, 0 }

  buf := bytes.NewBuffer(data)

  dbpf, err := Parse(buf)
  if err != nil {
    t.Error()
  }

  if dbpf.MajorVersion != 1 {
    t.Error()
  }

  if dbpf.MinorVersion != 0 {
    t.Error()
  }
}

// func TestCreateNewDBPFWithAnEntry(t *testing.T) {
//   dbpf := New()
//   data := []byte{ 0xA5, 0x5A, 0xA5, 0x5A,
//                   0xA5, 0x5A, 0xA5, 0x5A,
//                   0xA5, 0x5A, 0xA5, 0x5A,
//                   0xA5, 0x5A, 0xA5, 0x5A,
//                   0xA5, 0x5A, 0xA5, 0x5A }
//
//   tgi := &entry.DBPFEntryTGI{ TypeId: 0x7ab50e44, GroupId: 0x0986135e, InstanceId: 0xffff4000 }
//   entry := &entry.DBPFEntry{ TGI: tgi }
//   entry.SetData(data)
//
//   dbpf.Entries.PushBack(entry)
//
//   if dbpf.Index.Len() != 1 {
//     t.Error()
//   }
// }

func TestFindEntryUsingTGI(t *testing.T) {
  dbpf := New()
  data1 := []byte{ 0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA }

  data2 := []byte{ 0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55 }

  tgi1 := &entry.DBPFEntryTGI{ TypeId: 0x7ab50e44, GroupId: 0x0986135e, InstanceId: 0xffff4000 }
  tgi2 := &entry.DBPFEntryTGI{ TypeId: 0x7ab50e44, GroupId: 0x0986135e, InstanceId: 0xffff4005 }

  entry1 := entry.NewEntry(tgi1)
  entry2 := entry.NewEntry(tgi2)
  entry1.SetData(data1)
  entry2.SetData(data2)

  dbpf.AddEntry(entry1)
  dbpf.AddEntry(entry2)

  if entry := dbpf.Find(tgi1); entry.GetData()[0] != byte(0xAA) {
    t.Error()
  }

  if entry := dbpf.Find(tgi2); entry.GetData()[0] != byte(0x55) {
    t.Error()
  }
}

func TestFindIndexEntries(t *testing.T) {
  dbpf := New()

  data1 := []byte{ 0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA,
                  0xAA, 0xAA, 0xAA, 0xAA }

  data2 := []byte{ 0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55,
                  0x55, 0x55, 0x55, 0x55 }

  tgi1 := &entry.DBPFEntryTGI{ TypeId: 0x7ab50e44, GroupId: 0x0986135e, InstanceId: 0xffff4000 }
  tgi2 := &entry.DBPFEntryTGI{ TypeId: 0x7ab50e44, GroupId: 0x0986135e, InstanceId: 0xffff4005 }

  entry1 := entry.NewEntry(tgi1)
  entry2 := entry.NewEntry(tgi2)

  entry1.SetData(data1)
  entry2.SetData(data2)

  dbpf.AddEntry(entry1)
  dbpf.AddEntry(entry2)

  if dbpf.Len() != 2 {
    t.Error()
  }

  if entry := dbpf.Find(tgi1); !entry.TGI.Equals(tgi1) {
    t.Error()
  }

  if entry := dbpf.Find(tgi2); !entry.TGI.Equals(tgi2) {
    t.Error()
  }
}

func TestSaveEmptyDBPF(t *testing.T) {
  dbpf := New()
  dbpf.MajorVersion = 1
  dbpf.MinorVersion = 0
  dbpf.CreatedDate = time.Unix(3465168386, 0)
  dbpf.ModifiedDate = time.Unix(3751499539, 0)
  dbpf.IndexMajorVersion = 7
  dbpf.IndexMinorVersion = 0
  byteBuffer := []byte{}
  buffer := bytes.NewBuffer(byteBuffer)
  expected := []byte{ 'D', 'B', 'P', 'F',
                      0x1, 0x0, 0x0, 0x0,
                      0x0, 0x0, 0x0, 0x0,
                      0x0, 0x0, 0x0, 0x0,
                      0x0, 0x0, 0x0, 0x0,
                      0x0, 0x0, 0x0, 0x0,
                      0x02, 0x46, 0x8A, 0xCE,
                      0x13, 0x57, 0x9B, 0xDF,
                      0x7, 0x0, 0x0, 0x0,
                      0x0, 0x0, 0x0, 0x0 }

  dbpf.Save(buffer)

  for i := 0; i < 40; i++ {
    if c, _ := buffer.ReadByte(); c != expected[i] {
      t.Errorf("Assertion failed. Byte %d expecting %2x but was %2x\n", i, expected[i], c)
    }
  }
}

func TestSaveSingleEntryDBPF(t *testing.T) {
  dbpf := New()
  dbpf.MajorVersion = 1
  dbpf.MinorVersion = 0
  dbpf.CreatedDate = time.Unix(3465168386, 0)
  dbpf.ModifiedDate = time.Unix(3751499539, 0)
  dbpf.IndexMajorVersion = 7
  dbpf.IndexMinorVersion = 0
  byteBuffer := []byte{}
  buffer := bytes.NewBuffer(byteBuffer)
  expected := []byte{ 'D', 'B', 'P', 'F',
                       0x1, 0x0, 0x0, 0x0, // Major version
                       0x0, 0x0, 0x0, 0x0, // Minor version
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x02, 0x46, 0x8A, 0xCE, // Created date
                       0x13, 0x57, 0x9B, 0xDF, // Modified date
                       0x7, 0x0, 0x0, 0x0, // Index major version
                       0x2, 0x0, 0x0, 0x0, // Index entry count
                       0x7E, 0x0, 0x0, 0x0, // Offset of first index entry
                       0x28, 0x0, 0x0, 0x0, // Size of the index
                       0x0, 0x0, 0x0, 0x0, // Hole entry count
                       0x0, 0x0, 0x0, 0x0, // Hole offset
                       0x0, 0x0, 0x0, 0x0, // Hole size
                       0x0, 0x0, 0x0, 0x0, // Index minor version
                       0x0, 0x0, 0x0, 0x0, // Index offset
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       0x0, 0x0, 0x0, 0x0,
                       // First file
                       0xA, 0x0, 0x0, 0x0, // Compressed size
                       0x10, 0xFB, // Compression ID
                       0x0, 0x0, 0x14, // Uncompressed size
                       0x8F, 0x40, 0x0, 0xAA, // Compressed data
                       0xFC,
                       // DIR file
                       0x44, 0xE, 0xB5, 0x7A, // TypeId
                       0x5E, 0x13, 0x86, 0x9, // GroupId
                       0x0, 0x40, 0xFF, 0xFF, // InstanceId
                       0x14, 0x0, 0x0, 0x0, // Uncompressed size
                       // Index
                       //  First entry
                       0x44, 0xE, 0xB5, 0x7A, // TypeId
                       0x5E, 0x13, 0x86, 0x9, // GroupId
                       0x0, 0x40, 0xFF, 0xFF, // InstanceId
                       0x60, 0x0, 0x0, 0x0, // Location in file
                       0xE, 0x0, 0x0, 0x0, // Size
                       //  Dir entry
                       0xEF, 0x1E, 0x6B, 0xE8, // TypeId
                       0xEF, 0x1E, 0x6B, 0xE8, // GroupId
                       0x3, 0x1F, 0x6B, 0x28, // InstanceId
                       0x6E, 0x0, 0x0, 0x0, // Location in file
                       0x10, 0x0, 0x0, 0x0 } // Size

  tgi := &entry.DBPFEntryTGI{ TypeId: 0x7ab50e44, GroupId: 0x0986135e, InstanceId: 0xffff4000 }
  dbpf.AddCompressedEntry(tgi, []byte{ 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA })

  dbpf.Save(buffer)

  for i := 0; i < len(expected); i++ {
    if c, _ := buffer.ReadByte(); c != expected[i] {
      t.Errorf("Assertion failed. Byte %d expecting %2x but was %2x\n", i, expected[i], c)
    }
  }
}
