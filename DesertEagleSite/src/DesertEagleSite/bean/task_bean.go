package bean

import (
  "errors"
  "strconv"
)

type UrlTask struct {
  Url string
  Keywords []string
  Info DataItem
}

type UrlResult struct {
  Task UrlTask
  Eval []int
}

func (result *UrlResult) Compare(another *UrlResult) (int, error) {
  if len(result.Eval) != len(another.Eval) {
    return 0, errors.New("length is not match")
  }
  for i := 0; i < len(result.Eval); i++ {
    if result.Eval[i] > another.Eval[i] {
      return 1, nil
    } else if result.Eval[i] < another.Eval[i] {
      return -1, nil
    }
  }
  return 0, nil
}

func (result *UrlResult) ToString() (str string) {
  str = result.Task.Info.Title + ", " +result.Task.Url + ":"
  str += "["
  for _, item := range result.Eval {
    str += strconv.Itoa(item) + ","
  }
  str += "]"
  return
}
