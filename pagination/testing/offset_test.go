package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

// createNoLimitOffsetPager mimics a list endpoint queried without a limit
// query parameter: the API returns the whole collection (here 25 items, more
// than the legacy assumed page size of 10) in a single response.
func createNoLimitOffsetPager(t *testing.T, requestCount *int) pagination.Pager {
	testhelper.SetupHTTP()

	all := make([]string, 0, 25)
	for i := 1; i <= 25; i++ {
		all = append(all, strconv.Itoa(i))
	}

	testhelper.Mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		*requestCount++
		// The whole collection is returned only when no limit is requested.
		// Any follow-up paged request (the bug) is treated as unexpected.
		if r.URL.Query().Get("limit") != "" || r.URL.Query().Get("offset") != "" {
			t.Errorf("unexpected paged request: %v", r.URL.RawQuery)
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "count": 25, "results": [%s] }`, strings.Join(all, ", "))
	})

	client := createClient()

	createPage := func(r pagination.PageResult) pagination.Page {
		return OffsetPageResult{pagination.OffsetPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, testhelper.Server.URL+"/list", createPage)
}

// TestAllPagesOffsetNoLimitSingleRequest guards against the regression where a
// limit-less list whose full result set exceeds the default page size of 10
// triggered redundant offset/limit follow-up requests (and duplicate results).
func TestAllPagesOffsetNoLimitSingleRequest(t *testing.T) {
	requestCount := 0
	pager := createNoLimitOffsetPager(t, &requestCount)
	defer testhelper.TeardownHTTP()

	page, err := pager.AllPages()
	testhelper.AssertNoErr(t, err)

	actual, err := ExtractOffsetInts(page)
	testhelper.AssertNoErr(t, err)

	expected := make([]int, 0, 25)
	for i := 1; i <= 25; i++ {
		expected = append(expected, i)
	}
	testhelper.CheckDeepEquals(t, expected, actual)

	if requestCount != 1 {
		t.Errorf("expected exactly 1 request, but made %d", requestCount)
	}
}
