package compression

import "post/article/domain"

type Compression interface {
	Compressed(article *domain.Article) ([]byte, error)
	Decompress(data []byte, art *domain.Article) error
}
