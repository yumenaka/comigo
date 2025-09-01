package model

import (
	"encoding/json"
)

type ReadingProgress struct {
	// 当前页
	NowPageNum int `json:"nowPageNum"`
	// 当前章节
	NowChapterNum int `json:"nowChapterNum"`
	// 阅读时间，单位为秒
	ReadingTime int `json:"readingTime"`
}

func GetReadingProgress(progress string) (ReadingProgress, error) {
	// 创建一个ReadingProgress实例用于保存解析结果
	var rp ReadingProgress
	// 将JSON字符串解析到ReadingProgress结构体中
	err := json.Unmarshal([]byte(progress), &rp)
	if err != nil {
		// 如果解析出错，返回空的ReadingProgress和错误信息
		return ReadingProgress{}, err
	}
	// 返回解析后的ReadingProgress和nil错误
	return rp, nil
}
