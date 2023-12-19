package urlquery

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func ReadStringFromQuery(values url.Values, key string) *string {
	if values.Has(key) {
		str := values.Get(key)
		return &str
	}

	return nil
}

func ReadInt64FromQuery(values url.Values, key string) (*int64, error) {
	if values.Has(key) {
		str := values.Get(key)
		integer, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse '%s', expected an integer, got '%s'", key, str)
		}

		return &integer, nil
	}

	return nil, nil
}

func ReadIntFromQuery(values url.Values, key string) (*int, error) {
	if values.Has(key) {
		str := values.Get(key)
		integer, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("could not parse '%s', expected an integer, got '%s'", key, str)
		}

		return &integer, nil
	}

	return nil, nil
}

func ReadDateFromQuery(values url.Values, key string, format string) (*time.Time, error) {
	if values.Has(key) {
		date, err := time.Parse(format, values.Get(key))
		if err != nil {
			return nil, err
		}

		return &date, nil
	}

	return nil, nil
}
