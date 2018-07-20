package main

import "strings"

type UserSupplement struct {
	Locales []Locale `toml:"locales"`
}

type Locale struct {
	Code         string `toml:"code"`
	FunctionName string `toml:"functionName"`
}

func (l Locale) FunctionNameTitle() string { return strings.Title(l.FunctionName) }
