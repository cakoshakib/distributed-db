package joinrequest

import (
	"regexp"
)

const (
	Join = "join"
)

type JoinRequest struct {
	NodeID  string
	Address string
}

func NewRequest(s string) JoinRequest {
	req := JoinRequest{}

	re := regexp.MustCompile(`^join\s+([^\s]+)\s+(.+);$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) == 3 {
		req.NodeID = matches[1]
		req.Address = matches[2]
	}

	return req
}

func (r JoinRequest) Validate() bool {
	return r.NodeID != "" && r.Address != ""
}
