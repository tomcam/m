package util
import (
  "io/ioutil"
)

// fileToBytes converts a source file to a byte slice.
// In the spirit of web servers if the file can't be
// opened it just returns an empty slice. No error
// is generated.
func FileToBytes(filename string) []byte {
  bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
  return bytes
}


