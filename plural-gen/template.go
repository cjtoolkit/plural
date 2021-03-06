package main

type Context struct {
	PackageName string
	Source      string
	Glue        []Glue
}

type Glue struct {
	Locale             Locale
	PluralGroup        PluralGroup
	PluralGroupOrdinal PluralGroup
}

const codeTemplate = `// Code generated by plural-gen. DO NOT EDIT.
// Source: {{ .Source }}

package {{ .PackageName }}

import (
	"github.com/cjtoolkit/plural"
)

{{ $relationRegexp := relationRegexp }}{{ range .Glue }}
// {{ .Locale.Code }}
func {{ .Locale.FunctionName }}() plural.PluralGroup {
	return plural.PluralGroup{
		{{ with .PluralGroup }} Cardinal: plural.NewPluralSpec([]plural.Plural{  {{ range $i, $e := .PluralRules }}{{ if $i }}, {{ end }}plural.{{ $e.CountTitle }}{{ end }} }, func(ops *plural.Operands) plural.Plural { {{ range .PluralRules }}{{ if .GoCondition $relationRegexp }}
			// {{ .Condition }}
			if {{ .GoCondition $relationRegexp }} {
				return plural.{{ .CountTitle }}
			}{{ end }}{{ end }}
			return plural.Other
		}), {{ end }}
		{{ with .PluralGroupOrdinal }} Ordinal: plural.NewPluralSpec([]plural.Plural{  {{ range $i, $e := .PluralRules }}{{ if $i }}, {{ end }}plural.{{ $e.CountTitle }}{{ end }} }, func(ops *plural.Operands) plural.Plural { {{ range .PluralRules }}{{ if .GoCondition $relationRegexp }}
			// {{ .Condition }}
			if {{ .GoCondition $relationRegexp }} {
				return plural.{{ .CountTitle }}
			}{{ end }}{{ end }}
			return plural.Other
		}), {{ end }}
	}
}

{{ end }}
`

const testTemplate = `// Code generated by plural-gen. DO NOT EDIT.
// Source: {{ .Source }}

package {{ .PackageName }}

import (
	"testing"

	"github.com/cjtoolkit/plural"
	"github.com/cjtoolkit/plural/pluralTestUtil"
)

{{ range .Glue }}{{ $locale := .Locale }}
// Test {{ $locale.Code }}
func Test{{ $locale.FunctionNameTitle }}(t *testing.T) {
	group := {{ $locale.FunctionName }}()
	{{ with .PluralGroup }} t.Run("Cardinal", func(t *testing.T) {
		var tests []pluralTestUtil.PluralTest
		{{ range .PluralRules }}
		{{ if .IntegerExamples }}tests = pluralTestUtil.AppendIntegerTests(tests, plural.{{ .CountTitle }}, {{ printf "%#v" .IntegerExamples }}){{ end }}
		{{ if .DecimalExamples }}tests = pluralTestUtil.AppendDecimalTests(tests, plural.{{ .CountTitle }}, {{ printf "%#v" .DecimalExamples }}){{ end }}
		{{ end }}
		pluralTestUtil.Run(t, "{{ $locale.FunctionName }} (Cardinal)", group.Cardinal, tests)
	}){{ end }}

	{{ with .PluralGroupOrdinal }} t.Run("Ordinal", func(t *testing.T) {
		var tests []pluralTestUtil.PluralTest
		{{ range .PluralRules }}
		{{ if .IntegerExamples }}tests = pluralTestUtil.AppendIntegerTests(tests, plural.{{ .CountTitle }}, {{ printf "%#v" .IntegerExamples }}){{ end }}
		{{ if .DecimalExamples }}tests = pluralTestUtil.AppendDecimalTests(tests, plural.{{ .CountTitle }}, {{ printf "%#v" .DecimalExamples }}){{ end }}
		{{ end }}
		pluralTestUtil.Run(t, "{{ $locale.FunctionName }} (Ordinal)", group.Ordinal, tests)
	}){{ end }}
}

{{ end }}
`
