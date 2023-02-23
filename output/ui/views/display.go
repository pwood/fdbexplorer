package views

import (
	"fmt"
	"strings"
)

const Kibibyte float64 = 1024
const Mebibyte float64 = Kibibyte * 1024
const Gibibyte float64 = Mebibyte * 1024
const Tebibyte float64 = Gibibyte * 1024

func Titlify(s string) string {
	return strings.Title(strings.Replace(s, "_", " ", -1))
}

func Boolify(b bool) string {
	if b {
		return "Yes"
	} else {
		return "No"
	}
}

const None string = ""

func Convert(v float64, dp int, perUnit string) string {
	f := fmt.Sprintf("%%0.%df %%s%%s", dp)
	suffix := ""

	if perUnit != None {
		suffix = fmt.Sprintf("/%s", perUnit)
	}

	if v < Kibibyte {
		return fmt.Sprintf(f, v, "B", suffix)
	} else if v < Mebibyte {
		return fmt.Sprintf(f, v/Kibibyte, "KiB", suffix)
	} else if v < Gibibyte {
		return fmt.Sprintf(f, v/Mebibyte, "MiB", suffix)
	} else if v < Tebibyte {
		return fmt.Sprintf(f, v/Gibibyte, "GiB", suffix)
	} else {
		return fmt.Sprintf(f, v/Tebibyte, "TiB", suffix)
	}
}
