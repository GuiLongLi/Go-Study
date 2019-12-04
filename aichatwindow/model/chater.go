package model

import (
	"encoding/json"

	"aichatwindow/pkg/errno"
)

/*
传入参数 json
{
	"reqType":0,
    "perception": {
        "inputText": {
            "text": "附近的酒店"
        },
        "inputImage": {
            "url": "imageUrl"
        },
        "selfInfo": {
            "location": {
                "city": "北京",
                "province": "北京",
                "street": "信息路"
            }
        }
    },
    "userInfo": {
        "apiKey": "",
        "userId": ""
    }
}
*/

type Chatting struct {
	ReqType    int        `json:"reqType"`
	Perception Perception `json:"perception"`
	UserInfo   UserInfo   `json:"userInfo"`
}

type InputText struct {
	Text string `json:"text"`
}

type Perception struct {
	InputText InputText `json:"inputText"`
}

type UserInfo struct {
	ApiKey string `json:"apiKey"`
	UserId string `json:"userId"`
}
//更新 图灵参数构造体 信息
func UpdateChatting(userId string, text string, chattingInfo Chatting) Chatting {
	chattingInfo.UserInfo.UserId = userId
	chattingInfo.Perception.InputText.Text = text
	return chattingInfo
}
//建立 图灵参数构造体
func BuildChatting(text string, userId string,appKey string) Chatting {
	chatting := Chatting{ReqType: 0}
	chatting.Perception = buildPerception(text)
	chatting.UserInfo = buildUserInfo(userId,appKey)
	return chatting
}
//建立 Perception
func buildPerception(text string) Perception {
	perception := Perception{buildInputText(text)}
	return perception
}
//建立 InputText
func buildInputText(text string) InputText {
	inputText := InputText{text}
	return inputText
}
//建立 UserInfo
func buildUserInfo(userId string,appKey string) UserInfo {
	return UserInfo{appKey, userId}
}

//构造体转换成字符串
func ConvertJson(chattingInfo Chatting) (string,error) {
	jsons, errs := json.Marshal(chattingInfo)
	if errs != nil {
		return "", errno.ModelError
	}
	return string(jsons),nil
}