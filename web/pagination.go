package web

import (
	"fmt"
	"net/http"
	"strconv"
)

type PageResult struct {
	DisplayCount int
	HrefCount int

	Pages int
	CurrentPage int
}

func getPageFromURL(r *http.Request, itemCount int, countDefault int) *PageResult {
	var pr PageResult

	// get display count
	pr.DisplayCount = countDefault
	pr.HrefCount = 0
	if qCount, ok := r.URL.Query()["count"]; ok {
		if len(qCount[0]) >= 1 {
			uCount, err := strconv.ParseInt(qCount[0], 10, 32)
			if err != nil {
				logger.Debugf("invalid count: %s", qCount[0])
			} else {
				pr.DisplayCount = int(uCount)
				pr.HrefCount = int(uCount)
			}
		}
	}

	pr.Pages = roundUp(float64(itemCount) / float64(pr.DisplayCount))

	// get display page
	pr.CurrentPage = 1
	if qPage, ok := r.URL.Query()["page"]; ok {
		if len(qPage[0]) >= 1 {
			uPage, err := strconv.ParseInt(qPage[0], 10, 32)
			if err != nil {
				logger.Debugf("invalid page: %s", qPage[0])
			} else {
				pr.CurrentPage = int(uPage)
			}
		}
	}

	return &pr
}

func makePagination(page, count int, href string, hrefCount int) *[]templatePaginationItems {
	displayItems := 5
	startingNumber := 1

	if count < displayItems {
		// less than
		displayItems = count
	} else if page > count-displayItems/2 {
		// end of the
		startingNumber = count-displayItems+1
	} else if page > displayItems/2 {
		// center active
		startingNumber = page - displayItems/2
	}

	var items []templatePaginationItems

	// previous button
	prevItem := templatePaginationItems{
		Text:        "Previous",
		DisplayHTML: "<i class=\"fas fa-caret-left\"></i>",
	}
	if page == 1 {
		prevItem.Disabled = true
	} else if hrefCount > 0 {
		prevItem.HRef = fmt.Sprintf("%s?page=%d&count=%d", href, page-1, hrefCount)
	} else {
		prevItem.HRef = fmt.Sprintf("%s?page=%d", href, page-1)
	}
	items = append(items, prevItem)

	// add pages
	for i := 0; i < displayItems; i++ {
		newItem := templatePaginationItems{
			Text: fmt.Sprintf("%d", startingNumber+i),
		}

		if page == startingNumber+i {
			newItem.Active = true
		} else if hrefCount > 0 {
			newItem.HRef = fmt.Sprintf("%s?page=%d&count=%d", href, startingNumber+i, hrefCount)
		} else {
			newItem.HRef = fmt.Sprintf("%s?page=%d", href, startingNumber+i)
		}

		items = append(items, newItem)
	}

	// next button
	nextItem := templatePaginationItems{
		Text:        "Next",
		DisplayHTML: "<i class=\"fas fa-caret-right\"></i>",
	}
	if page == count {
		nextItem.Disabled = true
	} else if hrefCount > 0 {
		nextItem.HRef = fmt.Sprintf("%s?page=%d&count=%d", href, page+1, hrefCount)
	} else {
		nextItem.HRef = fmt.Sprintf("%s?page=%d", href, page+1)
	}
	items = append(items, nextItem)

	return &items

}

