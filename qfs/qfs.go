package qfs

import (
  "io"
  "encoding/binary"
)

// Encode takes the bytes contained in the provided byte slice, encodes them, and
// writes them to the Writer.
func Encode(w io.Writer, data []byte) error {
  offsetToLastOccurence := make([]int, 256)
  for i, _ := range offsetToLastOccurence { offsetToLastOccurence[i] = -1 }

  nextWritePos := 0 // The index of the next byte from data to be written to w

  if e := binary.Write(w, binary.LittleEndian, uint16(0xFB10)); e != nil {
    return e
  }

  size := []byte { byte(len(data) >> 16 & 0xFF), byte(len(data) >> 8 & 0xFF), byte(len(data) & 0xFF) }
  if _, e := w.Write(size); e != nil {
    return e
  }

  for i := 0; i < len(data); i++ {
    v := data[i]
    if -1 != offsetToLastOccurence[v] {
      repeatCount := firstByteRepeatCount(data[i:])
      if compressed, e := writeCompressible(w, data[nextWritePos:i], repeatCount, uint32(i - offsetToLastOccurence[v])); e != nil {
        return e
      } else if compressed {
        nextWritePos = i + int(repeatCount)
        i += int(repeatCount) - 1
      } else {
        offsetToLastOccurence[v] = i
      }
    } else {
      offsetToLastOccurence[v] = i
    }
  }

  // Need to write the remaining non-repeating bytes
  if e := writeFinalBlocks(w, data[nextWritePos:]); e != nil {
    return e
  }

  return nil
}

// firstByteRepeatCount counts how many times the first byte of the provided
// slice is repeated consecutively.
func firstByteRepeatCount(s []byte) uint32 {
  i := 1
  for ; i < len(s) && s[0] == s[i]; i++ {}

  return uint32(i)
}

// writeCompressible determines if the provided inputs data can be compressed. If
// it can, it encodes the data and writes it to the Writer; otherwise it does
// nothing.  writeCompressible returns the number of bytes written.
func writeCompressible(w io.Writer, data []byte, copyCount, copyOffset uint32) (bool, error) {
  var control []byte = nil
  compressed := false

  if copyCount >= 3 && copyCount <= 10 && copyOffset <= 1024 {
    // 3 to 10 repeating bytes and copyOffset no greater than 1024: can be compressed
    // with 2 control bytes
    control = createTwoControlByteBlock(uint32(len(data)), copyCount, copyOffset)
  } else if copyCount >= 4 && copyCount <= 67 && copyOffset <= 16384 {
    // 4 to 67 repeating bytes and copyOffset no greater than 16384: can be compressed
    // with 3 control bytes
    control = createThreeControlByteBlock(uint32(len(data)), copyCount, copyOffset)
  } else if copyCount >= 5 && copyCount <= 1028 && copyOffset <= 131072 {
    // 5 to 1028 repeating bytes and copyOffset no greater than 131072: can be compressed
    // with 4 control bytes
    control = createFourControlByteBlock(uint32(len(data)), copyCount, copyOffset)
  }

  if control != nil {
    // If the data can be compressed, encode it and write it to the Writer
    consumedBytes, e := writeNonRepeatingBlocks(data, w)
    if e != nil {
      return false, e
    }

    if e := writeCompressedBlock(control, data[consumedBytes:], w); e != nil {
      return false, e
    }

    compressed = true
  }

  return compressed, nil
}

// writeFinalBlocks encodes and writes all remaining bytes including the final
// terminating block.
func writeFinalBlocks(w io.Writer, data []byte) error {
  if consumedBytes, e := writeNonRepeatingBlocks(data, w); e != nil {
    return e
  } else {
    data = data[consumedBytes:]
  }

  if e := writeCompressedBlock(createFinalControlByteBlock(uint32(len(data)), 0, 0), data, w); e != nil {
    return e
  }

  return nil
}

// writeNonRepeatingBlocks encodes and writes blocks of bytes that do not include
// any copied bytes from earlier in the stream.
func writeNonRepeatingBlocks(data []byte, w io.Writer) (consumedBytes int, e error) {
  // Blocks of non-repeating bytes are limited to 4 to 112 bytes in length (in
  // increments of 4).  Write as many full blocks of 112 as possible first, then
  // write the remaining bytes.
  for len(data) >= 112 {
    if e := writeCompressedBlock(createOneControlByteBlock(112, 0, 0), data[0:112], w); e != nil {
      return consumedBytes, e
    }

    data = data[112:]
    consumedBytes += 112
  }

  blockSize := int(len(data) / 4) * 4
  if blockSize > 0 {
    if e := writeCompressedBlock(createOneControlByteBlock(uint32(blockSize), 0, 0), data[0:blockSize], w); e != nil {
      return consumedBytes, e
    }

    consumedBytes += blockSize
  }

  return consumedBytes, nil
}

// createOneControlByteBlock creates a slice of 1 byte that incorporates the
// dataLen (p) into this bit pattern:
//  byte0: 111ppppp
// This function ignores the copyCount and copyOffset.
func createOneControlByteBlock(dataLen, copyCount, copyOffset uint32) []byte {
  control := make([]byte, 1)

  control[0] = byte(0xE0 | dataLen / 4 - 1)

  return control
}

// createTwoControlByteBlock creates a slice of 2 bytes that incorporates the
// dataLen (p), copyCount (c), and copyOffset(o) into this bit pattern:
//  byte0: 0oocccpp
//  byte1: oooooooo
func createTwoControlByteBlock(dataLen, copyCount, copyOffset uint32) []byte {
  copyOffset -= 1
  copyCount -= 3

  control := make([]byte, 2)

  control[0] = byte(copyOffset >> 8 & 0x3 << 5 | copyCount & 0x7 << 2 | dataLen & 0x3)
  control[1] = byte(copyOffset & 0xFF)

  return control
}

// createThreeControlByteBlock creates a slice of 3 bytes that incorporates the
// dataLen (p), copyCount (c), and copyOffset (o) into this bit pattern:
//  byte0: 10cccccc
//  byte1: ppoooooo
//  byte2: oooooooo
func createThreeControlByteBlock(dataLen, copyCount, copyOffset uint32) []byte {
  copyOffset -= 1
  copyCount -= 4

  control := make([]byte, 3)

  control[0] = byte(0x80 | copyCount & 0x3F)
  control[1] = byte(dataLen & 0x3 << 6 | copyOffset >> 8 & 0x3F)
  control[2] = byte(copyOffset & 0xFF)

  return control
}

// createFourControlByteBlock creates a slice of 4 bytes that incorporates the
// dataLen (p), copyCount (c), and copyOffset (o) into this bit pattern:
//  byte0: 110occpp
//  byte1: oooooooo
//  byte2: oooooooo
//  byte3: cccccccc
func createFourControlByteBlock(dataLen, copyCount, copyOffset uint32) []byte {
  copyOffset -= 1
  copyCount -= 5

  control := make([]byte, 4)

  control[0] = byte(0xC0 | copyOffset >> 16 & 0x1 << 4 | copyCount >> 8 & 0x3 << 2 | dataLen & 0x3)
  control[1] = byte(copyOffset >> 8 & 0xFF)
  control[2] = byte(copyOffset & 0xFF)
  control[3] = byte(copyCount & 0xFF)

  return control
}

// createFinalControlByteBlock creates a slice of 1 byte that incorporates the
// dataLen (p) into this bit pattern:
//  byte0: 111111pp
// This function ignores the copyCount and copyOffset.  This block also serves
// as the terminating block for the compressed data stream.
func createFinalControlByteBlock(dataLen, copyCount, copyOffset uint32) []byte {
  control := make([]byte, 1)

  control[0] = byte(0xFC | dataLen & 0x3)

  return control
}

// writeCompressedBlock writes the provided control bytes and the proceeding
// bytes passed in as the control and data slices to the io.Writer passed in.
// The number of bytes written, including the control bytes is returned.
func writeCompressedBlock(control, data []byte, w io.Writer) error {
  buffer := make([]byte, len(control) + len(data))

  copy(buffer, control)
  copy(buffer[len(control):], data)

  if _, e := w.Write(buffer); e != nil {
    return e
  }

  return nil
}

// Decode reads bytes from the provided Reader and decodes them.  The decoded
// bytes are returned in a byte slice.
func Decode(r io.Reader) ([]byte, error) {
  // The buffer is 113 bytes long because the largest byte sequence from a single
  // control code is a 1-byte control byte with 112 proceeding bytes.
  buffer := make([]byte, 113)

  // Consume the compression header
  if _, e := r.Read(buffer[0:2]); e != nil {
    return nil, e
  }

  // Read the uncompressed data size
  if _, e := r.Read(buffer[0:3]); e != nil {
    return nil, e
  }
  size := uint32(buffer[0] << 16 | buffer[1] << 8 | buffer[2])

  // Create the output byte slice
  output := make([]byte, size)
  outputPos := 0

  for {
    if n, e := r.Read(buffer[0:1]); n == 0 && e != nil {
      break;
    }

    if e := decodeSequence(buffer[0], r, output, &outputPos); e != nil {
      return nil, e
    }
  }

  return output, nil
}

// decodeSequence decodes a sequence of bytes and writes the decoded bytes in
// the provided byte slice.
func decodeSequence(control byte, r io.Reader, output []byte, outputPos *int) error {
  buffer := make([]byte, 4)
  buffer[0] = control

  var f func([]byte) (int, int, int)

  switch {
  case control & 0x80 == 0x0:
    // 2-byte control code
    if _, e := r.Read(buffer[1:2]); e != nil {
      return e
    }
    f = decodeTwoByteSequence
    break
  case control & 0xC0 == 0x80:
    // 3-byte control code
    if _, e := r.Read(buffer[1:3]); e != nil {
      return e
    }
    f = decodeThreeByteSequence
    break
  case control & 0xE0 == 0xC0:
    // 4-byte control code
    if _, e := r.Read(buffer[1:4]); e != nil {
      return e
    }
    f = decodeFourByteSequence
    break
  case control & 0xE4 == 0xE0:
    // 1-byte control code
    f = decodeOneByteSequence
    break
  case control & 0xFC == 0xFC:
    // Special 1-byte control code (terminator)
    f = decodeFinalSequence
    break
  }

  proceeding, count, offset := f(buffer)
  pos := *outputPos
  if n, e := r.Read(output[pos:pos + proceeding]); e != nil {
    return e
  } else {
    *outputPos += n
  }
  if 0 < offset {
    b := output[*outputPos - offset]
    for i := 0; i < count; i++ {
      output[*outputPos + i] = b
    }

    *outputPos += count
  }

  return nil
}

// decodeFourByteSequence decodes a 4-byte control sequence to extract the number
// of proceeding bytes, the number of times to write the repeating byte, and the
// offset of the repeating byte.
func decodeFourByteSequence(control []byte) (proceeding, count, offset int) {
  proceeding = int(control[0] & 0x3)
  count = int(control[0] & 0xC << 6 + control[3] + 5)
  offset = int(control[0] & 0x10 << 12 + control[1] << 8 + control[2] + 1)

  return
}

// decodeThreeByteSequence decodes a 3-byte control sequence to extract the number
// of procceeding bytes, the number of times to write the repeating byte, and the
// offset of the repeating byte.
func decodeThreeByteSequence(control []byte) (proceeding, count, offset int) {
  proceeding = int(control[1] & 0xC0 >> 6 & 0x3)
  count = int(control[0] & 0x3F + 4)
  offset = int(control[1] & 0x3F << 8 + control[2] + 1)

  return
}

// decodeTwoByteSequence decodes a 2-byte control sequence to extract the number
// of procceeding bytes, the number of times to write the repeating byte, and the
// offset of the repeating byte.
func decodeTwoByteSequence(control []byte) (proceeding, count, offset int) {
  proceeding = int(control[0] & 0x3)
  count = int(control[0] & 0x1C >> 2 + 3)
  offset = int(control[0] & 0x60 << 3 + control[1] + 1)

  return
}

// decodeOneByteSequence decodes a 1-byte control sequence to extract the number
// of procceeding bytes.  The repeating count and offset is simply set to zero.
func decodeOneByteSequence(control []byte) (proceeding, count, offset int) {
  proceeding = int(control[0] & 0x1F << 2 + 4)
  count = 0
  offset = 0

  return
}

// decodeFinalSequence decodes a 1-byte control sequence to extract the number
// of proceeding bytes.  The repeating count and offset is simply set to zero.
func decodeFinalSequence(control []byte) (proceeding, count, offset int) {
  proceeding = int(control[0] & 0x3)
  count = 0
  offset = 0

  return
}
