package network

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
	Get         = operation("get")
	Del         = operation("del")
	Add         = operation("add")
	CreateTable = operation("ct")
	CreateUser  = operation("cu")
)

func StringToOperation(s string) {
	switch s {
	case str(Get):
		return Get
	case str(Del):
		return Del
	case str(Add):
		return Add
	case str(CreateTable):
		return CreateTable
	case str(CreateUser):
		return CreateUser
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

	matches := re.FindStringSubmatch(request)
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
	case Get, Del:
		return r.user != "" && r.table != "" && r.key != ""
	case Add:
		// r.value != "" is not included as the value could be the null string, ""
		return r.user != "" && r.table != "" && r.key != ""
	case CreateTable:
		return r.user != "" && r.table != ""
	case CreateUser:
		return r.user != ""
	default:
		return false
	}
}
