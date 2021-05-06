package dcmtime_test

import (
	"errors"
	"github.com/suyashkumar/dicom/pkg/dcmtime"
	"testing"
	"time"
)

func TestParseDatetime(t *testing.T) {
	testCases := []struct {
		Name              string
		DTValue           string
		Expected          time.Time
		ExpectedPrecision dcmtime.PrecisionLevel
		HasOffset         bool
	}{
		{
			Name:              "PrecisionFull-PositiveOffset",
			DTValue:           "10100203040506.456789+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionFull,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS5-PositiveOffset",
			DTValue:           "10100203040506.45678+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456780000, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMS5,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS4-PositiveOffset",
			DTValue:           "10100203040506.4567+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456700000, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMS4,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS3-PositiveOffset",
			DTValue:           "10100203040506.456+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456000000, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMS3,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS2-PositiveOffset",
			DTValue:           "10100203040506.45+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 450000000, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMS2,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS1-PositiveOffset",
			DTValue:           "10100203040506.4+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 400000000, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMS1,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionSeconds-PositiveOffset",
			DTValue:           "10100203040506+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 0, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionSeconds,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMinutes-PositiveOffset",
			DTValue:           "101002030405+0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 0, 0, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMinutes,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionHours-PositiveOffset",
			DTValue:           "1010020304+0102",
			Expected:          time.Date(1010, 2, 3, 4, 0, 0, 0, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionHours,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionDay-PositiveOffset",
			DTValue:           "10100203+0102",
			Expected:          time.Date(1010, 2, 3, 0, 0, 0, 0, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionDay,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMonth-PositiveOffset",
			DTValue:           "101002+0102",
			Expected:          time.Date(1010, 2, 1, 0, 0, 0, 0, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionMonth,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionYear-PositiveOffset",
			DTValue:           "1010+0102",
			Expected:          time.Date(1010, 1, 1, 0, 0, 0, 0, time.FixedZone("", 3720)),
			ExpectedPrecision: dcmtime.PrecisionYear,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionFull-NegativeOffset",
			DTValue:           "10100203040506.456789-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionFull,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS5-NegativeOffset",
			DTValue:           "10100203040506.45678-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456780000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMS5,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS4-NegativeOffset",
			DTValue:           "10100203040506.4567-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456700000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMS4,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS3-NegativeOffset",
			DTValue:           "10100203040506.456-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMS3,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS2-NegativeOffset",
			DTValue:           "10100203040506.45-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 450000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMS2,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMS1-NegativeOffset",
			DTValue:           "10100203040506.4-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 400000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMS1,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionSeconds-NegativeOffset",
			DTValue:           "10100203040506-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 000000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionSeconds,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMinutes-NegativeOffset",
			DTValue:           "101002030405-0102",
			Expected:          time.Date(1010, 2, 3, 4, 5, 0, 000000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMinutes,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionHours-NegativeOffset",
			DTValue:           "1010020304-0102",
			Expected:          time.Date(1010, 2, 3, 4, 0, 0, 000000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionHours,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionDay-NegativeOffset",
			DTValue:           "10100203-0102",
			Expected:          time.Date(1010, 2, 3, 0, 0, 0, 000000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionDay,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionMonth-NegativeOffset",
			DTValue:           "101002-0102",
			Expected:          time.Date(1010, 2, 1, 0, 0, 0, 000000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionMonth,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionYear-NegativeOffset",
			DTValue:           "1010-0102",
			Expected:          time.Date(1010, 1, 1, 0, 0, 0, 000000000, time.FixedZone("", -3720)),
			ExpectedPrecision: dcmtime.PrecisionYear,
			HasOffset:         true,
		},
		{
			Name:              "PrecisionFull-NoOffset",
			DTValue:           "10100203040506.456789",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionFull,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionMS5-NoOffset",
			DTValue:           "10100203040506.45678",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456780000, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMS5,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionMS4-NoOffset",
			DTValue:           "10100203040506.4567",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456700000, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMS4,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionMS3-NoOffset",
			DTValue:           "10100203040506.456",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 456000000, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMS3,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionMS2-NoOffset",
			DTValue:           "10100203040506.45",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 450000000, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMS2,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionMS1-NoOffset",
			DTValue:           "10100203040506.4",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 400000000, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMS1,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionSeconds-NoOffset",
			DTValue:           "10100203040506",
			Expected:          time.Date(1010, 2, 3, 4, 5, 6, 0, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionSeconds,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionMinutes-NoOffset",
			DTValue:           "101002030405",
			Expected:          time.Date(1010, 2, 3, 4, 5, 0, 0, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMinutes,
			HasOffset:         false,
		},
		{
			Name:              "PrecisionHours-NoOffset",
			DTValue:           "1010020304",
			Expected:          time.Date(1010, 2, 3, 4, 0, 0, 0, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionHours,
			HasOffset:         false,
		},

		// Full value, no offset, no hours
		{
			Name:              "PrecisionDay-NoOffset",
			DTValue:           "10100203",
			Expected:          time.Date(1010, 2, 3, 0, 0, 0, 0, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionDay,
			HasOffset:         false,
		},

		// Full value, no offset, no days
		{
			Name:              "PrecisionMonth-NoOffset",
			DTValue:           "101002",
			Expected:          time.Date(1010, 2, 1, 0, 0, 0, 0, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionMonth,
			HasOffset:         false,
		},

		// Full value, no offset, no month
		{
			Name:              "PrecisionYear-NoOffset",
			DTValue:           "1010",
			Expected:          time.Date(1010, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedPrecision: dcmtime.PrecisionYear,
			HasOffset:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			parsed, err := dcmtime.ParseDatetime(tc.DTValue)
			if err != nil {
				t.Fatal("parse err:", err)
			}

			if !tc.Expected.Equal(parsed.Time) {
				t.Errorf(
					"parsed time (%v) != expected (%v)",
					parsed.Time,
					tc.Expected,
				)

			}

			if parsed.Precision != tc.ExpectedPrecision {
				t.Errorf(
					"precision: expected %v, got %v",
					tc.ExpectedPrecision.String(),
					parsed.Precision.String(),
				)
			}
		})
	}
}

func TestParseDatetimeErr(t *testing.T) {
	testCases := []struct {
		Name     string
		BadValue string
	}{
		{
			Name:     "TotallyWrong",
			BadValue: "notaDT",
		},
		{
			Name:     "ContainsValidHead",
			BadValue: "10100203040506.456789+0102SomeText",
		},
		{
			Name:     "ContainsValidHead_LineBreak",
			BadValue: "10100203040506.456789+0102\nSomeText",
		},
		{
			Name:     "ContainsValidHead_WhiteSpace",
			BadValue: "10100203040506.456789+0102 SomeText",
		},
		{
			Name:     "ContainsValidTail",
			BadValue: "SomeText10100203040506.456789+0102",
		},
		{
			Name:     "ContainsValidTail_LineBreak",
			BadValue: "SomeText\n10100203040506.456789+0102",
		},
		{
			Name:     "ContainsValidTail_WhiteSpace",
			BadValue: "SomeText 10100203040506.456789+0102",
		},
		{
			Name:     "ExtraDigit_TZ",
			BadValue: "10100203040506.456789+01023",
		},
		{
			Name:     "MissingDigit_TZ",
			BadValue: "10100203040506.456789+010",
		},
		{
			Name:     "TZ_NoMinutes",
			BadValue: "10100203040506.456789+01",
		},
		{
			Name:     "TZ_SingleDigit",
			BadValue: "10100203040506.456789+1",
		},
		{
			Name:     "TZ_DoubleSign",
			BadValue: "10100203040506.456789++0102",
		},
		{
			Name:     "TZ_BadSign",
			BadValue: "10100203040506.456789&0102",
		},
		{
			Name:     "ExtraDigit_Milliseconds",
			BadValue: "10100203040506.4567891+0102",
		},
		{
			Name:     "ExtraDigit_Milliseconds_NoTZ",
			BadValue: "10100203040506.4567891",
		},
		{
			Name:     "ExtraDigit_Seconds",
			BadValue: "101002030405061.456789+0102",
		},
		{
			Name:     "ExtraDigit_Seconds_NoTZ",
			BadValue: "101002030405061.456789",
		},
		{
			Name:     "ExtraDigit_Seconds_NoMilliseconds",
			BadValue: "101002030405061",
		},
		{
			Name:     "MissingDigit_Seconds",
			BadValue: "1010020304056.456789+0102",
		},
		{
			Name:     "MissingDigit_Seconds_NoTZ",
			BadValue: "1010020304056.456789",
		},
		{
			Name:     "MissingDigit_Seconds_NoMilliseconds",
			BadValue: "1010020304056",
		},
		{
			Name:     "MissingDigit_Minutes",
			BadValue: "10100203045.456789+0102",
		},
		{
			Name:     "MissingDigit_Minutes_NoTZ",
			BadValue: "10100203045.456789",
		},
		{
			Name:     "MissingDigit_Minutes_NoMilliseconds",
			BadValue: "10100203045",
		},
		{
			Name:     "MissingDigit_Hours",
			BadValue: "101002034.456789+0102",
		},
		{
			Name:     "MissingDigit_Hours_NoTZ",
			BadValue: "101002034.456789",
		},
		{
			Name:     "MissingDigit_Hours_NoMilliseconds",
			BadValue: "101002034",
		},
		{
			Name:     "MissingDigit_Days",
			BadValue: "1010023.456789+0102",
		},
		{
			Name:     "MissingDigit_Days_NoTZ",
			BadValue: "1010023.456789",
		},
		{
			Name:     "MissingDigit_Days_NoMilliseconds",
			BadValue: "1010023",
		},
		{
			Name:     "MissingDigit_Months",
			BadValue: "10102.456789+0102",
		},
		{
			Name:     "MissingDigit_Months_NoTZ",
			BadValue: "10102.456789",
		},
		{
			Name:     "MissingDigit_Months_NoMilliseconds",
			BadValue: "10102.456789",
		},
		{
			Name:     "MissingDigit_Years",
			BadValue: "101.456789+0102",
		},
		{
			Name:     "MissingDigit_Years_NoTZ",
			BadValue: "101.456789",
		},
		{
			Name:     "MissingDigit_Years_NoMilliseconds",
			BadValue: "101.456789",
		},
		{
			Name:     "NoMillisecondsWithSeparator",
			BadValue: "10100203040506.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := dcmtime.ParseDatetime(tc.BadValue)
			if !errors.Is(err, dcmtime.ErrParseDT) {
				t.Errorf("expected ErrParseDT from ParseDatetime(), got %v", err)
			}
		})
	}
}

func TestDatetime_Methods(t *testing.T) {
	testCases := []struct {
		Name           string
		TimeVal        time.Time
		ExpectedDCM    string
		ExpectedString string
		Precision      dcmtime.PrecisionLevel
		NoOffset       bool
	}{
		{
			Name:           "PrecisionFull-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.456789-0102",
			ExpectedString: "1010-02-03 04:05:06.456789 -01:02",
			Precision:      dcmtime.PrecisionFull,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionFull-WithOffset-MSLeadingZeroes",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.000456-0102",
			ExpectedString: "1010-02-03 04:05:06.000456 -01:02",
			Precision:      dcmtime.PrecisionFull,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionFull-WithOffset-TruncateNanos",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789999, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.456789-0102",
			ExpectedString: "1010-02-03 04:05:06.456789 -01:02",
			Precision:      dcmtime.PrecisionFull,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionFull-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.456789",
			ExpectedString: "1010-02-03 04:05:06.456789",
			Precision:      dcmtime.PrecisionFull,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMS5-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.45678-0102",
			ExpectedString: "1010-02-03 04:05:06.45678 -01:02",
			Precision:      dcmtime.PrecisionMS5,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMS5-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.45678",
			ExpectedString: "1010-02-03 04:05:06.45678",
			Precision:      dcmtime.PrecisionMS5,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMS4-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.4567-0102",
			ExpectedString: "1010-02-03 04:05:06.4567 -01:02",
			Precision:      dcmtime.PrecisionMS4,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMS4-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.4567",
			ExpectedString: "1010-02-03 04:05:06.4567",
			Precision:      dcmtime.PrecisionMS4,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMS3-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.456-0102",
			ExpectedString: "1010-02-03 04:05:06.456 -01:02",
			Precision:      dcmtime.PrecisionMS3,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMS3-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.456",
			ExpectedString: "1010-02-03 04:05:06.456",
			Precision:      dcmtime.PrecisionMS3,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMS2-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.45-0102",
			ExpectedString: "1010-02-03 04:05:06.45 -01:02",
			Precision:      dcmtime.PrecisionMS2,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMS2-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.45",
			ExpectedString: "1010-02-03 04:05:06.45",
			Precision:      dcmtime.PrecisionMS2,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMS1-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.4-0102",
			ExpectedString: "1010-02-03 04:05:06.4 -01:02",
			Precision:      dcmtime.PrecisionMS1,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMS1-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506.4",
			ExpectedString: "1010-02-03 04:05:06.4",
			Precision:      dcmtime.PrecisionMS1,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionSeconds-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506-0102",
			ExpectedString: "1010-02-03 04:05:06 -01:02",
			Precision:      dcmtime.PrecisionSeconds,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionSeconds-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203040506",
			ExpectedString: "1010-02-03 04:05:06",
			Precision:      dcmtime.PrecisionSeconds,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMinutes-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "101002030405-0102",
			ExpectedString: "1010-02-03 04:05 -01:02",
			Precision:      dcmtime.PrecisionMinutes,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMinutes-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "101002030405",
			ExpectedString: "1010-02-03 04:05",
			Precision:      dcmtime.PrecisionMinutes,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionHours-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "1010020304-0102",
			ExpectedString: "1010-02-03 04 -01:02",
			Precision:      dcmtime.PrecisionHours,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionHours-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "1010020304",
			ExpectedString: "1010-02-03 04",
			Precision:      dcmtime.PrecisionHours,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionDay-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203-0102",
			ExpectedString: "1010-02-03 -01:02",
			Precision:      dcmtime.PrecisionDay,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionDay-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "10100203",
			ExpectedString: "1010-02-03",
			Precision:      dcmtime.PrecisionDay,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionMonth-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "101002-0102",
			ExpectedString: "1010-02 -01:02",
			Precision:      dcmtime.PrecisionMonth,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionMonth-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "101002",
			ExpectedString: "1010-02",
			Precision:      dcmtime.PrecisionMonth,
			NoOffset:       true,
		},
		{
			Name:           "PrecisionYear-WithOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "1010-0102",
			ExpectedString: "1010 -01:02",
			Precision:      dcmtime.PrecisionYear,
			NoOffset:       false,
		},
		{
			Name:           "PrecisionYear-NoOffset",
			TimeVal:        time.Date(1010, 2, 3, 4, 5, 6, 456789000, time.FixedZone("", -3720)),
			ExpectedDCM:    "1010",
			ExpectedString: "1010",
			Precision:      dcmtime.PrecisionYear,
			NoOffset:       true,
		},
	}

	for _, tc := range testCases {
		dt := dcmtime.Datetime{
			Time:      tc.TimeVal,
			Precision: tc.Precision,
			NoOffset:  tc.NoOffset,
		}

		// Run one master sub-test for each case
		t.Run(tc.Name, func(t *testing.T) {

			// Break out methods into sub-sub tests
			t.Run("DCM()", func(t *testing.T) {
				dcmVal := dt.DCM()
				if dcmVal != tc.ExpectedDCM {
					t.Errorf("DCM(): expected '%v', got '%v'", tc.ExpectedDCM, dcmVal)
				}
			})

			t.Run("String()", func(t *testing.T) {
				strVal := dt.String()
				if strVal != tc.ExpectedString {
					t.Errorf(
						"String(): expected '%v', got '%v'", tc.ExpectedString, strVal,
					)
				}
			})
		})
	}
}
