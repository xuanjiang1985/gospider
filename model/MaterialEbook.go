package model

import "learn/app/model"

type MaterialEbookModel struct {
	model.BaseModel
	Id int64 `json:"id" gorm:"primary_key"`
	CourseName string `json:"course_name"`
	MaterialName string `json:"meterial_name"`
	Volume int8 `json:"volume"`
	GradeName string `json:"grade_name"`
	KnowledgeName string `json:"knowledge_name"`
	QueryPage int `json:"query_page"`
	ImgLink string `jsong:"img_link"`
	Path string `json:"path"`
}

func (this *MaterialEbookModel)TableName() string{
	return this.GetTableName("material_ebook")
}