package flickr

import (
	"fmt"
)

func ExampleRequest() {
	client := NewClient("", "")

	resp, err := client.Request("people.findByUsername", Params{
		"username": "azerbike",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v", string(resp))
	// Output: {"user":{"id":"98269877@N00","nsid":"98269877@N00","username":{"_content":"azerbike"}},"stat":"ok"}
}
