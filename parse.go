package flickr

import (
	"encoding/json"
	"errors"
)

type FailResponse struct {
	Stat    string
	Code    int
	Message string
}

func Parse(data []byte, v interface{}) error {
	fail := Fail(data)

	if fail != nil {
		return fail
	}

	err := json.Unmarshal(data, v)

	if err != nil {
		fail := Fail(data)

		if fail != nil {
			return fail
		}

		return err
	}

	return nil
}

type PaginatedResult struct {
	Page       int `json:"page"`
	NumPages   int `json:"pages"`
	NumPerPage int `json:"perpage"`
	Total      int `json:"total"`
}

func ParsePaginated(data []byte, p *PaginatedResult, v interface{}) error {
	fail := Fail(data)

	if fail != nil {
		return fail
	}

	err := json.Unmarshal(data, v)

	if err != nil {
		fail := Fail(data)

		if fail != nil {
			return fail
		}

		return err
	}

	// err = json.Unmarshal(data, p)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func Fail(data []byte) error {
	fail := &FailResponse{}
	err := json.Unmarshal(data, fail)

	if err == nil && fail.Stat == "fail" {
		return errors.New(fail.Message)
	}

	return nil
}
