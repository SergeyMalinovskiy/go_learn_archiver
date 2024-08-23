package vlc

import (
	"awesome-archiver/lib/compression/vlc/table"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"strings"
	"unicode"
)

type EncoderDecoder struct {
	tblGenerator table.Generator
}

func New(tblGenerator table.Generator) EncoderDecoder {
	return EncoderDecoder{tblGenerator: tblGenerator}
}

func (ed EncoderDecoder) Encode(str string) []byte {
	//str = prepareText(str)

	tbl := ed.tblGenerator.NewTable(str)

	encoded := encodeBin(str, tbl)

	return buildEncodedFile(tbl, encoded)
}

func (ed EncoderDecoder) Decode(encodedData []byte) string {
	tbl, data := parseFile(encodedData)

	return tbl.Decode(data)
}

func parseFile(data []byte) (table.EncodingTable, string) {
	const (
		tableSizeBytesCount = 4
		dataSizeBytesCount  = 4
	)

	tableSizeBinary, data := data[:tableSizeBytesCount], data[tableSizeBytesCount:]
	dataSizeBinary, data := data[:dataSizeBytesCount], data[dataSizeBytesCount:]

	tableSize := binary.BigEndian.Uint32(tableSizeBinary)
	dataSize := binary.BigEndian.Uint32(dataSizeBinary)

	tblBinary, data := data[:tableSize], data[tableSize:]

	tbl := decodeTable(tblBinary)

	body := NewBinChunks(data).Join()

	return tbl, body[:dataSize]

}

func buildEncodedFile(tbl table.EncodingTable, data string) []byte {
	encodedTbl := encodeTable(tbl)

	var buf bytes.Buffer

	buf.Write(encodeInt(len(encodedTbl)))
	buf.Write(encodeInt(len(data)))
	buf.Write(encodedTbl)
	buf.Write(splitByChunks(data, chunksSize).Bytes())

	return buf.Bytes()
}

func encodeInt(num int) []byte {
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, uint32(num))

	return res
}

func encodeTable(tbl table.EncodingTable) []byte {
	var tableBuf bytes.Buffer

	if err := gob.NewEncoder(&tableBuf).Encode(tbl); err != nil {
		log.Fatal("can't serialize table: ", err)
	}

	return tableBuf.Bytes()
}

func decodeTable(tblBinary []byte) table.EncodingTable {
	var tbl table.EncodingTable

	r := bytes.NewReader(tblBinary)
	if err := gob.NewDecoder(r).Decode(&tbl); err != nil {
		log.Fatal("can't serialize table: ", err)
	}

	return tbl
}

// encodeBin encodes str into binary codes string without spaces
func encodeBin(str string, table table.EncodingTable) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch, table))
	}

	return buf.String()
}

func bin(ch rune, table table.EncodingTable) string {
	res, ok := table[ch]
	if !ok {
		panic("unknown character: " + string(ch))
	}

	return res
}

//func getEncodingTable() encodingTable {
//	return encodingTable{
//		' ': "11",
//		't': "1001",
//		'n': "10000",
//		's': "0101",
//		'r': "01000",
//		'd': "00101",
//		'!': "001000",
//		'c': "000101",
//		'm': "000011",
//		'g': "0000100",
//		'b': "0000010",
//		'v': "00000001",
//		'k': "0000000001",
//		'q': "000000000001",
//		'e': "101",
//		'o': "10001",
//		'a': "011",
//		'i': "01001",
//		'h': "0011",
//		'l': "001001",
//		'u': "00011",
//		'f': "000100",
//		'p': "0000101",
//		'w': "0000011",
//		'y': "0000001",
//		'j': "000000001",
//		'x': "00000000001",
//		'z': "000000000000",
//	}
//}

// prepareText prepares text to be fit for encode
// changes upper case letters to: ! + lower case letter
func prepareText(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(ch))
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

// exportText is opposite to prepareText, it prepares decoded text to export:
// it changes: ! + <lower case letter> -> to upper case letter.
// i.g.: !my name is !ted -> My name is Ted
func exportText(str string) string {
	var buf strings.Builder

	var isCapital bool

	for _, ch := range str {
		if isCapital {
			buf.WriteRune(unicode.ToUpper(ch))
			isCapital = false

			continue
		}

		if ch == '!' {
			isCapital = true

			continue
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}
