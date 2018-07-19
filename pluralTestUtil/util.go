package pluralTestUtil

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cjtoolkit/plural"
)

type PluralTest struct {
	num    interface{}
	plural plural.Plural
}

func AppendIntegerTests(tests []PluralTest, plural plural.Plural, examples []string) []PluralTest {
	for _, ex := range expandExamples(examples) {
		i, err := strconv.ParseInt(ex, 10, 64)
		if err != nil {
			panic(err)
		}
		tests = append(tests, PluralTest{ex, plural}, PluralTest{i, plural})
	}
	return tests
}

func AppendDecimalTests(tests []PluralTest, plural plural.Plural, examples []string) []PluralTest {
	for _, ex := range expandExamples(examples) {
		tests = append(tests, PluralTest{ex, plural})
	}
	return tests
}

func expandExamples(examples []string) []string {
	var expanded []string
	for _, ex := range examples {
		if parts := strings.Split(ex, "~"); len(parts) == 2 {
			for ex := parts[0]; ; ex = increment(ex) {
				expanded = append(expanded, ex)
				if ex == parts[1] {
					break
				}
			}
		} else {
			expanded = append(expanded, ex)
		}
	}
	return expanded
}
func increment(dec string) string {
	runes := []rune(dec)
	carry := true
	for i := len(runes) - 1; carry && i >= 0; i-- {
		switch runes[i] {
		case '.':
			continue
		case '9':
			runes[i] = '0'
		default:
			runes[i]++
			carry = false
		}
	}
	if carry {
		runes = append([]rune{'1'}, runes...)
	}
	return string(runes)
}

func Run(t *testing.T, name string, spec *plural.PluralSpec, tests []PluralTest) {
	for _, test := range tests {
		if plural2, err := spec.Plural(test.num); plural2 != test.plural {
			t.Errorf("%s: PluralCategory(%#v) returned %s, %v; expected %s", name, test.num, plural2, err, test.plural)
		}
	}
}
