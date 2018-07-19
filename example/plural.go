// Code generated by plural-gen. DO NOT EDIT.
// Source: plural.toml

package example

import (
	"github.com/cjtoolkit/plural"
)

// en
func English() *plural.PluralSpec {
	fn := func(ops *plural.Operands) plural.Plural {
		// i = 1 and v = 0
		if plural.IntEqualsAny(ops.I, 1) && plural.IntEqualsAny(ops.V, 0) {
			return plural.One
		}
		return plural.Other
	}
	return plural.NewPluralSpec([]plural.Plural{plural.One, plural.Other}, fn)
}

// fr
func French() *plural.PluralSpec {
	fn := func(ops *plural.Operands) plural.Plural {
		// i = 0,1
		if plural.IntEqualsAny(ops.I, 0, 1) {
			return plural.One
		}
		return plural.Other
	}
	return plural.NewPluralSpec([]plural.Plural{plural.One, plural.Other}, fn)
}
