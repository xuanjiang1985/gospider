package main

import (
	"time"
  "learn/app/library/database"
  "fmt"
  "log"
  "goquery-book/model"
  "github.com/PuerkitoBio/goquery"
  "github.com/djimenez/iconv-go"
  "strconv"
  "strings"
  "net/http"
  "encoding/json"
  "os"
  "io/ioutil"
  "io"
  "bytes"
)

//语文
//人教小学教材
var start = 23095
var end = 23142

func BookScrape() {

  for i := start; i <= end; i++ {
    fmt.Println("正在采集语文页面",i)

    resource := "http://www.yuwenziyuan.com/rjb/6x/dzkb/" + strconv.Itoa(i) + ".html"
    doc, err := goquery.NewDocument(resource)
    if err != nil {
      log.Fatal(err)
    }

    materialName := utfString(doc.Find("div.crumb a").Eq(1).Text())
    gradeAndVolume := utfString(doc.Find("div.crumb a").Eq(2).Text())
    knowledgeName := utfString(doc.Find("div.crumb h1").Eq(0).Text())
    gradeName := ""
    volume := 0
    imgLink := []string{}
    doc.Find("div.dzkb img").Each(func(index int, s *goquery.Selection){
      src, ok := s.Attr("data-original")
      if ok != true {
        fmt.Println("src-err",err)
      } else {
        imgLink = append(imgLink, getImage(src))
      }
      
    })

    imgJson, _ := json.Marshal(imgLink)

    if gradeAndVolume != "" {
      slice := strings.Split(gradeAndVolume, "")
      gradeName = slice[0] + "年级"
      if (slice[1] == "上") {
        volume = 0
      } else {
        volume = 1
      }
    } else {
      gradeName = gradeAndVolume
      volume = 0
    }

    //写入数据库
    ebook := model.MaterialEbookModel{}
    ebook.CourseName = "语文"
    ebook.MaterialName = materialName
    ebook.Volume = int8(volume)
    ebook.GradeName = gradeName
    ebook.KnowledgeName = knowledgeName
    ebook.QueryPage = i
    ebook.ImgLink = string(imgJson)

    db := database.DB()
    db.Create(&ebook)
    time.Sleep(3)
  } 
  fmt.Println("采集完毕。")
}

func main() {
  BookScrape()
}

func utfString(str string) string {
  word, err:= iconv.ConvertString(str,  "gbk","utf-8")
  if err != nil {
    word = ""
    return word
  }
  return word
}

func getImage(src string) string {

    path := strings.Split(src, "/")
    name := strconv.Itoa(int(time.Now().Unix()))
    str := "null.jpg"
    if len(path) > 1 {
        str = path[len(path)-1]
    }
    imgName := name + "-" + str
    out, err := os.Create("images/" + imgName)
    defer out.Close()
    resp, err := http.Get(src)
    if err != nil {
      fmt.Println("获取图片地址错误", err)
      return str
    }
    defer resp.Body.Close()
    byteString, _ := ioutil.ReadAll(resp.Body)
    io.Copy(out, bytes.NewReader(byteString))
    return imgName
}