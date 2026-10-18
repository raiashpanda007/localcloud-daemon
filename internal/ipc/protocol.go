package ipc

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RequestMethod string

const (
	GET    RequestMethod = "get"
	POST   RequestMethod = "post"
	DELETE RequestMethod = "delete"
	PUT    RequestMethod = "put"
)

// TODO:Resouce will be in payload

type RequestType struct {
	ID      string        `json:"id"`
	PID     string        `json:"pid"`
	Method  RequestMethod `json:"method"`
	Payload string        `json:"payload,omitempty"`
}

type ResponseType struct {
	ID      string `json:"id"`
	ReqID   string `json:"reqID"`
	Payload string `json:"payload,omitempty"`
}

func RequestParser[T any](reqStr string) (RequestType, T, error) {

	var reqObj RequestType
	var payloadObj T
	err := json.Unmarshal([]byte(reqStr), &reqObj)

	if err != nil {
		return reqObj, payloadObj, err
	}

	err = json.Unmarshal([]byte(reqObj.Payload), &payloadObj)

	if err != nil {
		return reqObj, payloadObj, err
	}

	return reqObj, payloadObj, nil
}

func ResponseParserObj[T any](reqId string, payload T) (ResponseType, error) {
	var result ResponseType

	payloadByte, err := json.Marshal(payload)

	if err != nil {
		return result, err
	}

	result = ResponseType{
		ID:      uuid.New().String(),
		ReqID:   reqId,
		Payload: string(payloadByte),
	}

	return result, nil
}

func SendResponseMsg(respObj ResponseType) (string, error) {

	result, err := json.Marshal(respObj)

	if err != nil {
		return "", nil
	}

	return string(result), nil

}
