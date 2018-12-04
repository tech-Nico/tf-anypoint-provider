package sdk

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	tabwriterMinWidth = 10
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = 0
)

func tabularize(elems []string) string {
	toReturn := ""

	for _, val := range elems {
		toReturn = toReturn + strings.TrimSpace(val) + "\t"
	}

	return toReturn
}

func PrintTabular(headers []string, data [][]string) {

	w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
	defer w.Flush()
	fmt.Println("")
	headersStr := tabularize(headers)
	fmt.Fprintln(w, headersStr)

	for _, row := range data {
		lineStr := tabularize(row)
		fmt.Fprintln(w, lineStr)
	}

}

func Debug(doSomething func()) {
	if debug := os.Getenv("ANYPOINT_DEBUG"); debug != "" && strings.EqualFold(debug, "true") {
		doSomething()
	}
}

func PrintAsList(objects []interface{}, extractTabularDataFunc func([]interface{}) [][]string, headers []string) {

	data := extractTabularDataFunc(objects)

	PrintTabular(headers, data)
}

func PrintAsJSON(objects []interface{}) {
	b, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		fmt.Println("Error while marshalling output:", err)
	}
	os.Stdout.Write(b)
}

func OpenYAMLFile(f string, t interface{}) error {

	fileContent, err := ioutil.ReadFile(f)

	if err != nil {
		return fmt.Errorf("Error while opening yaml file %q. Error: %s", f, err)
	}

	err = yaml.Unmarshal(fileContent, t)

	if err != nil {
		return fmt.Errorf("Error while parsing YAML file %q . Error: %s", f, err)
	}

	return nil
}
