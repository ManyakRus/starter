// Package data_packer -- компрессор данных движка хранилища
package data_packer

import (
	"bytes"
	"fmt"

	"github.com/klauspost/compress/zstd"
)

// DataPacker -- архиватор данных движка хранилища
type DataPacker struct {
	bufIn  *bytes.Buffer // Временный буфер для входящих данных
	bufOut *bytes.Buffer // Временный буфер для исходящих данных
	enc    *zstd.Encoder
	dec    *zstd.Decoder
}

// NewDataPacker -- возвращает новый *DataPacker
func NewDataPacker() *DataPacker {
	sf := &DataPacker{
		bufIn:  bytes.NewBuffer([]byte{}),
		bufOut: bytes.NewBuffer([]byte{}),
	}
	// Здесь ошибки быть не может
	sf.enc, _ = zstd.NewWriter(sf.bufOut, zstd.WithEncoderLevel(3))
	sf.dec, _ = zstd.NewReader(sf.bufOut) // Куда писать выход
	return sf
}

// Pack -- сжимает данные для сохранения
func (sf *DataPacker) Pack(binIn []byte) (binOut []byte, err error) {
	binOut = sf.enc.EncodeAll(binIn, binOut)
	return binOut, nil
}

// Unpack -- разжимает данные для отдачи
func (sf *DataPacker) Unpack(binIn []byte) (binOut []byte, err error) {
	binOut, err = sf.dec.DecodeAll(binIn, binOut)
	if err != nil {
		return nil, fmt.Errorf("DataPacker.Unpack(): in decode zstd, err=\n\t%w", err)
	}
	return binOut, nil
}
