package main

type UserSupplement struct {
	Locales []Locale `toml:"locales"`
}

type Locale struct {
	Locale       string `toml:"locale"`
	FunctionName string `toml:"functionName"`
}
