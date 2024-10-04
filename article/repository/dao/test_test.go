package dao

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getRandomChineseText(seed int64, length int) string {
	r := rand.New(rand.NewSource(seed))
	chineseChars := "的一是不了人我在有他这为之大来以个中上们到说时地要就出会可也你对生能而子那得于着下自之年过发后作里用道行所然家种事方多经么去法学如都同现当没动面起看定天分还进好小部其些主样理心她本前开但因只从想实日军者意无力它与长把机十民第公此已工使情明性知全三又关点正业外将两高间由问很最重并物手应战向头文体政美相见被利什二等产新己制身果加西斯月话合回特代内信表化老给世位次度门任常先海通教儿原东声提立及比员解水名真论处走义各入几口认条平系气题活尔更别打女变四神总何电数安少报才结反受目太量再感建务做或接必场件计管期市直德资命山金指克李路风接"

	var result string
	for i := 0; i < length; i++ {
		index := r.Intn(len(chineseChars))
		result += string(chineseChars[index])
	}
	return result
}

func generateTestData(db *gorm.DB, count int) {
	seed := time.Now().UnixNano()
	for i := 0; i < count; i++ {
		article := ArticleReader{
			Title:    fmt.Sprintf("Title %d", i),
			Content:  getRandomChineseText(seed, 2000),
			Authorid: uint64(rand.Intn(1000)),
			Ctime:    rand.Int63(),
			Utime:    rand.Int63(),
			SnowID:   rand.Int63(),
		}
		db.Create(&article)
	}
}

func TestT(t *testing.T) {
	// 初始化 MySQL 连接
	dsn := "root:20031214pzw!@tcp(127.0.0.1:3306)/garden_article"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移
	db.AutoMigrate(&ArticleReader{})

	// 生成 100,000 条测试数据
	generateTestData(db, 100000)
}
