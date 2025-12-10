# [`tabellchen`](https://owncloud.gwdg.de/index.php/apps/files/?dir=/github_docs&fileid=3024168943): dealing with tables in `Go`

## Quick tutorial

A table has a header and rows:

```
type Table struct {
     Header []string
     Rows [][]string
}
```

To read a headerless tab-separated table from a file:

```
readConfig := tabellchen.ReadConfig{
	   FilePath: table.tsv,
	   Separator: '\t',
	   HasHeader: false,
}

table := tabellchen.ReadTable(readConfig)
```

To filter the table, specify the column to filter for and the
filtering function that accepts a `string` and returns a `bool`. The
package contains factories for some filtering functions, for example:

```go
// keep rows with values >= 0.95 in the fourth column
err := table.Filter(3, tabellchen.GreaterOrEqual(0.95))

// keep rows distinct by the first column
err = table.Filter(0, tabellchen.DistinctByCol())
```

To extract a column:

```go
col, err := table.GetCol(1)
```

To extract unique values from a column:

```go
uniq := table.ColUniq(1)

```

The columns can be reordered:

```go
err := table.ReorderColumns(3, 4, 5, 0, 1, 2)
```

*This package is in development.*