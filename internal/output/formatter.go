package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/olekukonko/tablewriter"
)

func FormatUnixTime(unix int) string {
	return time.Unix(int64(unix)/1000, 0).Format("2006-01-02 15:04:05")
}

type Formatter interface {
	Format(data interface{}) error
}

type JSONFormatter struct {
	writer io.Writer
}

type TableFormatter struct {
	writer   io.Writer
	noHeader bool
}

func NewJSONFormatter(writer io.Writer) *JSONFormatter {
	if writer == nil {
		writer = os.Stdout
	}
	return &JSONFormatter{writer: writer}
}

func NewTableFormatter(writer io.Writer, noHeader bool) *TableFormatter {
	if writer == nil {
		writer = os.Stdout
	}
	return &TableFormatter{writer: writer, noHeader: noHeader}
}

func PrintAsJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func NewTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	return table
}

func (f *JSONFormatter) Format(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *TableFormatter) Format(data interface{}) error {
	if data == nil {
		return nil
	}

	// Handle slice of structs
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		return f.formatSlice(v)
	}

	// Handle single struct
	if v.Kind() == reflect.Struct {
		return f.formatStruct(v)
	}

	// Fallback to JSON for complex types
	jsonFormatter := NewJSONFormatter(f.writer)
	return jsonFormatter.Format(data)
}

func (f *TableFormatter) formatSlice(v reflect.Value) error {
	if v.Len() == 0 {
		_, _ = fmt.Fprintln(f.writer, "No data found")
		return nil
	}

	first := v.Index(0)
	if first.Kind() != reflect.Struct {
		return f.formatSimpleSlice(v)
	}

	table := tablewriter.NewWriter(f.writer)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetCenterSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Extract headers from struct fields
	headers := f.extractHeaders(first.Type())
	if !f.noHeader {
		table.SetHeader(headers)
	}

	// Add rows
	for i := 0; i < v.Len(); i++ {
		row := f.extractRow(v.Index(i))
		table.Append(row)
	}

	table.Render()
	return nil
}

func (f *TableFormatter) formatStruct(v reflect.Value) error {
	table := tablewriter.NewWriter(f.writer)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetCenterSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	if !f.noHeader {
		table.SetHeader([]string{"Field", "Value"})
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		fieldValue := f.formatValue(v.Field(i))
		table.Append([]string{fieldName, fieldValue})
	}

	table.Render()
	return nil
}

func (f *TableFormatter) formatSimpleSlice(v reflect.Value) error {
	table := tablewriter.NewWriter(f.writer)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetCenterSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	if !f.noHeader {
		table.SetHeader([]string{"Value"})
	}

	for i := 0; i < v.Len(); i++ {
		value := f.formatValue(v.Index(i))
		table.Append([]string{value})
	}

	table.Render()
	return nil
}

func (f *TableFormatter) extractHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		headers = append(headers, field.Name)
	}
	return headers
}

func (f *TableFormatter) extractRow(v reflect.Value) []string {
	var row []string
	for i := 0; i < v.NumField(); i++ {
		if !v.Type().Field(i).IsExported() {
			continue
		}
		value := f.formatValue(v.Field(i))
		row = append(row, value)
	}
	return row
}

func (f *TableFormatter) formatValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return f.formatValue(v.Elem())
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", v.Float())
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return ""
		}
		return fmt.Sprintf("[%d items]", v.Len())
	case reflect.Struct:
		// For nested structs, show a summary instead of full structure
		typeName := v.Type().Name()
		if typeName == "" {
			typeName = "struct"
		}
		return fmt.Sprintf("<%s>", typeName)
	default:
		// Avoid showing memory addresses and internal representations
		typeName := v.Type().Name()
		if typeName != "" {
			return fmt.Sprintf("<%s>", typeName)
		}
		return fmt.Sprintf("<%s>", v.Type().String())
	}
}
