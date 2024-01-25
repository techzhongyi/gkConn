package helper

// BccChecksum bcc校验
func BccChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum ^= b
	}
	return checksum
}
