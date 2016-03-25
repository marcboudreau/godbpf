package util

// ReadUint32 reads up to 4 bytes in little endian order and constructs an
// unsigned 32-bit integer value.
func ReadUint32(bytes []byte) (result uint32) {
  for i := 0; i < 4 && i < len(bytes); i++ {
    result |= uint32(bytes[i]) << uint32(8 * i)
  }

  return result
}

// WriteUint32 writes am unsigned 32-bit integer value into a byte slice.
func WriteUint32(value uint32) (result []byte) {
  result = make([]byte, 4)

  for i := 0; i < 4; i++ {
    result[i] = byte(value >> uint32(8 * i) & 0xFF)
  }

  return result
}
