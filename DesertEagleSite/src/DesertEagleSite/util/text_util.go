package util

import (
  "os"
  "fmt"
  "unicode/utf8"
)

func DecodeUtf8String(encodedString string) (decodedString string) {
  decodedString = ""
  if utf8.ValidString(encodedString) {
    for len(encodedString) > 0 {
      r, size := utf8.DecodeRuneInString(encodedString)
      decodedString = decodedString + string(r)
      encodedString = encodedString[size:]
    }
  } else {
    decodedString = encodedString
  }
  return
}

func Write2File(filename string, data []byte) {
  f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  defer f.Close()
  f.Write(data)
}
