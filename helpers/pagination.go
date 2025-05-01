package helpers

import (
	"net/url"
	"strconv"
)

func GetOffsetLimit(query url.Values) (page, offset, limit int64) {
	// validate offset & limit
	pageString := query.Get("page")
	limitString := query.Get("limit")
	pageInt, _ := strconv.Atoi(pageString)
	limitInt, _ := strconv.Atoi(limitString)

	// convert to int64
	pageInt64 := int64(pageInt)
	limitInt64 := int64(limitInt)

	// set default page & limit
	if pageInt64 <= 0 {
		pageInt64 = 1
	}
	if limitInt64 <= 0 {
		limitInt64 = 10
	}

	// set offset
	offset = (pageInt64 - 1) * limitInt64

	return pageInt64, offset, limitInt64
}
