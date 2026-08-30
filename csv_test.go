package dataexchange

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCSVContractRoundTripAndBounds(t *testing.T) {
	var output bytes.Buffer
	encoder := NewCSVEncoder(&output, 64)
	if err := encoder.Write([]string{"name", "amount"}); err != nil {
		t.Fatal(err)
	}
	if err := encoder.Write([]string{"Ada", "12"}); err != nil {
		t.Fatal(err)
	}
	if err := encoder.Close(); err != nil {
		t.Fatal(err)
	}
	var rows []CSVRecord
	headers, err := DecodeCSV(t.Context(), strings.NewReader(output.String()), CSVDecodeLimits{MaxBytes: 64, MaxRows: 1, MaxColumns: 2}, func(_ []string, row CSVRecord) error {
		rows = append(rows, row)
		return nil
	})
	if err != nil || len(headers) != 2 || len(rows) != 1 || rows[0].Values[0] != "Ada" {
		t.Fatalf("headers=%v rows=%v error=%v", headers, rows, err)
	}
	var bounded bytes.Buffer
	writer := BoundedWriter{Writer: &bounded, Limit: 3}
	if _, err := writer.Write([]byte("four")); !errors.Is(err, ErrPayloadTooLarge) || bounded.Len() != 0 {
		t.Fatalf("bounded write error=%v bytes=%q", err, bounded.Bytes())
	}
}
