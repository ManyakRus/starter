// Package data_packer -- компрессор данных движка хранилища
package data_packer

import (
	"fmt"
	"github.com/klauspost/compress/zstd"
	//"sync"
)

// DataPacker -- архиватор данных движка хранилища
type DataPacker struct {
	enc *zstd.Encoder
	dec *zstd.Decoder
}

// NewDataPacker -- возвращает новый *DataPacker
func NewDataPacker() *DataPacker {
	dp := &DataPacker{}

	dp.enc, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
	dp.dec, _ = zstd.NewReader(nil)

	return dp
}

func (dp *DataPacker) Close() error {
	var errs []error

	if dp.enc != nil {
		if err := dp.enc.Close(); err != nil {
			errs = append(errs, fmt.Errorf("encoder close: %w", err))
		}
	}

	if dp.dec != nil {
		dp.dec.Close()
	}

	// Возвращаем первую ошибку (если есть)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Pack -- сжимает данные для сохранения
func (dp *DataPacker) Pack(binIn []byte) (binOut []byte, err error) {
	return dp.enc.EncodeAll(binIn, nil), nil
}

// Unpack -- разжимает данные для отдачи
func (dp *DataPacker) Unpack(binIn []byte) ([]byte, error) {
	return dp.dec.DecodeAll(binIn, nil)
}
