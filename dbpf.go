package godbpf

import (
    "io"
    "time"
    "encoding/binary"
    "errors"
    "bytes"
    "container/list"
    "fmt"

    "github.com/marcboudreau/godbpf/entry"
    "github.com/marcboudreau/godbpf/qfs"
    "github.com/marcboudreau/godbpf/util"
)

// DBPF is a structure that encompasses all of the contents of a DBPF file.
type DBPF struct {
  // The major version of the DBPF schema used.
  MajorVersion uint32

  // The minor version of the DBPF schema used.
  MinorVersion uint32

  // The major version of the index schema used in this DBPF instance.
  IndexMajorVersion uint32

  // The minor version of the index schema used in this DBPF instance.
  IndexMinorVersion uint32

  // Timestamp of when this DBPF was created.
  CreatedDate time.Time

  // Timestamp of when this DBPF was last modified.
  ModifiedDate time.Time

  // entries points to a list of DBPFEntry (exluding the DBPFDirEntry) if one
  // was created), contained in the DBPF instance.
  entries *list.List
}

// New creates a new DBPF instance.
func New() *DBPF {
  dbpf := new(DBPF)
  dbpf.entries = list.New()

  return dbpf
}

// Len returns the number of entries in this DBPF instance.
func (dbpf *DBPF) Len() uint32 {
  return uint32(dbpf.entries.Len())
}

// Parse creates a DBPF from the data read from the provided reader.
func Parse(r io.Reader) (*DBPF, error) {
  dbpf := New()
  count, offset, _ := parseHeader(r, dbpf)

  contentBuf := make([]byte, offset)
  if offset > 96 {
    // Start reading in bytes at position 96 to account for the header data, that
    // way all of the offset values read in the file will align in this slice as
    // well.
    if _, e := r.Read(contentBuf[96:]); e != nil {
      return nil, e
    }
  }

  dbpf.parseEntries(r, offset, count, contentBuf)

  return dbpf, nil
}

// parseHeader parses the header section of the DBPF file and populates the values
// in the provided DBPF instance.
func parseHeader(r io.Reader, dbpf *DBPF) (count, offset uint32, e error) {
  fmt.Println("ENTER parseHeader")
  defer fmt.Println("EXIT parseHeader")

  magic := make([]byte, 4)
  if _, e := r.Read(magic); e != nil {
    return 0, 0, e
  }

  if string(magic) != "DBPF" {
    return 0, 0, errors.New("Invalid magic number")
  }

  binary.Read(r, binary.LittleEndian, &dbpf.MajorVersion)
  binary.Read(r, binary.LittleEndian, &dbpf.MinorVersion)

  fmt.Printf("  dbpf.MajorVersion = %d\n", dbpf.MajorVersion)
  fmt.Printf("  dbpf.MinorVersion = %d\n", dbpf.MinorVersion)

  r.Read(make([]byte, 12))

  var timestamp uint32
  binary.Read(r, binary.LittleEndian, &timestamp)
  dbpf.CreatedDate = time.Unix(int64(timestamp), 0)
  binary.Read(r, binary.LittleEndian, &timestamp)
  dbpf.ModifiedDate = time.Unix(int64(timestamp), 0)

  fmt.Printf("  dbpf.CreatedDate = 0x%08X\n", dbpf.CreatedDate.Unix())
  fmt.Printf("  dbpf.ModifiedDate = 0x%08X\n", dbpf.ModifiedDate.Unix())

  binary.Read(r, binary.LittleEndian, &dbpf.IndexMajorVersion)
  binary.Read(r, binary.LittleEndian, &count)
  binary.Read(r, binary.LittleEndian, &offset)

  fmt.Printf("  index.MajorVersion = %d\n", dbpf.IndexMajorVersion)
  fmt.Printf("  count = %d\n", count)
  fmt.Printf("  offset = %d\n", offset)

  // Gobble up the index size and 3 other  unused uint32 values
  r.Read(make([]byte, 16))

  binary.Read(r, binary.LittleEndian, &dbpf.IndexMinorVersion)

  fmt.Printf("  index.MinorVersion = %d\n", dbpf.IndexMinorVersion)

  // Gobble up 10 unused uint32 values
  r.Read(make([]byte, 40))

  return count, offset, nil
}

func (dbpf *DBPF) parseEntries(r io.Reader, indexOffset, indexCount uint32, content []byte) {
  for i := 0; i < int(indexCount); i++ {
    var typeId, groupId, instanceId, location, size uint32

    binary.Read(r, binary.LittleEndian, &typeId)
    binary.Read(r, binary.LittleEndian, &groupId)
    binary.Read(r, binary.LittleEndian, &instanceId)
    binary.Read(r, binary.LittleEndian, &location)
    binary.Read(r, binary.LittleEndian, &size)

    entry := &entry.DBPFEntry{TGI: &entry.DBPFEntryTGI{TypeId: typeId, GroupId: groupId, InstanceId: instanceId}}
    entry.SetData(content[location:location + size])

    dbpf.entries.PushBack(entry)
  }
}

func (dbpf *DBPF) Save(w io.Writer) error {
  w.Write([]byte("DBPF"))

  binary.Write(w, binary.LittleEndian, dbpf.MajorVersion)
  binary.Write(w, binary.LittleEndian, dbpf.MinorVersion)

  w.Write(make([]byte, 12))

  binary.Write(w, binary.LittleEndian, uint32(dbpf.CreatedDate.Unix()))
  binary.Write(w, binary.LittleEndian, uint32(dbpf.ModifiedDate.Unix()))

  binary.Write(w, binary.LittleEndian, dbpf.IndexMajorVersion)
  binary.Write(w, binary.LittleEndian, dbpf.Len())

  contentBuf := dbpf.encodeContent()

  binary.Write(w, binary.LittleEndian, uint32(contentBuf.Len() + 96))
  binary.Write(w, binary.LittleEndian, uint32(20 * dbpf.Len()))

  w.Write(make([]byte, 48))

  // Write the file entries
  contentBuf.WriteTo(w)

  // // Write the index table
  dbpf.encodeIndex(w)

  return nil
}

func (dbpf *DBPF) encodeContent() *bytes.Buffer {
  buf := new(bytes.Buffer)

  for elem := dbpf.entries.Front(); elem != nil; elem = elem.Next() {
    if entry, ok := elem.Value.(*entry.DBPFEntry); ok {
      buf.Write(entry.GetData())
    }
  }

  return buf
}

// encodeIndex encodes the index entries and writes the binary data to the
// provided io.Writer.
func (dbpf *DBPF) encodeIndex(w io.Writer) {
  location := uint32(96)
  for elem := dbpf.entries.Front(); elem != nil; elem = elem.Next() {
    if entry, ok := elem.Value.(*entry.DBPFEntry); ok {
      binary.Write(w, binary.LittleEndian, entry.TGI.TypeId)
      binary.Write(w, binary.LittleEndian, entry.TGI.GroupId)
      binary.Write(w, binary.LittleEndian, entry.TGI.InstanceId)
      binary.Write(w, binary.LittleEndian, location)
      binary.Write(w, binary.LittleEndian, entry.Size())

      location += entry.Size()
    }
  }
}

// AddEntry adds the provided DBPFEntry instance to the DBPF instance and creates
// an entry in the DBPFIndex instance as well.
func (dbpf *DBPF) AddEntry(e *entry.DBPFEntry) {
  dbpf.entries.PushBack(e)
}

func (dbpf *DBPF) AddCompressedEntry(tgi *entry.DBPFEntryTGI, uncompressedData []byte) {
  entry := &entry.DBPFEntry{TGI: tgi}
  buf := new(bytes.Buffer)
  qfs.Encode(buf, uncompressedData)

  data := make([]byte, 4 + buf.Len())
  copy(data[0:4], util.WriteUint32(uint32(buf.Len())))
  copy(data[4:], buf.Bytes())
  entry.SetData(data)

  dbpf.entries.PushBack(entry)

  dirEntry := dbpf.GetDirEntry()
  dirEntry.AddEntry(tgi, uint32(len(uncompressedData)))
}

// GetDirEntry locates and returns the DIR entry in the receiver.  If there isn't
// one, then one is created.
func (dbpf *DBPF) GetDirEntry() *entry.DBPFEntry {
  dirEntry := dbpf.Find(entry.DIR_ENTRY_TGI)
  if dirEntry == nil {
    dirEntry = entry.CreateDirEntry()
    dbpf.AddEntry(dirEntry)
  }

  return dirEntry
}

// Find searches the index and returns the related DBPFEntry instance.
func (dbpf *DBPF) Find(tgi *entry.DBPFEntryTGI) *entry.DBPFEntry {
  fmt.Printf("ENTER Find(tgi: %s)\n", tgi.String())
  defer fmt.Println("EXIT Find")

  for elem := dbpf.entries.Front(); elem != nil; elem = elem.Next() {
    if entry, ok := elem.Value.(*entry.DBPFEntry); ok {
      fmt.Printf("  entry: %s\n", entry.String())
      if entry.TGI.Equals(tgi) {
        fmt.Println("  found the match")
        return entry
      }
    }
  }

  return nil
}
