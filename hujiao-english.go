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

//英语
//沪教版
var url = "http://www.171english.cn/oxford/shoxford/9A/unit"
var host = "http://www.171english.cn/oxford/shoxford/9A/"
func BookScrape() {

  for j := 4; j < 8; j++ {
    fmt.Println("正在采集英语页面unit", j)

	resource := url + strconv.Itoa(j) + ".html"
    doc, err := goquery.NewDocument(resource)
    if err != nil {
      log.Fatal(err)
	}
	
     // 解析 <a href="index.html">三年级上册英语</a>
	 material := utfString(doc.Find("div#nowPosition").Find("div.texts a").Last().Text())
	 materialNameSlice := strings.Split(material, "") 
	 materialName := "沪教版"
	 gradeName := materialNameSlice[0] + materialNameSlice[1] + materialNameSlice[2]
	 volumeName := materialNameSlice[3] + materialNameSlice[4]
	 volume := 0
	 if volumeName == "上册" {
		 volume = 0
	 } else {
	   volume = 1
	 }

	 knowledge := utfString(doc.Find("div#nowPosition").Find("div.texts").Text())
	 knowledgeSlice := strings.Split(knowledge, ">")
	 fmt.Println(knowledgeSlice)
	 knowledgeName := knowledgeSlice[2]

    imgLink := []string{}
    doc.Find("div.Contentbox img").Each(func(index int, s *goquery.Selection){
      src, ok := s.Attr("src")
      if ok != true {
        fmt.Println("src-err",err)
      } else {
        imgLink = append(imgLink, getImage(host + src))
      }   
	})
		
    imgJson, _ := json.Marshal(imgLink)


    //写入数据库
    ebook := model.MaterialEbookModel{}
    ebook.CourseName = "英语"
    ebook.MaterialName = materialName
    ebook.Volume = int8(volume)
    ebook.GradeName = gradeName
    ebook.KnowledgeName = knowledgeName
    ebook.QueryPage = j
    ebook.ImgLink = string(imgJson)

    db := database.DB()
    db.Create(&ebook)
    time.Sleep(3)
  } 
  fmt.Println("采集完毕---------------------------")
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