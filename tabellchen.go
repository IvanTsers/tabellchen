package tabellchen

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ReadConfig struct {
	FilePath    string
	Separator   rune
	CommentChar rune
	HasHeader   bool
}
type WriteConfig struct {
	File      *os.File
	Separator rune
}
type Table struct {
	Header []string
	Rows   [][]string
}

// The method NewTable creates a table with specified header and rows.
func NewTable(header []string, rows [][]string) Table {
	t := Table{
		Header: header,
		Rows:   rows,
	}
	return t
}
func (t Table) ColIdByName(colname string) (int, error) {
	id := -1
	for i, h := range t.Header {
		if h == colname {
			id = i
			break
		}
	}
	var err error = nil
	if id == -1 {
		err = fmt.Errorf("ColIdByName: "+
			"column '%s' not found\n", colname)
	}
	return id, err
}

// The method WriteTable writes a Table into an os.File.
func (t Table) WriteTable(config WriteConfig) error {
	file := config.File
	separator := config.Separator

	var writeErr error
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	if len(t.Header) > 0 {
		stringHeader := strings.Join(t.Header, string(separator))
		_, err := writer.WriteString(stringHeader + "\n")
		if err != nil {
			writeErr := fmt.Errorf("WriteTable: failed "+
				"to write header: %v\n", err)
			return writeErr
		}
	}
	rn := 0
	for _, row := range t.Rows {
		stringRow := strings.Join(row, string(separator))
		_, err := writer.WriteString(stringRow + "\n")
		if err != nil {
			writeErr := fmt.Errorf("WriteTable: failed "+
				"to write row %d: %v\n", rn, err)
			return writeErr
		}
		rn += 1
	}
	return writeErr
}

// The method NewColumn creates a new empty named column in a Table.
func (t *Table) NewColumn(colName string) {
	t.Header = append(t.Header, colName)

	for i := range t.Rows {
		t.Rows[i] = append(t.Rows[i], "")
	}
}

// The method ReorderColumns changes the order of table's columns according to the order of specified indices.
func (t *Table) ReorderColumns(newOrder ...int) error {
	fmt.Println(newOrder)
	fmt.Println(t)
	numCols := len(t.Header)
	if len(newOrder) != numCols {
		return fmt.Errorf("ReorderColumns: "+
			"expected %d indices, got %d",
			numCols, len(newOrder))
	}
	seen := make(map[int]bool)
	for _, index := range newOrder {
		if index < 0 || index >= numCols {
			return fmt.Errorf("ReorderColumns: "+
				"Index out of range: %d", index)
		}
		if seen[index] {
			return fmt.Errorf("ReorderColumns: "+
				"Duplicate index: %d", index)
		}
		seen[index] = true
	}
	newHeader := make([]string, numCols)
	for i, newIndex := range newOrder {
		newHeader[i] = t.Header[newIndex]
	}
	t.Header = newHeader
	for i := range t.Rows {
		newRow := make([]string, numCols)
		for j, newIndex := range newOrder {
			newRow[j] = t.Rows[i][newIndex]
		}
		t.Rows[i] = newRow
	}
	return nil
}

// The function ReadTable reads data from a file and populates a Table. It accepts a ReadConfig struct.
func ReadTable(config ReadConfig) (Table, error) {
	path := config.FilePath
	separator := config.Separator
	commentChar := config.CommentChar
	hasHeader := config.HasHeader

	file, err := os.Open(path)
	if err != nil {
		return Table{},
			fmt.Errorf("Failed to open file: %v\n", err)
	}
	defer file.Close()
	table := Table{}
	numFields := 0
	currentLine := 0
	scanner := bufio.NewScanner(file)
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()
		currentLine += 1
		if len(line) == 0 {
			continue
		}
		if rune(line[0]) == commentChar {
			continue
		}
		fields := strings.FieldsFunc(line, func(c rune) bool {
			return c == separator
		})
		if firstLine {
			numFields = len(fields)
			if hasHeader {
				table.Header = fields
			} else {
				table.Header = []string{}
				table.Rows = append(table.Rows, fields)

			}
			firstLine = false
			continue
		}
		if len(fields) != numFields {
			return Table{},
				fmt.Errorf("Line %d has %d fields, expected %d\n",
					currentLine,
					len(fields),
					numFields)
		}
		table.Rows = append(table.Rows, fields)
	}

	if err := scanner.Err(); err != nil {
		return Table{},
			fmt.Errorf("Error reading file: %v\n", err)
	}

	return table, nil
}

// The function Filter returns rows of a Table that satisfy a condition.
func Filter(t Table,
	col interface{},
	cond func(string) bool) (Table, error) {
	colId := -1
	switch v := col.(type) {
	case int:
		colId = v
		nc := len(t.Rows[0])
		if colId > nc {
			return t,
				fmt.Errorf("Filter: tried to access "+
					"column %d, but there are "+
					"only %d columns\n", colId, nc)
		}
	case string:
		var errColName error = nil
		colId, errColName = t.ColIdByName(v)
		if errColName != nil {
			return t, errColName
		}
	default:
		return t,
			fmt.Errorf("Filter: can't handle column "+
				"index of type %v\n", v)
	}
	filteredRows := [][]string{}
	for _, row := range t.Rows {
		if cond(row[colId]) {
			filteredRows = append(filteredRows, row)
		}
	}
	return Table{Header: t.Header, Rows: filteredRows}, nil
}

// The function GreaterOrEqual returns a cond filtering function that checks if a string converted to a float is greater or equal to a value.
func GreaterOrEqual(threshold float64) func(string) bool {
	return func(s string) bool {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Errorf("GreaterOrEqual: failed "+
				"to parse the string %s: "+
				"%v\n", s, err)
			return false
		}
		return v >= threshold
	}
}
