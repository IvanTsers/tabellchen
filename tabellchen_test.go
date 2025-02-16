package tabellchen

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestReadTable(t *testing.T) {
	// Expected tables
	h := []string{"c1", "c2", "c3"}
	f := [][]string{{"f1", "f2", "f3"}, {"f4", "f5", "f6"}}
	h6fields := NewTable(h, f)
	noh6fields := NewTable([]string{}, f)
	var errProper error = nil
	var errBadWidth error = fmt.Errorf("Line 2 has 4 fields," +
		" expected 3\n")

	// Expected output slices
	proper := []interface{}{h6fields, errProper}
	noheader := []interface{}{noh6fields, errProper}
	badWidth := []interface{}{Table{}, errBadWidth}
	testCases := []struct {
		name  string
		sep   rune
		input string
		want  []interface{}
	}{
		{"properTsv", '\t', "data/proper.tsv", proper},
		{"properCsv", ',', "data/proper.csv", proper},
		{"noheader", '\t', "data/noheader.tsv", noheader},
		{"badWidth", '\t', "data/badWidth.tsv", badWidth},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hasHeader := tc.name != "noheader"

			config := ReadConfig{
				FilePath:    tc.input,
				Separator:   tc.sep,
				CommentChar: '#',
				HasHeader:   hasHeader,
			}

			tab, err := ReadTable(config)
			get := []interface{}{tab, err}
			if !reflect.DeepEqual(get, tc.want) {
				t.Errorf("want:\n%v\nget:\n%v\n",
					tc.want, get)
			}
		})
	}
}
func isBob(s string) bool {
	return s == "Bob"
}

func TestFilter(t *testing.T) {
	h := []string{"Name", "Year", "Color"}
	f := [][]string{
		{"Bob", "2022", "red"},
		{"Bob", "2024", "yellow"},
	}
	want := NewTable(h, f)

	config := ReadConfig{
		FilePath:    "data/bob.csv",
		Separator:   ',',
		CommentChar: '#',
		HasHeader:   true,
	}

	tab, _ := ReadTable(config)
	get, _ := Filter(tab, "Name", isBob)
	if !reflect.DeepEqual(get, want) {
		t.Errorf("want:\n%v\nget:\n%v\n",
			want, get)
	}
}
func TestWriteTable(t *testing.T) {
	header := []string{"c1", "c2", "c3"}
	rows := [][]string{
		{"f1", "f2", "f3"},
		{"f4", "f5", "f6"},
	}

	tab := NewTable(header, rows)
	outputFile := "data/output.csv"
	file, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	wconfig := WriteConfig{
		File:      file,
		Separator: ',',
	}

	err = tab.WriteTable(wconfig)

	if err != nil {
		t.Fatalf("Failed to write table: %v", err)
	}
	rconfig := ReadConfig{
		FilePath:    "data/output.csv",
		Separator:   ',',
		CommentChar: '#',
		HasHeader:   true,
	}

	writtenTab, err := ReadTable(rconfig)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	tabsAreEqual := reflect.DeepEqual(tab, writtenTab)
	if !tabsAreEqual {
		t.Errorf("Written content mismatch: "+
			"\nwant:\n%s\nget:\n%s\n",
			tab, writtenTab)
	}
	err = os.Remove(outputFile)
	if err != nil {
		t.Fatalf("Failed to remove output file: %v",
			err)
	}
}
func TestNewColumn(t *testing.T) {
	header := []string{"c1", "c2", "c3"}
	rows := [][]string{
		{"f1", "f2", "f3"},
		{"f4", "f5", "f6"},
	}

	tab := NewTable(header, rows)
	oldNumCol := len(tab.Header)
	tab.NewColumn("c4")
	newNumCol := len(tab.Header)
	if newNumCol != oldNumCol+1 {
		t.Errorf("expected %d columns, got %d\n",
			oldNumCol+1, newNumCol)
	}
	if tab.Header[3] != "c4" {
		t.Errorf("expected %s as the last columns, got %v\n",
			"c4", tab.Header[3])
	}
}
func TestReorderColumns(t *testing.T) {
	header := []string{"c1", "c2", "c3"}
	rows := [][]string{
		{"f1", "f2", "f3"},
		{"f4", "f5", "f6"},
	}

	tab := NewTable(header, rows)
	tab.ReorderColumns(1, 0, 2)

	//Check the headers' order
	hwant := []string{"c2", "c1", "c3"}
	hget := tab.Header
	if !reflect.DeepEqual(hwant, hget) {
		t.Errorf("wrong new column order:\n"+
			"want: %v\nget: %v\n",
			hwant, hget)
	}

	//Check the columns' order
	rwant := []string{"f5", "f4", "f6"}
	rget := tab.Rows[1]
	if !reflect.DeepEqual(rwant, rget) {
		t.Errorf("wrong new row order:\n"+
			"want: %v\nget: %v\n",
			rwant, rget)
	}
}
