package compiler

import (
	"regexp"
	"runtime"
	"strings"
)

func countoccurances(indata string) bool {
	if strings.Contains(indata, "><") || strings.Contains(indata, "<>") || strings.Contains(indata, "+-") || strings.Contains(indata, "-+") || strings.Contains(indata, "[]") || strings.Contains(indata, "[->-[-]<]") {
		return true
	} else {
		return false
	}
}

func stripoccurances(indata *string) string {
	outdata := strings.Replace(strings.Replace(*indata, "><", "", -1), "<>", "", -1)
	outdata = strings.Replace(strings.Replace(outdata, "+-", "", -1), "-+", "", -1)
	outdata = strings.Replace(outdata, "[]", "", -1)
	outdata = strings.Replace(outdata, "[->-[-]<]", "[-]>[-]<", -1)
	return outdata
}

func filterchunk(s string, filter *regexp.Regexp, rch chan string) {
	data := filter.FindAllStringSubmatch(s, -1)
	joindata := make([]string, len(data))
	for i, item := range data {
		joindata[i] = strings.Join(item, "")
	}
	rch <- strings.Join(joindata, "")
}

func filterbfc(indata *string) string {
	gex, err := regexp.Compile("[\\[\\]\\-\\+\\>\\<\\,\\.]+")
	if err != nil {
		panic(err)
	}

	returns := make([]chan string, runtime.NumCPU())
	splitwidth := (len(*indata) + (len(returns) - 1)) / len(returns)
	for i := 0; i < len(returns); i++ {
		returns[i] = make(chan string)
		nextwidth := (splitwidth * i) + splitwidth
		if nextwidth > len(*indata) {
			nextwidth = len(*indata)
		}
		go filterchunk((*indata)[splitwidth*i:nextwidth], gex, returns[i])
	}

	joindata := make([]string, len(returns))
	for i, item := range returns {
		joindata[i] = <-item
	}
	return strings.Join(joindata, "")
}

func CompressBFC(indata string) string {
	export := filterbfc(&indata)
	for countoccurances(export) {
		export = stripoccurances(&export)
	}
	return export
}
