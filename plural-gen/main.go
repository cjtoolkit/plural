package main

import (
	"encoding/xml"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"

	toml "github.com/pelletier/go-toml"
)

func main() {
	if len(os.Args) <= 2 {
		return
	}

	var data SupplementalData
	if err := xml.Unmarshal([]byte(pluralsXml), &data); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	var ordinalData SupplementalData
	if err := xml.Unmarshal([]byte(ordinalsXml), &ordinalData); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	localeMap := map[string]PluralGroup{}
	for _, pg := range data.PluralGroups {
		for _, locale := range pg.SplitLocales() {
			localeMap[locale] = pg
		}
	}

	localeOrdinalMap := map[string]PluralGroup{}
	for _, pg := range ordinalData.PluralGroups {
		for _, locale := range pg.SplitLocales() {
			localeOrdinalMap[locale] = pg
		}
	}

	packageName := os.Args[1]
	userFileName := os.Args[2]
	file, err := os.Open(userFileName)
	if err != nil {
		log.Fatalf("Unable to open user file: %s", err)
	}

	var user UserSupplement
	err = toml.NewDecoder(file).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to parse user file with toml: %s", err)
	}
	file.Close()

	glue := []Glue{}
	for _, locale := range user.Locales {
		group, ok := localeMap[locale.Code]
		if !ok {
			log.Fatalf("User Locale not found: %s", locale.Code)
		}

		groupOridinal := localeOrdinalMap[locale.Code]

		glue = append(glue, Glue{
			Locale:             locale,
			PluralGroup:        group,
			PluralGroupOrdinal: groupOridinal,
		})
	}

	ctx := Context{
		PackageName: packageName,
		Source:      userFileName,
		Glue:        glue,
	}

	userFileNameSplit := strings.Split(path.Base(userFileName), ".")
	codeName := strings.Join(userFileNameSplit[:len(userFileNameSplit)-1], ".") + ".go"
	testName := strings.Join(userFileNameSplit[:len(userFileNameSplit)-1], ".") + "_test.go"

	relationRegexp := regexp.MustCompile("([niftvw])(?: % ([0-9]+))? (!=|=)(.*)")
	funcs := template.FuncMap{
		"relationRegexp": func() *regexp.Regexp { return relationRegexp },
	}

	file, err = os.Create(codeName)
	if err != nil {
		log.Fatalf("Unable to create file: %s", err)
	}

	template.Must(template.New("code").Funcs(funcs).Parse(codeTemplate)).Execute(file, ctx)
	file.Close()
	exec.Command("go", "fmt", codeName).Run()

	file, err = os.Create(testName)
	if err != nil {
		log.Fatalf("Unable to create file: %s", err)
	}

	template.Must(template.New("code").Parse(testTemplate)).Execute(file, ctx)
	file.Close()
	exec.Command("go", "fmt", testName).Run()
}
