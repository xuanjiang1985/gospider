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
  "math/rand"
)

//语文 英语 物理 历史
//北师大版
var start = 18
var listHost = "http://gbjc.bnup.com/eduresource.php?action=showcatalog&subjectid=2328&id=2355"
var host = "http://gbjc.bnup.com"
func BookScrape() {
	// 抓取小学所有目录放入 切片
	primaryPage := make([]map[string]string, 0 , 20)
	resource := listHost
    doc, err := goquery.NewDocument(resource)
    if err != nil {
      log.Fatal(err)
	}
	
	doc.Find("div#container_txtimg table").Find("li a").Each(func(index int, s *goquery.Selection){
		eleText := s.Find("span").Text()
		eleLink, _ := s.Attr("href")
		result := make(map[string]string , 2)
		result["knowledge"] = eleText
		result["link"] = eleLink
		primaryPage = append(primaryPage, result)
	  })
	  // 抓取年级上下册
	  material := doc.Find("div.home_title").Find("div.home_title_name a").Text()
	  //开始采集
	  fmt.Println("一共收集到采集页面map总数:", len(primaryPage))
	  restPage := primaryPage[start:len(primaryPage)]
	  for index, value := range restPage {
				start = index
				fmt.Println("开始采集页面map：", start)
				volume, gradeName := collectText(material)
				imgJson := collectImage(host + "/" + value["link"])
				 //写入数据库
				 ebook := model.MaterialEbookModel{}
				 ebook.CourseName = "历史"
				 ebook.MaterialName = "北师大"
				 ebook.Volume = volume
				 ebook.GradeName = gradeName
				 ebook.KnowledgeName = value["knowledge"]
				 ebook.QueryPage = start
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

func collectText(str string) (volume int8, gradeName string) {
	// <a>苏教版一年级语文上册——目录</a> 或者 <a>苏教版二年级上册语文：识字4</a> 或者 <a>小学六年级苏教版语文下册 生字表</a>
	volume = 0
	
	// 上下册
	if strings.Contains(str, "上") {
		volume = int8(0)
	} else {
		volume = int8(1)
	}
	// 年级
	if strings.Contains(str, "七年级") {
		gradeName = "七年级"
	} else if strings.Contains(str, "八年级") {
		gradeName = "八年级"
	}else {
		gradeName = "九年级"
	}
	return int8(volume), gradeName
}

func collectImage(src string) []byte {
	doc, err := goquery.NewDocument(src)
	if err != nil {
		log.Fatal(err)
	}
	// 抓取图片
	imgLink := []string{}
	// 第一页
	doc.Find("div#gallery").Find("div.ad-nav ul li").Find("a img").Each(func(index int, s *goquery.Selection){
		imgsrc, ok := s.Attr("src")
		if ok != true {
		fmt.Println("img-src-err", err)
		} else {
		imgLink = append(imgLink, getImage(host + imgsrc))
		} 
	})
	//判断是否为空
	if len(imgLink)  == 0 {
		img := doc.Find("div.min_right").Find("div.home_c").Children().Eq(3).Find("img").First();
		imgsrc, ok := img.Attr("src")
		if ok != true {
		fmt.Println("img-src-err2", err)
		} else {
		imgLink = append(imgLink, getImage(host + imgsrc))
		} 
	}
	imgJson, _ := json.Marshal(imgLink)
	return imgJson
}