package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", tagHandle)          //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type tags struct {
	color          []string
	size           []string
	price          []string
	neckline       []string
	dresses_length []string
	fabric_type    []string
	quality_level  []string
}

var (
	tag tags
	db  *sql.DB
)

type product struct {
	id     int
	name   string
	cat_id string
}
type product_tag struct {
	id             int
	produtc_id     string
	cat_id         string
	color          string
	size           string
	price          string
	neckline       string
	dresses_length string
	fabric_type    string
	pattern_type   string
	material       string
	silhouette     string
	quality_level  string
}
type tag_color struct {
	id         int
	produtc_id string
	color      string
}
type tag_size struct {
	id         int
	produtc_id string
	size       string
}
type tag_price struct {
	id         int
	produtc_id string
	price      string
}

func init() {

	tag.color = []string{1: "红", 2: "绿", 3: "蓝", 4: "黄", 5: "青", 6: "紫", 7: "橙", 8: "黑", 9: "白", 10: "灰"}
	tag.size = []string{1: "S", 2: "M", 3: "L", 4: "XL", 5: "XXL"}
	tag.price = []string{1: "0~5", 2: "5~10", 3: "10~15", 4: "15~20", 5: "20~25", 6: "25~30", 7: "30~35", 8: "35~"}
	tag.neckline = []string{1: "圆领", 2: "V领", 3: "中领", 4: "低领", 5: "高领", 6: "无领", 7: "立领", 8: "翻领", 9: "尖领", 10: "方领"}
	tag.dresses_length = []string{1: "直身裙", 2: "鱼尾裙", 3: "A裙", 4: "喇叭裙", 5: "一步裙"}
	tag.fabric_type = []string{1: "纤维", 2: "涤纶", 3: "兔毛", 4: "绒布", 5: "缎"}
	tag.quality_level = []string{1: "1星", 2: "2星", 3: "3星", 4: "4星", 5: "5星"}

	//db, _ = sql.Open("mysql", "test:123456@tcp(127.0.0.1:3306)/demo3?charset=utf8")
	db, _ = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/demo?charset=utf8")

}

func tagHandle(w http.ResponseWriter, r *http.Request) {
	get := r.URL.Query().Get

	var orderby string
	get_orderby := get("orderby")
	if get_orderby == "1" {
		orderby = " order by name asc"
	} else if get_orderby == "2" {
		orderby = " order by name desc"
	}

	//fmt.Println("SELECT * FROM product where cat_id=1 " + orderby)
	rows, err := db.Query("SELECT * FROM product where cat_id=1 " + orderby)

	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	products := []product{}
	for rows.Next() {
		p := product{}
		rows.Scan(&p.id, &p.name, &p.cat_id)
		products = append(products, p)
	}

	rows_tag, err := db.Query("SELECT * FROM product_tag where cat_id=1 ")
	defer rows_tag.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	product_tags := []product_tag{}
	for rows_tag.Next() {
		p := product_tag{}
		rows_tag.Scan(&p.id, &p.produtc_id, &p.cat_id, &p.color, &p.size, &p.price, &p.neckline, &p.dresses_length, &p.fabric_type, &p.pattern_type, &p.material, &p.silhouette, &p.quality_level)
		product_tags = append(product_tags, p)
	}
	//fmt.Println(product_tags)
	//fmt.Println(reflect.TypeOf(get))
	tagHtml(get, products, product_tags)
}

//--------
type tag_html struct {
	url   string
	total int
	tag   string
	val   string
}
type tag_rs struct {
	tags    map[string]tag_html
	product product
}

func tagHtml(get func(string) string, p []product, pt []product_tag) {

	//all:=[]interface{}
	fmt.Println(get("color"))

	th := map[string]map[string]int{}
	th1 := map[string]int{}
	th2 := map[string]int{}
	th3 := map[string]int{}
	c1 := get("color")
	sc1 := strings.Split(c1, "_")

	//c2 := get("dresses_length")
	//sc2 := strings.Split(c2, "_")
	c3 := get("size")
	sc3 := strings.Split(c3, "_")

	for _, v := range pt {

		if v.color != "" {
			e1 := false
			if len(sc3) > 0 {
				for _, v3 := range sc3 {

					if v.size == v3 || v3 == "" {
						e1 = true
					}
				}
			} else {
				e1 = true
			}

			if e1 {
				th1[v.color] = th1[v.color] + 1
			}

		}
		if v.dresses_length != "" {
			th2[v.dresses_length] = th1[v.dresses_length] + 1
		}
		if v.size != "" {
			e3 := false
			if len(sc3) > 0 {
				for _, v1 := range sc1 {

					//os.Exit(-1)
					if v.color == v1 || v1 == "" {
						e3 = true
						break
					}
				}
			} else {
				e3 = true
			}

			if e3 {
				th3[v.size] = th3[v.size] + 1
			}
		}

	}
	th["color"] = th1
	//th["dresses_length"] = th2
	th["size"] = th3
	fmt.Println(th)

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
