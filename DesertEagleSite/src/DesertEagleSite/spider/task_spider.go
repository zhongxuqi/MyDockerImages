package spider

import (
  "fmt"
  "time"
  . "DesertEagleSite/bean"
  "DesertEagleSite/wordtool"
  "DesertEagleSite/evaluator"
)

func execTask(taskQueue <-chan *UrlTask, resultQueue chan<- *UrlResult) {
END_LABEL:
  for {
    select {
    case task := <- taskQueue:
      fmt.Println("search: ", task.Info.Title)
      eval, err := evaluator.EvaluateUrlByKeyWords(task.Url, task.Keywords)
      if err != nil {
        continue
      }
      result := &UrlResult {
        Task: *task,
        Eval: eval,
      }
      resultQueue <- result
    case <- time.After(time.Second):
      break END_LABEL
    }
  }
}

func submitTask(taskQueue chan<- *UrlTask, UrlList []DataItem, keywords []string) {
  tmpMap := make(map[string]DataItem)
  for _, item := range UrlList {
    if _, ok := tmpMap[item.Link]; ok {
      continue
    }
    tmpMap[item.Link] = item
    task := &UrlTask{
      Info: item,
      Url: item.Link,
      Keywords: keywords,
    }
    taskQueue <- task
  }
}

func execTasks(UrlList []DataItem, keyword string) ([]*UrlResult) {
  taskQueue := make(chan *UrlTask, 64)
  resultQueue := make(chan *UrlResult, 64)
  ExecNum := 4
  for i := 0; i < ExecNum; i++ {
    go execTask(taskQueue, resultQueue)
  }
  keywords := wordtool.SplitContent2Words(keyword)
  go submitTask(taskQueue, UrlList, keywords)
  resultList := make([]*UrlResult, 0)
MAIN_END_LABEL:
  for {
END_INSERT:
    select {
    case resultItem := <- resultQueue:
      index := len(resultList)
      for i := 0; i < len(resultList); i++ {
        val, err := resultItem.Compare(resultList[i])
        if err != nil {
          break END_INSERT
        }
        if val > 0 {
          index = i
          break
        } else if val == 0 && resultItem.Task.Url == resultList[i].Task.Url {
          break END_INSERT
        }
      }
      tmpList := resultList
      resultList = make([]*UrlResult, 0)
      for _, item := range tmpList[0:index] {
        resultList = append(resultList, item)
      }
      resultList = append(resultList, resultItem)
      for _, item := range tmpList[index:] {
        resultList = append(resultList, item)
      }
      if len(resultList) >= len(UrlList) {
        break MAIN_END_LABEL
      }
    case <- time.After(10 * time.Second):
      break MAIN_END_LABEL
    }
  }
  return resultList
}
