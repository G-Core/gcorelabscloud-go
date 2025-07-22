package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	"github.com/G-Core/gcorelabscloud-go/testhelper"
)

// OffsetPager sample and test cases.

type OffsetPageResult struct {
	pagination.OffsetPageBase
}

func ExtractOffsetInts(r pagination.Page) ([]int, error) {
	var s struct {
		Results []int `json:"results"`
	}
	err := (r.(OffsetPageResult)).ExtractInto(&s)
	return s.Results, err
}

func createOffsetPager(t *testing.T) pagination.Pager {
	testhelper.SetupHTTP()

	testhelper.Mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		switch offset := r.URL.Query().Get("offset"); offset {
		case "0":
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{ "count": 9, "results": [1, 2, 3] }`)
		case "3":
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{ "count": 9, "results": [4, 5, 6] }`)
		case "6":
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{ "count": 9, "results": [7, 8, 9] }`)
		case "9":
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{ "count": 9, "results": [] }`)
		default:
			t.Errorf("Request with unexpected offset: %v", offset)
		}
	})

	client := createClient()

	createPage := func(r pagination.PageResult) pagination.Page {
		return OffsetPageResult{pagination.OffsetPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, testhelper.Server.URL+"/list?limit=3&offset=0", createPage)
}

func TestEnumerateOffset(t *testing.T) {
	pager := createOffsetPager(t)
	defer testhelper.TeardownHTTP()

	callCount := 0
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		actual, err := ExtractOffsetInts(page)
		if err != nil {
			return false, err
		}

		t.Logf("Handler invoked with %v", actual)

		var expected []int
		switch callCount {
		case 0:
			expected = []int{1, 2, 3}
		case 1:
			expected = []int{4, 5, 6}
		case 2:
			expected = []int{7, 8, 9}
		case 3:
			expected = nil
		default:
			t.Fatalf("Unexpected call count: %d", callCount)
			return false, nil
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Call %d: Expected %#v, but was %#v", callCount, expected, actual)
		}

		callCount++
		return true, nil
	})
	if err != nil {
		t.Errorf("Unexpected error for page iteration: %v", err)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, but was %d", callCount)
	}
}

func TestAllPagesOffset(t *testing.T) {
	pager := createOffsetPager(t)
	defer testhelper.TeardownHTTP()

	page, err := pager.AllPages()
	testhelper.AssertNoErr(t, err)

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	actual, err := ExtractOffsetInts(page)
	testhelper.AssertNoErr(t, err)
	testhelper.CheckDeepEquals(t, expected, actual)
}
