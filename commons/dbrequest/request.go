package dbrequest

import (
	"strings"
	"regexp"
)

type operation string

type DBRequest struct {
	op    operation
	user  string
	table string
	key   string
	value string
}

const (
	Unspecified = operation("")
	GetKV       = operation("get")
	DelKV       = operation("delkv")
	AddKV       = operation("add")
	CreateTable = operation("ct")
	DeleteTable = operation("deltable")
	CreateUser  = operation("cu")
	DeleteUser  = operation("deluser")
)

func StringToOperation(s string) operation {
	switch s {
	case string(GetKV):
		return GetKV
	case string(DelKV):
		return DelKV
	case string(AddKV):
		return AddKV
	case string(CreateTable):
		return CreateTable
	case string(DeleteTable):
		return DeleteTable
	case string(CreateUser):
		return CreateUser
	case string(DeleteUser):
		return DeleteUser
	default:
		return Unspecified
	}
}

func NewRequest(s string) DBRequest {
	req := DBRequest{}
	operation := Unspecified
	params := []string{}

	// this pattern is chatgpt'd, i couldnt be bothered to do this lol
	pattern := `^(\w+)\s+((?:\s*(?:"[^"\\]*(?:\\.[^"\\]*)*"|\S+))*)\s*;?$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(s)
	if len(matches) > 2 {
		operation = StringToOperation(matches[1])
		req.op = operation
		paramsBlob := matches[2]

		paramsMatches := regexp.MustCompile(`"[^"\\]*(?:\\.[^"\\]*)*"|\S+`).FindAllString(paramsBlob, -1)
		for _, param := range paramsMatches {
			param = strings.Trim(param, `"`)
			param = strings.Replace(param, `\"`, `"`, -1)
			param = strings.Replace(param, `\\`, `\`, -1)
			params = append(params, param)
		}
	}

	paramFields := []*string{&req.user, &req.table, &req.key, &req.value}
	for i, param := range params {
		if i < len(paramFields) {
			*paramFields[i] = param
		}
	}

	return req
}

func (r DBRequest) Validate() bool {
	switch r.op {
	case GetKV, DelKV:
		return r.user != "" && r.table != "" && r.key != ""
	case AddKV:
		// r.value != "" is not included as the value could be the null string, ""
		return r.user != "" && r.table != "" && r.key != ""
	case CreateTable, DeleteTable:
		return r.user != "" && r.table != ""
	case CreateUser, DeleteUser:
		return r.user != ""
	default:
		return false
	}
}