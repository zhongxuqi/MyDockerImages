package spider

import (
  "fmt"
  "time"
  . "DesertEagleSite/bean"
  "DesertEagleSite/wordtool"
  "DesertEagleSite/evaluator"
  "DesertEagleSite/util"
  "DesertEagleSite/push_manager"
)

var mSearchResultMap = make(map[string] UnionResponse)

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
  for _, item := range UrlList {
    task := &UrlTask{
      Info: item,
      Url: item.Link,
      Keywords: keywords,
    }
    taskQueue <- task
  }
}

func GetUnionData(keyword, parser_names, registration_id string, spiderList []SpiderObject) {
  UrlList := make([]DataItem, 0)
  for _, spider := range spiderList {
    fmt.Println("search in: ", spider.Name)
    itemList, nextPage, err := spider.GetDataFunc(keyword)
    if err == nil && itemList != nil {
      for _, item := range itemList {
        UrlList = append(UrlList, item)
      }
    }

    // get next page data
    itemList, _, err = spider.ParseFunc(nextPage)
    if err == nil && itemList != nil {
      for _, item := range itemList {
        UrlList = append(UrlList, item)
      }
    }
  }
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
  mapKey := registration_id + "-" + util.GetFormatTimeNow()
  var response UnionResponse
  response.ResultData = make([]UrlResult, 0)
  for _, item := range resultList {
    fmt.Println(item.ToString())
    response.ResultData = append(response.ResultData, *item)
  }
	response.Status = "200"
	response.Message = "search success"
  response.ParserNames = parser_names
  response.Keyword = keyword
  mSearchResultMap[mapKey] = response
	var message PushMessage
	message.MapKey = mapKey
  message.Keyword = keyword
	push_manager.PushJPushMessage(registration_id, util.ConvObject2Json(message))
}

func GetResultByKey(mapkey string) UnionResponse {
  resp, ok := mSearchResultMap[mapkey]
  if ok {
    delete(mSearchResultMap, mapkey)
    return resp
  } else {
    var response UnionResponse
    response.Status = "400"
  	response.Message = "has not the map key"
    return response
  }
}
