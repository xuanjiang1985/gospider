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
//翼教版
var lesson = [][]int{{1,2,3,4,5,6}, {1,2,3,4,5,6}, {1,2,3,4,5,6}, {1,2,3,4,5,6}}
var host = "http://www.171english.cn/yijiao/xiaoxue/6B/unit"
func BookScrape() {
for j := 1; j <= 4; j++ {
	for index, value := range lesson[j-1] {
		fmt.Println("正在采集英语页面lesson", value)

		resource := host + strconv.Itoa(j) + "/lesson" + strconv.Itoa(value) + ".html"
		doc, err := goquery.NewDocument(resource)
		if err != nil {
		log.Fatal(err)
		}
		
		// 解析 
		material := utfString(doc.Find("div#nowPosition").Find("div.texts a").First().Text())
		materialName := "冀教版"
		gradeName := ""
		volume := 0
		// 上下册
		if strings.Contains(material, "上") {
			volume = 0
		} else {
			volume = 1
		}
		// 年级
		if strings.Contains(material, "一年级") {
			gradeName = "一年级"
		} else if strings.Contains(material, "二年级") {
			gradeName = "二年级"
		} else if strings.Contains(material, "三年级") {
			gradeName = "三年级"
		} else if strings.Contains(material, "四年级") {
			gradeName = "四年级"
		} else if strings.Contains(material, "五年级") {
			gradeName = "五年级"
		} else {
			gradeName = "六年级"
		}

		knowledge := utfString(doc.Find("div#nowPosition").Find("div.texts").Text())
		knowledgeSlice := strings.Split(knowledge, ">")
		fmt.Println(knowledgeSlice)
		knowledgeName := knowledgeSlice[len(knowledgeSlice) - 1]

		imgLink := []string{}
		doc.Find("div.Contentbox img").Each(func(index int, s *goquery.Selection){
		src, ok := s.Attr("src")
		if ok != true {
			fmt.Println("src-err",err)
		} else {
			imgLink = append(imgLink, getImage(host + strconv.Itoa(j) + "/" + src))
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
		ebook.QueryPage = index
		ebook.ImgLink = string(imgJson)

		db := database.DB()
		db.Create(&ebook)
		time.Sleep(3)
	} 
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