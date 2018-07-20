package plural

import "strings"

// PluralSpec defines the CLDR plural rules for a language.
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type PluralSpec struct {
	PluralsSlice []Plural
	Plurals      map[Plural]struct{}
	PluralFunc   func(*Operands) Plural
}

func NewPluralSpec(plurals []Plural, pluralFunc func(*Operands) Plural) *PluralSpec {
	return &PluralSpec{
		PluralsSlice: plurals,
		Plurals:      newPluralSet(plurals...),
		PluralFunc:   pluralFunc,
	}
}

func normalizePluralSpecID(id string) string {
	id = strings.Replace(id, "_", "-", -1)
	id = strings.ToLower(id)
	return id
}

// Plural returns the plural category for number as defined by
// the language's CLDR plural rules.
func (ps *PluralSpec) Plural(number interface{}) (Plural, error) {
	ops, err := newOperands(number)
	if err != nil {
		return Invalid, err
	}
	return ps.PluralFunc(ops), nil
}

// Create Plural Map from Slice
func (ps *PluralSpec) CreatePluralMapFromSlice(s []string) map[Plural]string {
	m := stringSliceToMap(s)
	fallback := ""
	pluralMap := map[Plural]string{}
	for k, plural := range ps.PluralsSlice {
		vv, ok := m[k]
		if !ok {
			pluralMap[plural] = fallback
			continue
		}
		fallback = vv
		pluralMap[plural] = vv
	}
	return pluralMap
}

func newPluralSet(plurals ...Plural) map[Plural]struct{} {
	set := make(map[Plural]struct{}, len(plurals))
	for _, plural := range plurals {
		set[plural] = struct{}{}
	}
	return set
}

func IntInRange(i, from, to int64) bool {
	return from <= i && i <= to
}

func IntEqualsAny(i int64, any ...int64) bool {
	for _, a := range any {
		if i == a {
			return true
		}
	}
	return false
}

func stringSliceToMap(s []string) map[int]string {
	m := map[int]string{}
	for k, v := range s {
		m[k] = v
	}
	return m
}
