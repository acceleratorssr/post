package integration

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	searchv1 "post/api/proto/gen/search/v1"
	"post/search/grpc"
	"post/search/integration/startup"
	"testing"
	"time"
)

type Article struct {
	ID      uint64   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Author  Author   `json:"author"`
	Tags    []string `json:"tags"`
}

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

var artsTest = []Article{
	{
		ID:      1,
		Title:   "health",
		Content: "这段描述主要讨论了长寿与癌症之间的关系，并探讨了细胞老化、干细胞功能衰退以及细胞自我修复机制如程序性细胞死亡和自噬在其中的作用。随着年龄的增长，人体内的细胞会经历一系列变化，包括基因突变的积累、细胞损伤、代谢废物增多以及端粒的缩短，这些因素共同作用导致细胞衰老。细胞更新是通过新细胞替代旧细胞来完成的，这一过程中干细胞起着关键作用。然而，干细胞自身也会随着年龄的增长而失去活力，这限制了它们修复机体的能力。尽管干细胞会衰老，但研究发现一些长寿个体（通常指百岁以上的老人）具有与年轻个体相似的某些生物标志物水平，表明他们的身体在某种程度上保持了年轻状态。这些长寿个体可能具有一些保护性的遗传特征，使他们能够更好地清除损伤细胞，抵抗外界压力，并维持干细胞的功能.癌症的发展涉及细胞的重新编程，这一过程与胚胎发育期间发生的细胞分化相似。例如，甲胎蛋白（AFP）这种通常在胎儿阶段高水平表达的蛋白质，在肝癌患者中也会重新表达。同样，端粒酶在癌症中被激活，使癌细胞能够持续分裂而不死亡。这表明癌症细胞在某种程度上恢复了“年轻”的特性。因此，研究长寿和癌症可以互相借鉴。长寿研究的一个策略是避免重大疾病，因为许多长寿者往往没有或延迟发生诸如癌症、心血管疾病等严重的老年性疾病。这一观察支持了“病患压缩”（compression of morbidity）的概念，即健康寿命的延长与疾病发作的延迟相关联。综上所述，长寿与癌症的研究揭示了两者之间的复杂关系，并为理解衰老、疾病预防以及健康寿命的延长提供了重要线索。",
		Author: Author{
			Id:   1,
			Name: "gopher",
		},
	},
	{
		ID:      2,
		Title:   "backend AD",
		Content: "为何考软考？\n提升就业竞争力：\n软考证书是国家认可的专业资质，有助于提升求职者的简历吸引力。\n考试内容涵盖广泛的软件工程知识，有助于拓宽技术视野。\n提升后端开发技能：\n考试内容由经验丰富的专家编写，系统地涵盖了软件设计、开发和测试等方面的知识。\n通过备考，可以系统地学习并提高技术水平。\n积分落户：\n在一些大城市，软考证书可以帮助满足积分落户的要求。\n个人所得税减免：\n根据相关政策，取得软考证书后可享受个人所得税专项附加扣除。\n软考级别及内容\n软考分为不同的级别，每个级别有不同的侧重点：\n\n初级：基础概念，如软件工程、数据库原理、计算机网络等。\n中级：具体技术和方法，如UML建模、需求分析、软件测试等。\n高级：解决实际问题、大型系统设计和管理等。\n后端开发所需技术\n除了软考之外，后端开发者还需掌握以下技术：\n\n编程语言：如Java、Python、C#等。\n框架：如Spring、Django、ASP.NET等。\n数据库：如MySQL、PostgreSQL、MongoDB等。服务器配置与部署：如Nginx、Apache等。前端基础：如HTML、CSS、JavaScript等。软考与后端开发的关系软考的知识点与后端开发密切相关，特别是系统设计和软件工程理论部分，有助于提升项目管理和需求分析能力，从而使项目架构更合理、系统更稳定可靠。结论对于希望成为优秀的后端开发者而言，软考不仅是简历上的加分项，更是帮助系统化学习软件工程知识的有效途径。结合实际项目经验和深入研究相关技术，将有助于在后端开发领域走得更远。",
		Author: Author{
			Id:   2,
			Name: "AD",
		},
	},
	{
		ID:      3,
		Title:   "backend golang",
		Content: "必须掌握的技术：AI\nAI应用：可以用来获取后端开发的知识体系，如通过提问来获取详细的介绍。\n学习AI：建议转向学习AI技术，因为它可以帮助新人更好地适应IT行业，并且AI可以极大提高开发效率。\n基础网络知识\nTCP/IP协议：了解TCP、UDP、HTTP、HTTPS、HTTP2等协议及其区别。\n网络层次：理解TCP与HTTP处于网络七层模型中的不同层次。\n数据库知识\n常用数据库：熟悉PostgreSQL和MySQL，并了解其优化、稳定性和性能相关的技术。\n其他类型数据库：了解图数据库（如Neo4J）、文档型数据库（如MongoDB）、键值存储（如Redis）以及矢量数据库（如Faiss）。\n服务器知识\nWeb服务器：了解Nginx等Web服务器的配置，以及针对不同语言的服务器（如Tomcat、Gunicorn）。\n配置调试：掌握服务器配置，能够解决因配置不当导致的问题。\n某一门开发语言及框架\n语言选择：掌握Java、Python或Go等语言的基本语法、核心库、设计模式、虚拟机等。\n框架学习：熟悉Web框架（如Spring、Django）、ORM框架、SSO/OAuth等。\n安全相关的知识\n防止注入攻击：使用PreparedStatement、ORM或参数化查询等方式防止SQL注入等。\n防御CSRF和XSS攻击：使用特定令牌或内置框架功能防护CSRF，对输出内容进行转义或清理防止XSS。\n认证与授权：确保只有授权用户能访问敏感资源，并对敏感数据进行加密。\n日志与监控：记录应用程序活动，并避免在错误消息中泄露过多信息。\n依赖库管理：定期检查和更新依赖库，确保没有已知的安全漏洞。\n总结\n本文作者认为，随着AI技术的发展，学习AI已经成为后端开发者的一项重要技能。掌握了AI工具后，可以更有效地学习和掌握其他后端开发所需的各项技术。此外，网络安全的重要性日益凸显，后端开发者需要具备一定的网络安全意识和技术来保障系统的安全。",
		Author: Author{
			Id:   3,
			Name: "golang",
		},
	},
	{
		ID:      4,
		Title:   "city",
		Content: "广东城市相关的第一热度话题永远都是，广州衰落了吗？广州已经跌出四大一线城市之列了吗？广州GDP被重庆超越，还会被XXXXX超越吗？广州为什么......广州都已经肉眼可见地越来越拉胯了，旁边那个以“广佛同城”为最大卖点，指望吸收广州发展红利和资源福利的佛山，又能好到哪里去呢？但是在我看来，无论是广州衰落论，还是佛山不行论，最主要的原因还是这几年整个经济大环境就这样，资源本来就有限，以前广东各个城市GDP雄霸全国的时代，那是个什么时代？那是本世纪初，加入WTO，世界工厂地位已定，疯狂做贸易搞基建全国大发展的制造业大增量时代。广东依赖地理位置以及各种政策优势，率先跑到了前面。我依稀记得，当年按照GDP统计城市排名时，很多人在前面看到东莞时的震撼，觉得这不就是一个打工人聚集打工挣钱的地儿嘛，怎么城市排名会这么高。其实对于很多来珠三角打工的人来说，整个珠三角都是东莞一样的存在，只不过不同地方聚集了不同的产业，去打工的工厂类型不一样罢了。而到了后面，长三角等其他地方发展也跟上来了，以前的那种遥遥领先的感觉就没有了，但至少还是相提并论的存在。但是到了这两年，整个大环境都不好，质量高的产业来来回回就那么几个，这是一个争夺存量市场的格局。当别的城市各出奇招争夺这些优质产业，而有的城市不做动作，或者昏招频出时，结果自然就是那么几个好产业被别的城市抢走了。广州跟上海争夺特斯拉工厂失败就是典型的案例。广州有限的资源不集中搞一两个产业，而是选择撒芝麻一样，啥都搞一点，也是典型的案例。说起来广州也是有金融城，生物岛，互联网电商总部区，汽车新能源小镇，各种产业园的地方。实际上呢？金融城是卖家具的，生物岛……至于佛山…嗯，佛山也有金融城，大概还不如广州的卖家具的…再往前倒，还有老生常谈的，网易总部在哪里，现在网易究竟是哪里的企业这种烂话题。深圳这几年能先超越广州，同时还保持着作为一线城市的竞争力，最重要的就是十多年前开始押宝新能源，生物技术等产业，同时也利用自身的各种政策优势巩固几个传统优势产业，在争夺优质产业资源大战中，保住了自己本来就有的优势产业的同时，新增市场的红利也吃到了不少。而广州这些年给人的整体印象就是，自身孵化的还不错的新兴产业公司，都存在着一个做大了之后，到底是搬去深圳，还是京沪的问题，主打的就是一个留不住。试想一下，如果深圳和广州一样，金融电子互联网没保住，新能源生物科技也一个都没发展起来，那深圳也早就被踢出一线城市之列了，从此之后中国只有京沪，没有其他。至于佛山，我觉得把广州各个方面被吐槽的事儿，调一下参数，就是佛山的版本。",
		Author: Author{
			Id:   4,
			Name: "people",
		},
	},
	{
		ID:      5,
		Title:   "city2",
		Content: "人矿跟煤矿一样，都有挖到头的时候。山东也有人抱怨说的，腾笼换鸟，吊用没有。实际上，对于这个问题的理解，全国没有比江苏感触更深的。江苏腾笼换鸟太早了，我小时候就开始了，十多年前就开始了，投资门槛全国最高，各种门槛。你想投资，想花钱，也要看看，你有没有那个水平。这都多少年了，砸了多少钱了？说转型升级，说腾笼换鸟。这是一天两天的事情？很明显，广东不可能再走过去的老路了。但转型升级，也不是一蹴而就的。我上周，路上遇到一个老头子，在红绿灯这里，卖他的不知道什么产品。都多少年，没遇到过这种人了。这都是近段时间流入的人口质量，严重下滑。这样的人，在我们这里没什么用的。满大街都是体面人是有原因的。因为不体面的人，连就业机会都没有。这不是今天才这样，十多年前，我们这里的厂子就开始往东南亚走了。我小时候，早就有移民潮了。广东改开最早，但一直拖到今天。江苏发展很快吗？一直都是相对平稳的速度发展。只不过更看重发展质量。我们也有过，所谓发大财的年头，大量的终端大企业，90年代，家喻户晓。但是这些东西，吊用不大，只在发展初期，快速扩张有用。我小时候，就开始按照平均面积算GDP了。你光绝对值高没有用。你需要在有限的空间，做大才行。那势必是会淘汰垃圾产业的。要么走人，要么完蛋，两条路给你。你看房地产均价，江苏连前五都进不了。炒房也给你拉闸。在这种条件下，短期之所以还没有出大事。完全是因为过去的超高增速的惯性。所有的商业模式，早就习惯了这样的大环境。你现在再差能差到哪里去？我们一直都是这样差的。过去有钱不要，现在没钱，还是不要。有什么区别？但是产业升级，国产替代，实打实的确实带来了，极为丰厚的利润。一部分产业在消亡，一部分产业加速成长。这对江苏来说，不是什么新鲜事。过去一直都是这样的。过去十多年，一直都是淘汰赛。那么多企业完蛋，太多太多太多企业完蛋了。广东一直拖到了今天。我对两地产业是有些理解的，我是产业一线实打实呆了很久的。我几年前就说过，什么时候广东开始拆那些作坊了，什么时候广东开始走上正轨。否则，那些产业，我连看都不要看一眼的。就没几个能留下来的。",
		Author: Author{
			Id:   5,
			Name: "TPeople",
		},
	},
}

type SearchTestSuite struct {
	suite.Suite
	searchSvc *grpc.SearchServiceServer
	syncSvc   *grpc.SyncServiceServer
}

func (s *SearchTestSuite) SetupSuite() {
	s.searchSvc = startup.InitSearchServer()
	s.syncSvc = startup.InitSyncServer()
}

func (s *SearchTestSuite) TestSearch() {
	tests := []struct {
		name          string
		expression    string
		limit         int32
		expectedNum   int
		expectedTitle string
	}{
		{
			"search",
			"软考需要学什么 我想考证书",
			10,
			2,
			"backend AD",
		},
	}

	//s.initData()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := s.searchSvc.Search(ctx, &searchv1.SearchRequest{
				Expression: tt.expression,
				Limit:      tt.limit,
			})
			require.NoError(s.T(), err)
			for _, v := range resp.Article.Articles {
				fmt.Println("resp: ", v.Title)
			}
			assert.Equal(s.T(), tt.expectedTitle, resp.Article.Articles[0].Title)
		})
	}
}

func (s *SearchTestSuite) initData() {
	t := s.T()
	for _, v := range artsTest {
		_, err := s.syncSvc.InputArticle(context.Background(), &searchv1.InputArticleRequest{
			Article: &searchv1.Article{
				Id:      v.ID,
				Title:   v.Title,
				Content: v.Content,
				Author: &searchv1.Author{
					Id:   v.Author.Id,
					Name: v.Author.Name,
				},
			},
		})
		require.NoError(t, err)
	}
}

type Tags struct {
	ObjType string   `json:"obj_type"`
	ObjId   int64    `json:"obj_id"`
	Tags    []string `json:"tags"`
}

func TestSearchService(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}
