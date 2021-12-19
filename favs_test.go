//go:build unit
// +build unit

package flickr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFavsFirstPage(t *testing.T) {
	client, err := NewPhotosClient()
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var sampleResp = []byte(`{
			"photos": {
			  "page": 1,
			  "pages": 43,
			  "perpage": 10,
			  "total": 421,
			  "photo": [
				{
				  "id": "51737769056",
				  "owner": "38134034@N04",
				  "secret": "0a0e8bf24c",
				  "server": "65535",
				  "farm": 66,
				  "title": "Take a seat",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639261222"
				},
				{
				  "id": "51739666614",
				  "owner": "80728884@N06",
				  "secret": "5388db9d26",
				  "server": "65535",
				  "farm": 66,
				  "title": "If things are shaking, you need a Friend...",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639261194"
				},
				{
				  "id": "51739069301",
				  "owner": "92986809@N05",
				  "secret": "c218dd924d",
				  "server": "65535",
				  "farm": 66,
				  "title": "Zooks Mill",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639261169"
				},
				{
				  "id": "51738575977",
				  "owner": "30430801@N06",
				  "secret": "903639dd92",
				  "server": "65535",
				  "farm": 66,
				  "title": "Sirmione, Lago di Garda",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639261159"
				},
				{
				  "id": "51179405040",
				  "owner": "53033551@N03",
				  "secret": "de2d6b3faa",
				  "server": "65535",
				  "farm": 66,
				  "title": "Berlin | (2015)",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639255345"
				},
				{
				  "id": "51646533603",
				  "owner": "53033551@N03",
				  "secret": "47a4d4923a",
				  "server": "65535",
				  "farm": 66,
				  "title": "DSC02170",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639255225"
				},
				{
				  "id": "50657384952",
				  "owner": "190970715@N04",
				  "secret": "46477d87ca",
				  "server": "65535",
				  "farm": 66,
				  "title": "Dans le Santerre, sans terre, pas de pommes de terre",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639244884"
				},
				{
				  "id": "51733295191",
				  "owner": "190970715@N04",
				  "secret": "eb954a0cb4",
				  "server": "65535",
				  "farm": 66,
				  "title": "Brumes sur la mégalopole (explore)",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639239960"
				},
				{
				  "id": "51734204877",
				  "owner": "76093456@N04",
				  "secret": "4f58067045",
				  "server": "65535",
				  "farm": 66,
				  "title": "IXPE by SpaceX",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639239570"
				},
				{
				  "id": "51735604885",
				  "owner": "28497307@N05",
				  "secret": "0e0c1642c4",
				  "server": "65535",
				  "farm": 66,
				  "title": "Making Haste",
				  "ispublic": 1,
				  "isfriend": 0,
				  "isfamily": 0,
				  "date_faved": "1639239417"
				}
			  ]
			},
			"stat": "ok"
		  }`)
		compactedResp := new(bytes.Buffer)
		if err := json.Compact(compactedResp, sampleResp); err != nil {
			t.Fatalf("Unable to compact JSON: %s", err)
		}
		rw.Write(compactedResp.Bytes())
	}))
	defer server.Close()

	client.URL = server.URL
	favs, err := client.Favs("")
	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	if favs.Page != 1 {
		t.Fatalf("Expecting page to be 1, but got %d", favs.Page)
	}
	if favs.Pages != 43 {
		t.Fatalf("Expecting at 43 pages, but got %d", favs.Pages)
	}
	if favs.PerPage != 10 {
		t.Fatalf("Expecting 100 photos per page, but got %d", favs.PerPage)

	}
	if len(favs.Photos) != 10 {
		t.Fatalf("Expecting 10 photos, got %d", len(favs.Photos))
	}
}
