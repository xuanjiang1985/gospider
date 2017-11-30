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

//物理
//人教小学教材
var start = 3115
var end = 3124
var url = "http://www.wsbedu.com/wu51/keben.asp?wai="
var host = "http://www.wsbedu.com"
func BookScrape() {

  for i := start; i <= end; i++ {
    fmt.Println("正在采集物理页面", i)

	resource := url + strconv.Itoa(i)
    doc, err := goquery.NewDocument(resource)
    if err != nil {
      log.Fatal(err)
	}
	
     // 解析 <a href="/xiaox/rjs11kb.asp">人教版一年级数学上册电子课本</a>
	 material := utfString(doc.Find("div.cle").Find("div.ggao2 a").Last().Text())
	 materialNameSlice := strings.Split(material, "") 
	 materialName := materialNameSlice[0] + materialNameSlice[1] + materialNameSlice[2]
	 gradeName := materialNameSlice[3] + materialNameSlice[4] + materialNameSlice[5]
	 volumeName := materialNameSlice[8] + materialNameSlice[9]
	 courseName := materialNameSlice[6] + materialNameSlice[7]
	 volume := 0
	 if volumeName == "上册" {
		 volume = 0
	 } else {
	   volume = 1
	 }

	 knowledge := utfString(doc.Find("div.main_left2 h4").First().Text())
	 knowledgeSlice := strings.Split(knowledge, "　")
	 fmt.Println(knowledgeSlice)
	 knowledgeName := knowledgeSlice[0]

    imgLink := []string{}
    doc.Find("div.main_left2 img").Each(func(index int, s *goquery.Selection){
      src, ok := s.Attr("src")
      if ok != true {
        fmt.Println("src-err",err)
      } else {
        imgLink = append(imgLink, getImage(host + src))
      }   
	})
	
	//判断是否有下一页
	pages := doc.Find("div.main table").Find("font a").Length()
	pages = pages -1
	if (pages > 1) {
		for j := 2; j <= pages; j++ {
			resource = url + strconv.Itoa(i) + "&page=" + strconv.Itoa(j)
			doc2, err := goquery.NewDocument(resource)
			if err != nil {
			  log.Fatal(err)
			}

			doc2.Find("div.main_left2 img").Each(func(index int, s *goquery.Selection){
				src, ok := s.Attr("src")
				if ok != true {
				  fmt.Println("src-err",err)
				} else {
				  imgLink = append(imgLink, getImage(host + src))
				}   
			  })

		}
	}

    imgJson, _ := json.Marshal(imgLink)


    //写入数据库
    ebook := model.MaterialEbookModel{}
    ebook.CourseName = courseName
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