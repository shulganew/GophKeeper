package services

import (
	"io"

	"go.uber.org/zap"
)

const Bblock = 16384

// Decode Gfiles reader.
type DecodeFileReader struct {
	io.ReadCloser
	dKey []byte
}

func (ft DecodeFileReader) Read(b []byte) (totlal int, err error) {
	dataBlock := len(b)
	if len(b)%Bblock != 0 {
		dataBlock = len(b)/Bblock + Bblock
	}
	datas := make([]byte, dataBlock)

	delta := make([]byte, len(b)-dataBlock)
	b = append(b, delta...)

	zap.S().Infoln("Codding start", len(b))
	for i := 0; i < len(b); i += Bblock {
		//time.Sleep(time.Second)
		//	zap.S().Infoln("Codding", i, " ", len(b))
		data, err := DecodeData(ft.dKey, b[i:i+Bblock])
		if err != nil {
			zap.S().Errorln("Can't decode data: ", err)
			return 0, err
		}
		datas = append(datas, data...)
	}

	return ft.ReadCloser.Read(datas)
}

func (ft DecodeFileReader) Close() error {
	return ft.ReadCloser.Close()
}

// Encode Gfiles reader.
type EncodeFileReader struct {
	io.ReadCloser
	dKey []byte
}

func (ft EncodeFileReader) Read(b []byte) (totlal int, err error) {
	dataBlock := len(b)
	if len(b)%Bblock != 0 {
		dataBlock = len(b)/Bblock + Bblock
	}
	datas := make([]byte, dataBlock)

	delta := make([]byte, len(b)-dataBlock)
	b = append(b, delta...)

	zap.S().Infoln("Codding start", len(b))
	for i := 0; i < len(b); i += Bblock {
		//time.Sleep(time.Second)
		//	zap.S().Infoln("Codding", i, " ", len(b))
		data, err := EncodeData(ft.dKey, b[i:i+Bblock])
		if err != nil {
			zap.S().Errorln("Can't decode data: ", err)
			return 0, err
		}
		datas = append(datas, data...)
	}

	return ft.ReadCloser.Read(datas)
}

func (ft EncodeFileReader) Close() error {
	return ft.ReadCloser.Close()
}
