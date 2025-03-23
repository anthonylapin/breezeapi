package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func gzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gzipWriter := gzip.NewWriter(&b)
	gzipWriter.Write(data)
	gzipWriter.Close()

	compressedData := b.Bytes()
	
	return compressedData, nil
}

var COMPRESS_STRATEGIES = map[string]func([]byte) ([]byte, error){
	"gzip": gzipCompress,
}

func compressionSupported(compressionType string) bool {
	_, strategySupported := COMPRESS_STRATEGIES[compressionType]
	return strategySupported
}

func getCompressedData(compressionType string, data []byte) ([]byte, error) {
	strategy, strategySupported := COMPRESS_STRATEGIES[compressionType]

	if !strategySupported {
		return nil, fmt.Errorf("Compression strategy not supported")
	}

	return strategy(data)
}