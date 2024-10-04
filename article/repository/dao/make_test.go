package dao

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestS(t *testing.T) {
	file, err := os.Create("articles.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	seed := time.Now().UnixNano()
	for i := 1; i <= 10000; i++ {
		article := ArticleReader{
			ID:       uint64(i),
			Title:    getRandomChineseText(seed, 20),
			Content:  getRandomChineseText(seed, 2000),
			Authorid: uint64(i % 100),
			Ctime:    int64(i * 1000),
			Utime:    int64(i * 1000),
			SnowID:   int64(i),
		}
		// 写入CSV，使用逗号分隔
		line := fmt.Sprintf("%d,%s,%s,%d,%d,%d,%d\n",
			article.ID, article.Title, article.Content, article.Authorid, article.Ctime, article.Utime, article.SnowID)
		file.WriteString(line)
	}

	fmt.Println("Data generation completed.")
}
