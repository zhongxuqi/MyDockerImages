package util

import (
  "os"
  "fmt"
  "bytes"
  "io/ioutil"
  "unicode/utf8"
  "net/http"
  "golang.org/x/text/encoding/simplifiedchinese"
  "golang.org/x/text/transform"
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

func WriteResp2File(filename string, resp *http.Response) {
  buf := bytes.NewBuffer([]byte(""))
  resp.Write(buf)
  Write2File(filename, buf.Bytes())
}

func GbkToUtf8(s []byte) ([]byte, error) {
  reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
  d, e := ioutil.ReadAll(reader)
  if e != nil {
      return nil, e
  }
  return d, nil
}

func DecodeResponse2Utf8Bytes(resp *http.Response) ([]byte, error) {
  buf := bytes.NewBuffer([]byte(""))
	resp.Write(buf)
  reader := transform.NewReader(bytes.NewReader(buf.Bytes()),
		simplifiedchinese.GBK.NewDecoder())
  b, err := ioutil.ReadAll(reader)
  if err != nil {
			return []byte{}, err
  }
	return b, nil
}

func ConvBytes2Reader(data []byte) (*bytes.Reader) {
  return  bytes.NewReader(data)
}
