package web

import (
	"github.com/BurntSushi/toml"
	"github.com/markbates/pkger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"regexp"
)

func compileLanguages() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Files muse be listed with Include for pkger to pull them in
	files := map[string]string{
		"active.es.toml": pkger.Include("/active.es.toml"),
	}

	for filename, file := range files {
		langFile, err := pkger.Open(file)
		if err != nil {
			return nil, err
		}
		defer langFile.Close()

		fileinfo, err := langFile.Stat()
		if err != nil {
			return nil, err
		}

		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)

		_, err = langFile.Read(buffer)
		if err != nil {
			return nil, err
		}

		bundle.MustParseMessageFileBytes(buffer, filename)
	}

	return bundle, nil
}

func isValidUUID4(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func roundUp(f float64) uint64 {
	fint := uint64(f)

	if f > float64(fint) {
		fint++
	}

	return fint
}
