package feeder

import (
	"testing"
	"time"
)

func Test_InvalidDate(t *testing.T) {
	date, err := parseTime("invalid")
	if !date.IsZero() {
		t.Errorf("Invalid date should parse to zero")
	}
	if err == nil {
		t.Errorf("error should not be nil")
	}
}

func Test_ParseLayout0(t *testing.T) {
	date, err := parseTime("2014-03-07T05:38:00-05:00")
	expected := time.Date(2014, time.March, 7, 5, 38, 0, 0, time.FixedZone("-0500", -18000))
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout1(t *testing.T) {
	date, err := parseTime("Fri, 07 Mar 2014 17:42:51 GMT")
	expected := time.Date(2014, time.March, 7, 17, 42, 51, 0, time.UTC)
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout2(t *testing.T) {
	date, err := parseTime("2014-02-05T23:33:34Z")
	expected := time.Date(2014, time.February, 5, 23, 33, 34, 0, time.UTC)
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout3(t *testing.T) {
	date, err := parseTime("Mon, 03 Mar 2014 02:12:25 +0000")
	expected := time.Date(2014, time.March, 3, 2, 12, 25, 0, time.UTC)
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout4(t *testing.T) {
	date, err := parseTime("Fri, 21, Mar 2014 10:41")
	expected := time.Date(2014, time.March, 21, 10, 41, 0, 0, time.UTC)
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout4_1(t *testing.T) {
	date, err := parseTime("Fri, 17, Jan 2014 11:1")
	expected := time.Date(2014, time.January, 17, 11, 1, 0, 0, time.UTC)
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout4_2(t *testing.T) {
	date, err := parseTime("Thu, 9, Jan 2014 10:19")
	expected := time.Date(2014, time.January, 9, 10, 19, 0, 0, time.UTC)
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func Test_ParseLayout5(t *testing.T) {
	date, err := parseTime("22 Jul 2013 14:55:01 EST")
	expected := time.Date(2013, time.July, 22, 14, 55, 1, 0, time.FixedZone("EST", -18000))
	assertEqualTime(t, expected, date)
	if err != nil {
		t.Errorf("err should be nil")
	}
}

func assertEqualTime(t *testing.T, expected, actual time.Time) {
	if !expected.Equal(actual) {
		t.Errorf("expected %v but was %v", expected, actual)
	}
}
