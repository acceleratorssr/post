package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"post/article/domain"
)

type ArticleCompressionByGZIP struct {
}

func (a *ArticleCompressionByGZIP) Compressed(article *domain.Article) ([]byte, error) {
	data, err := json.Marshal(article)
	if err != nil {
		return nil, err
	}

	compressedData, err := compress(data)
	if err != nil {
		return nil, err
	}
	return compressedData, nil
}

func (a *ArticleCompressionByGZIP) Decompress(data []byte, art *domain.Article) error {
	decompressedData, err := decompress(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decompressedData, &art)
	if err != nil {
		return err
	}
	return nil
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reader.Close()
	return result, nil
}

func NewArticleCompressionByGZIP() Compression {
	return &ArticleCompressionByGZIP{}
}
