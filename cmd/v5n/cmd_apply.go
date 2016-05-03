package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"limbo.services/version"
)

func apply(app *kingpin.Application, binPath, versionStr, authorStr, commitStr, dateStr string) {
	payloadPlaceholder := bytes.ToLower(payloadPlaceholderUpper)

	data, err := ioutil.ReadFile(binPath)
	app.FatalIfError(err, "v5n")

	if !bytes.Contains(data, payloadPlaceholder) {
		err := errors.New("placeholder not found. please import \"limbo.services/version\"")
		app.FatalIfError(err, "v5n")
	}

	vsn, err := version.ParseSemver(versionStr)
	app.FatalIfError(err, "v5n")

	if dateStr != "" {
		dateStr = strings.TrimSpace(dateStr)
		var formats = []string{
			time.UnixDate,
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02 15",
			"2006-01-02",
		}
		for _, f := range formats {
			d, dateErr := time.Parse(f, dateStr)
			if dateErr == nil {
				vsn.Released = d.UTC()
				break
			}
		}
		if vsn.Released.IsZero() {
			err = fmt.Errorf("invalid release date format (%q)", dateStr)
			app.FatalIfError(err, "v5n")
		}
	}

	if vsn.Released.IsZero() {
		vsn.Released = time.Now().UTC()
	}

	vsn.ReleasedBy = authorStr
	vsn.Commit = commitStr

	payload, err := json.Marshal(vsn)
	app.FatalIfError(err, "v5n")

	if len(payload) > len(payloadPlaceholder) {
		err := fmt.Errorf("version info is to large. must be smaller than %d bytes", len(payloadPlaceholder))
		app.FatalIfError(err, "v5n")
	}

	payload = append(payload, bytes.Repeat([]byte(" "), len(payloadPlaceholder)-len(payload))...)
	if len(payload) != len(payloadPlaceholder) {
		panic("should not happen")
	}

	data = bytes.Replace(data, payloadPlaceholder, payload, -1)

	err = ioutil.WriteFile(binPath, data, 0755)
	app.FatalIfError(err, "v5n")
}

var payloadPlaceholderUpper = []byte(strings.Join([]string{
	"E927CC4876808AB86054E3489A04EFD20BC9CF9F3FE2356E56B1274AA8FF4FC0DFA8F972031",
	"53FB75A3E6F274C84094B5B20A0306943E121CE818B5AF8333C9EBAF084ABF27F78EFFAF7EA",
	"1C36ED89BDF8FF8A369DA01388206D987A52ED22CB29FA600D61DA0772C5822499337BC8AD8",
	"655EBE185BFFF5C4EABA1D4DE5A577863ED661607379003B94374DD85B0C35E24DCFC3DABB0",
	"147607582C7402A782BE5FC0A19B7A92FB0C91599ED5DFE25CF180BC675CEF87CCD1F79BA86",
	"C72768AB7862831BFAAC0DA54BCA6166C78DBC558E2E324E5F85EE22156901B0E82C2AC9D2A",
	"CF29ED11AE86852A57C3C53EF75D292D0C5D21CF1DEB7E5FB1BD641FB46A97718F7983260B9",
	"415EB0B240731DE0359BF1E3764954CFB94277BDF972B13EAA6A38C0E3BB0BE58FA850857B7",
	"74F325E336DCF2550644ECAA1EDCDEA3B44E7632C5AE7D723D4E8692C04E1D1A9FA64CDCE23",
	"082A9032F34027A994AC6A13DDC7B9D3204350908FE1567BFF31ECE702051446E42E8EA7DFE",
	"D88F88AD42DE1B0EFFB19CCD1DA94462D17411EDC6FB510175912EA455BC7387E2040CBF0CD",
	"79EF4BAAE27375FAE38F5351F5CF4EBD540D7D560EB7CBE8D6AA2E040FAA0A2C00F8A32759F",
	"5A1BAFEE6BA690192B64BEE612DFE6E142D3FF53854ADB91E2DA8F86F0A58685D12B832E083",
	"BAF61AD80F6A353224E16BC7693585E2894147B8286985032",
}, ""))
