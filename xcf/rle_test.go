package xcf

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"vimagination.zapto.org/byteio"
	"vimagination.zapto.org/memio"
)

func TestReads(t *testing.T) {
	tests := []struct {
		Input, Output string
	}{
		{},
		{
			"\x00A",
			"A",
		},
		{
			"\x01A",
			"AA",
		},
		{
			"\x7eA",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		},
		{
			"\x7f\x00\x00A",
			"",
		},
		{
			"\x7f\x00\x01A",
			"A",
		},
		{
			"\x7f\x00\x0aA",
			"AAAAAAAAAA",
		},
		{
			"\x7f\x01\x00A",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		},
		{
			"\x80\x00\x01A",
			"A",
		},
		{
			"\x80\x00\x02AA",
			"AA",
		},
		{
			"\xffA",
			"A",
		},
		{
			"\xfeAB",
			"AB",
		},
		{
			"\x00A\x01B\x7f\x00\x01C\x7f\x00\x0aD\x80\x00\x0a1234567890",
			"ABBCDDDDDDDDDD1234567890",
		},
	}
	for n, test := range tests {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, &rle{
			Reader: &byteio.StickyBigEndianReader{
				Reader: strings.NewReader(test.Input),
			},
		})
		if err != nil {
			t.Errorf("test %d: unexpected error: %q", n+1, err)
		} else if str := buf.String(); str != test.Output {
			t.Errorf("test %d: expecting %q, got %q", n+1, test.Output, str)
		}
	}
}

func TestWrites(t *testing.T) {
	tests := []struct {
		Input, Output string
	}{
		{},
		{
			"A",
			"\x00A",
		},
		{
			"AA",
			"\x01A",
		},
		{
			"AAA",
			"\x02A",
		},
		{
			"AB",
			"\xfeAB",
		},
		{
			"AAAAABBCCCCCFFGGHHIIII",
			"\x04A\x01B\x04C\xfaFFGGHH\x03I",
		},
		{
			"ABVGGFHSDFGHDFHGFHDGZDBGNHJSDGSDFGHHTGHDDVFVDBFBDBVF",
			"\xccABVGGFHSDFGHDFHGFHDGZDBGNHJSDGSDFGHHTGHDDVFVDBFBDBVF",
		},
	}
	d := make([]byte, 0, 2048)
	for n, test := range tests {
		w := newWriter(memio.Create(&d))
		w.WriteRLE([]byte(test.Input))
		if w.Err != nil {
			t.Errorf("test %d: unexpected error: %q", n+1, w.Err)
		} else if str := string(d); str != test.Output {
			t.Errorf("test %d: expecting %q, got %q", n+1, test.Output, str)
		}
	}
}
