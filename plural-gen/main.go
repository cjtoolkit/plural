package main

import (
	"encoding/xml"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
)

func main() {
	if len(os.Args) <= 2 {
		return
	}

	var data SupplementalData
	if err := xml.Unmarshal([]byte(pluralsXml), &data); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	localeMap := map[string]PluralGroup{}
	for _, pg := range data.PluralGroups {
		for _, locale := range pg.SplitLocales() {
			localeMap[locale] = pg
		}
	}

	packageName := os.Args[1]
	userFileName := os.Args[2]
	file, err := os.Open(userFileName)
	if err != nil {
		log.Fatalf("Unable to open user file: %s", err)
	}

	var user UserSupplement
	_, err = toml.DecodeReader(file, &user)
	if err != nil {
		log.Fatalf("Unable to parse user file with toml: %s", err)
	}
	file.Close()

	glue := []Glue{}
	for _, locale := range user.Locales {
		group, ok := localeMap[locale.Locale]
		if !ok {
			log.Fatalf("User Locale not found: %s", locale.Locale)
		}

		glue = append(glue, Glue{
			Locale:      locale,
			PluralGroup: group,
		})
	}

	ctx := Context{
		PackageName: packageName,
		Source:      userFileName,
		Glue:        glue,
	}

	userFileNameSplit := strings.Split(userFileName, ".")
	codeName := strings.Join(userFileNameSplit[:len(userFileNameSplit)-1], ".") + ".go"
	testName := strings.Join(userFileNameSplit[:len(userFileNameSplit)-1], ".") + "_test.go"

	file, err = os.Create(codeName)
	if err != nil {
		log.Fatalf("Unable to create file: %s", err)
	}

	template.Must(template.New("code").Parse(codeTemplate)).Execute(file, ctx)
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
