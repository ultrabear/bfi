package compiler

import (
	"regexp"
	"strings"
)

func countoccurances(indata *string) int {
	outdata := 0
	outdata += strings.Count(*indata, "><") + strings.Count(*indata, "<>")
	outdata += strings.Count(*indata, "+-") + strings.Count(*indata, "-+")
	outdata += strings.Count(*indata, "[]")
	outdata += strings.Count(*indata, "[->-[-]<]")
	return outdata
}

func stripoccurances(indata *string) string {
	outdata := strings.Replace(strings.Replace(*indata, "><", "", -1), "<>", "", -1)
	outdata = strings.Replace(strings.Replace(outdata, "+-", "", -1), "-+", "", -1)
	outdata = strings.Replace(outdata, "[]", "", -1)
	outdata = strings.Replace(outdata, "[->-[-]<]", "[-]>[-]<", -1)
	return outdata
}

func filterbfc(indata *string) string {
	gex, _ := regexp.Compile("[\\[\\]\\-\\+\\>\\<\\,\\.]+")
	data := gex.FindAllStringSubmatch(*indata, -1)
	joindata := make([]string, 0)
	for _, item := range data {
		joindata = append(joindata, strings.Join(item, ""))
	}
	return strings.Join(joindata, "")
}

func CompressBFC(indata string) string {
	export := filterbfc(&indata)
	for countoccurances(&export) != 0 {
		export = stripoccurances(&export)
	}
	return export
}
