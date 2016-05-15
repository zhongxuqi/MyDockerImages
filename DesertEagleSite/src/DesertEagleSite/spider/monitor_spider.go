package spider

import (
  "fmt"
  "time"
  "sync"
  . "DesertEagleSite/bean"
  "DesertEagleSite/util"
  "DesertEagleSite/wordtool"
  "DesertEagleSite/evaluator"
  "github.com/PuerkitoBio/goquery"
  "DesertEagleSite/push_manager"
)

var mutex = &sync.Mutex{}
var mTaskList = make([]*MonitorTask, 0)
var mMonitorResultMap = make(map[string] MonitorResponse)

func init() {
  go loop()
}

func loop() {
  for {
    time.Sleep(30 * time.Second)

    for _, task := range mTaskList {

      // get response struct
      resp, ok := mMonitorResultMap[task.RespMapKey]
      if !ok {
        continue
      }
      resp.ResultData = make([]UrlResult, 0)

      // eval main url
      subtask := UrlResult{}
      subtask.Task.Url = task.Url
      subtask.Task.Keywords = task.Keywords
      subtask.Task.Info.Title = "Main Page"
      subtask.Task.Info.Link = task.Url
      doc, err := goquery.NewDocument(subtask.Task.Url)
    	if err != nil {
    		continue
    	}
      subtask.Eval = evaluator.EvaluateContentByKeyWords(doc.Find("body").Text(), subtask.Task.Keywords)
      resp.ResultData = append(resp.ResultData, subtask)

      // parser a tag from url
      UrlList := make([]DataItem, 0)
      doc.Find("a").Each(func(i int, s *goquery.Selection) {
        subtask := DataItem{}
        subtask.Title = s.Text()
        subtask.Link = s.First().AttrOr("href", "")
        if len(subtask.Link) == 0 {
          return
        }
        UrlList = append(UrlList, subtask)
      })
      fmt.Println(subtask.Task.Info.Title, " size: ", len(UrlList))

      // get result list
      resultList := execTasks(UrlList, task.Keyword)

      // store the response to map
      mapKey := task.RegistrationId + "-" + util.GetFormatTimeNow()
      var response MonitorResponse
      response.ResultData = make([]UrlResult, 0)
      for _, item := range resultList {
        fmt.Println(item.ToString())
        response.ResultData = append(response.ResultData, *item)
      }
    	response.Status = "200"
    	response.Message = "search success"
      response.Task = *task
      mMonitorResultMap[mapKey] = response

      // notice the client
    	var message PushMessage
    	message.MapKey = mapKey
      message.Keyword = task.Keyword
      message.Type = MONITOR_TYPE
    	push_manager.PushJPushMessage(task.RegistrationId, util.ConvObject2Json(message))
    }
  }
}

func SubmitMonitorTask(task *MonitorTask) bool {
  for _, item := range mTaskList {
    if task.IsEqual(item) {
      return false
    }
  }

  // add task to list
  mutex.Lock()
  task.Keywords = wordtool.SplitContent2Words(task.Keyword)
  mTaskList = append(mTaskList, task)
  mutex.Unlock()

  // add response to list
  mapKey := task.RegistrationId + "-" + util.GetFormatTimeNow()
  task.RespMapKey = mapKey
  var response MonitorResponse
  response.Status = "200"
	response.Message = "search success"
  response.Task = *task
  mMonitorResultMap[mapKey] = response

  return true
}

func SubmitRawMonitorTask(keyword, registration_id, target_url string) {
  task := &MonitorTask{}
  task.Url = target_url
  task.Keyword = keyword
  task.RegistrationId = registration_id
  SubmitMonitorTask(task)
}

func GetMonitorResultByKey(mapkey string) MonitorResponse {
  resp, ok := mMonitorResultMap[mapkey]
  if ok {
    delete(mMonitorResultMap, mapkey)
    return resp
  } else {
    var response MonitorResponse
    response.Status = "400"
  	response.Message = "has not the map key"
    return response
  }
}
