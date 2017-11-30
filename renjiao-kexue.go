package main

import (
	"time"
  "learn/app/library/database"
  "fmt"
  "log"
  "math/rand"
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

//语文 数学 科学 品德与生活 音乐
//沪教小学教材
var start = 6129
var end = 6301
var host = "http://www.szxuexiao.com/"
func BookScrape() {

  for i := start; i <= end; i++ {
    fmt.Println("正在采集页面",i)

    resource := "http://www.szxuexiao.com/keben/html/" + strconv.Itoa(i) + ".html"
    doc, err := goquery.NewDocument(resource)
    if err != nil {
      log.Fatal(err)
    }
	// <h2>小学一年级科学上册1各式各样的服装</h2>
	knowledgeString := utfString(doc.Find("div.RightBox div.title").Find("h2").Eq(0).Text())
	knowledgeSlice := strings.Split(knowledgeString, "")

	knowledgeName := strings.Join(knowledgeSlice[9:len(knowledgeSlice)], "")

	//上下册
	volume := 0
	if strings.Contains(knowledgeString, "上册") {
		volume = 0
	} else {
		volume = 1
	}

	// 年级
	gradeName := ""
	if strings.Contains(knowledgeString, "一年级") {
		gradeName = "一年级"
	} else if strings.Contains(knowledgeString, "二年级") {
		gradeName = "二年级"
	} else if strings.Contains(knowledgeString, "三年级") {
		gradeName = "三年级"
	} else if strings.Contains(knowledgeString, "四年级") {
		gradeName = "四年级"
	} else if strings.Contains(knowledgeString, "五年级") {
		gradeName = "五年级"
	} else {
		gradeName = "六年级"
	}
	
	// 抓取图片
    imgLink := []string{}
    doc.Find("div#content img").Each(func(index int, s *goquery.Selection){
	  src, ok := s.Attr("src")
      if ok != true {
        fmt.Println("src-err",err)
      } else {
		  //判断是否带http
		  if strings.Contains(src, "http"){
			imgLink = append(imgLink, getImage(src))
		  } else {
			imgLink = append(imgLink, getImage(host + src))
		  }
      }
      
    })

    imgJson, _ := json.Marshal(imgLink)
	
    //写入数据库
    ebook := model.MaterialEbookModel{}
    ebook.CourseName = "美术"
    ebook.MaterialName = "人教版"
    ebook.Volume = int8(volume)
    ebook.GradeName = gradeName
    ebook.KnowledgeName = knowledgeName
    ebook.QueryPage = i
    ebook.ImgLink = string(imgJson)

    db := database.DB()
    db.Create(&ebook)
    time.Sleep(3)
  } 
  fmt.Println("采集完毕--------------------------")
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
	
		name := strconv.Itoa(int(time.Now().Unix()))
		randNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNum := strconv.Itoa(randNumber.Intn(10000))
		imgName := name + "-" + randNum + ".jpg"
		str := "null.jpg"
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