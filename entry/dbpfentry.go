package entry

import (
  "bytes"
  "fmt"

  "github.com/marcboudreau/godbpf/qfs"
)


// The DBPFEntry struct defines three methods common to all entries within a
// DBPF file.
type DBPFEntry struct {
  // The TGI field is a pointer to a DBPFEntryTGI struct, which encapsulates a
  // TypeId, a GroupId, and an InstanceId value.
  TGI *DBPFEntryTGI

  // The Data method returns the bytes that the make up the entry.
  data []byte
}

// NewEntry creates a new empty entry with the provided DBPFEntryTGI instance.
func NewEntry(tgi *DBPFEntryTGI) *DBPFEntry {
  return &DBPFEntry{ TGI: tgi, data: nil }
}

// Size returns the size of the data stored in this entry.
func (e *DBPFEntry) Size() uint32 {
  return uint32(len(e.data))
}

// SetData updates the data stored in this entry, without performing any kind
// of encoding.
func (e *DBPFEntry) SetData(data []byte) {
  e.data = make([]byte, len(data))
  copy(e.data, data)
}

// GetData retrieves the data stored in this entry as is, without doing any
// decoding.
func (e *DBPFEntry) GetData() []byte {
  return e.data
}

// ReadData decodes the data stored in this entry, if that's necessary (e.g. QFS
// compressed data).
func (e *DBPFEntry) ReadData() ([]byte, error) {
  buffer := bytes.NewBuffer(e.data)

  return qfs.Decode(buffer)
}

func (e *DBPFEntry) WriteData(data []byte) {
  buffer := new(bytes.Buffer)

  qfs.Encode(buffer, data)
  e.data = buffer.Bytes()
}

// String returns a string representation of the receiver.
func (e *DBPFEntry) String() string {
  return fmt.Sprintf("TGI: %v, data: %v\n", e.TGI, e.data)
}
