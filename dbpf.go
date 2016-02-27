package godbpf

import (
    "io"
    "bufio"
    "time"
    "container/list"
    "bytes"
    "encoding/binary"
    "errors"
)

// DBPF is a structure that encompasses all of the contents of a DBPF file.
type DBPF struct {
  // The major version of the DBPF schema used.
  MajorVersion uint32

  // The minor version of the DBPF schema used.
  MinorVersion uint32

  // Timestamp of when this DBPF was created.
  CreatedDate time.Time

  // Timestamp of when this DBPF was last modified.
  ModifiedDate time.Time

  // List of entries contained in this DBPF.
  Entries *list.List
}

// New creates a new empty DBPF.
func New() *DBPF {
  return &DBPF{MajorVersion: 7,
               MinorVersion: 0,
               CreatedDate: time.Now(),
               ModifiedDate: time.Now(),
               Entries: list.New()}
}

func littleEndianReadUint32(buf io.Reader, data *uint32) {
  binary.Read(buf, binary.LittleEndian, data)
}

// Parse creates a DBPF from the data read from the provided reader.
func Parse(reader io.Reader) (*DBPF, error) {
  bufReader := bufio.NewReader(reader)
  magic := make([]byte, 4)
  if _, err := bufReader.Read(magic); err != nil {
    return nil, err
  }

  if magic[0] != 'D' && magic[1] != 'B' && magic[2] != 'P' && magic[3] != 'F' {
    return nil, errors.New("Invalid magic number")
  }

  var majorVersion uint32
  var minorVersion uint32

  littleEndianReadUint32(bufReader, &majorVersion)
  littleEndianReadUint32(bufReader, &minorVersion)

  bufReader.Discard(12)

  var createdDate uint32
  var modifiedDate uint32

  littleEndianReadUint32(bufReader, &createdDate)
  littleEndianReadUint32(bufReader, &modifiedDate)

  dbpf := &DBPF{majorVersion, minorVersion, time.Unix(int64(createdDate), 0), time.Unix(int64(modifiedDate), 0), list.New()}

  return dbpf, nil
}

func littleEndianWrite(buf io.Writer, value interface{}) {
  binary.Write(buf, binary.LittleEndian, value)
}

func (dbpf *DBPF) Save(writer io.Writer) error {
  buf := bytes.NewBuffer(make([]byte, 28))
  buf.WriteString("DBPF")
  littleEndianWrite(buf, dbpf.MajorVersion)
  littleEndianWrite(buf, dbpf.MinorVersion)
  littleEndianWrite(buf, []uint32{0, 0, 0})
  littleEndianWrite(buf, uint32(dbpf.CreatedDate.Unix()))
  littleEndianWrite(buf, uint32(dbpf.ModifiedDate.Unix()))
  // littleEndianWrite(buf, uint32(7))
  // littleEndianWrite(buf, uint32(dbpf.Entries.Len()))
  //
  // contentBuf := dbpf.encodeContent()
  //
  // littleEndianWrite(buf, uint32(contentBuf.Len() + 96))
  // littleEndianWrite(buf, uint32(20 * (dbpf.Entries.Len() + 1)))
  //
  // var filler [48]byte
  // buf.Write(filler)
  //
  // // Write the file entries
  // buf.Write(contentBuf.Bytes())
  //
  // // Write the index table
  // buf.Write(encodeIndex().Bytes())

  writer.Write(buf.Bytes())

  return nil
}

// func (dbpf *DBPF) encodeContent() *bytes.Buffer {
//   buf := new(bytes.Buffer)
//
//   return buf
// }

// func (dbpf *DBPF) encodeIndex() *bytes.Buffer {
//   buf := new(bytes.Buffer)
//
//   for elem := dbpf.Entries.Front(); elem.Next() != nil; elem = elem.Next() {
//     dbpfEntry := elem.Value
//
//     littleEndianWrite(buf, dbpfEntry.TypeId)
//     littleEndianWrite(buf, dbpfEntry.GroupId)
//     littleEndianWrite(buf, dbpfEntry.InstanceId)
//     littleEndianWrite(buf, dbpfEntry.Location)
//     littleEndianWrite(buf, uint32(dbpfEntry.Data.Len()))
//   }
//
//   return buf
// }


// EncodeIndexHeader encodes the receiver *DBPFIndex object into its binary form
//  and writes the data into a slice of bytes that's returned.
// func (index *DBPFIndex) EncodeIndexHeader() []byte, error {
//   buf := bytes.Buffer
//   binary.Write(buf, binary.LittleEndian, index.MajorVersion)
//   binary.Write(buf, binary.LittleEndian, index.MinorVersion)
//   binary.Write(buf, binary.LittleEndian, uint32(index.Entries.Len()))
//
//
//
//   return data
// }

// func (dbpf *DBPF) Encode() []byte, error {
//   buf := bytes.Buffer
//   buf.WriteString("DBPF")
//   binary.Write(buf, binary.LittleEndian, dbpf.MajorVersion)
//   binary.Write(buf, binary.LittleEndian, dbpf.MinorVersion)
//   binary.Write(buf, binary.LittleEndian, []uint32{0, 0, 0})
//   binary.Write(buf, binary.LittleEndian, uint32(dbpf.createdDate.Unix()))
//   binary.Write(buf, binary.LittleEndian, uint32(dbpf.modifiedDate.Unix()))
//   binary.Write(buf, binary.LittleEndian, dbpf.index.MajorVersion)
//   binary.Write(buf, binary.LittleEndian, uint32(dbpf.index.Entries.Len()))
//
//
//
//   binary.Write(buf, binary.LittleEndian, dbpf.index.MinorVersion)
//
// }

// type DBPFEntry struct {
//   TypeId uint32
//   GroupId uint32
//   InstanceId uint32
//   Location uint32
//   Data *bytes.Buffer
// }
