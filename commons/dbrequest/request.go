package dbrequest

import (
	"regexp"
	"strings"
)

type operation string

type DBRequest struct {
	Op    operation
	User  string
	Table string
	Key   string
	Value string
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
		req.Op = operation
		paramsBlob := matches[2]

		paramsMatches := regexp.MustCompile(`"[^"\\]*(?:\\.[^"\\]*)*"|\S+`).FindAllString(paramsBlob, -1)
		for _, param := range paramsMatches {
			param = strings.Trim(param, `"`)
			param = strings.Replace(param, `\"`, `"`, -1)
			param = strings.Replace(param, `\\`, `\`, -1)
			param = strings.Replace(param, `\;`, ";", -1)
			params = append(params, param)
		}
	}

	paramFields := []*string{&req.User, &req.Table, &req.Key, &req.Value}
	for i, param := range params {
		if i < len(paramFields) {
			// removing trailing semicolon on last parameter
			if i == len(params)-1 {
				param = strings.TrimSuffix(param, ";")
			}
			*paramFields[i] = param
		}
	}

	return req
}

func (r DBRequest) Validate() bool {
	switch r.Op {
	case GetKV, DelKV:
		return r.User != "" && r.Table != "" && r.Key != ""
	case AddKV:
		// r.value != "" is not included as the value could be the null string, ""
		return r.User != "" && r.Table != "" && r.Key != ""
	case CreateTable, DeleteTable:
		return r.User != "" && r.Table != ""
	case CreateUser, DeleteUser:
		return r.User != ""
	default:
		return false
	}
}
