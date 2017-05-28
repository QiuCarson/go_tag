package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func init() {

	tag.color = []string{1: "红", 2: "绿", 3: "蓝", 4: "黄", 5: "青", 6: "紫", 7: "橙", 8: "黑", 9: "白", 10: "灰"}
	tag.size = []string{1: "S", 2: "M", 3: "L", 4: "XL", 5: "XXL"}
	tag.price = []string{1: "0~5", 2: "5~10", 3: "10~15", 4: "15~20", 5: "20~25", 6: "25~30", 7: "30~35", 8: "35~"}
	tag.neckline = []string{1: "圆领", 2: "V领", 3: "中领", 4: "低领", 5: "高领", 6: "无领", 7: "立领", 8: "翻领", 9: "尖领", 10: "方领"}
	tag.dresses_length = []string{1: "直身裙", 2: "鱼尾裙", 3: "A裙", 4: "喇叭裙", 5: "一步裙"}
	tag.fabric_type = []string{1: "纤维", 2: "涤纶", 3: "兔毛", 4: "绒布", 5: "缎"}
	tag.quality_level = []string{1: "1星", 2: "2星", 3: "3星", 4: "4星", 5: "5星"}

	db, _ = sql.Open("mysql", "test:123456@tcp(127.0.0.1:3306)/demo3?charset=utf8")
	//db, _ = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/demo?charset=utf8")

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
	str := tagHtml(get, products, product_tags)
	io.WriteString(w, str)
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

func tagHtml(get func(string) string, p []product, pt []product_tag) string {
	th := map[string]map[string]int{}
	//all:=[]interface{}
	fmt.Println(get("color"))
	tag_color_get := get("color")
	tag_size_get := get("size")
	tag_neckline_get := get("neckline")
	tag_dresses_length_get := get("dresses_length")
	tag_fabric_type_get := get("fabric_type")
	tag_quality_level_get := get("quality_level")

	tag_color := map[string]int{}
	tag_size := map[string]int{}
	tag_neckline := map[string]int{}
	tag_dresses_length := map[string]int{}
	tag_fabric_type := map[string]int{}
	tag_quality_level := map[string]int{}

	tag_color_split := strings.Split(tag_color_get, "_")
	tag_size_split := strings.Split(tag_size_get, "_")
	tag_neckline_split := strings.Split(tag_neckline_get, "_")
	tag_dresses_length_split := strings.Split(tag_dresses_length_get, "_")
	tag_fabric_type_split := strings.Split(tag_fabric_type_get, "_")
	tag_quality_level_split := strings.Split(tag_quality_level_get, "_")

	var tag_bool_color, tag_bool_size, tag_bool_neckline, tag_bool_dresses_length, tag_bool_fabric, tag_bool_quality_level bool

	for _, v := range pt {
		//color
		if v.color != "" {
			tag_bool_size = false
			for _, k := range tag_size_split {
				if v.size == k || k == "" {
					tag_bool_size = true
				}
			}
			tag_bool_neckline = false
			for _, k := range tag_neckline_split {
				if v.neckline == k || k == "" {
					tag_bool_neckline = true
				}
			}
			tag_bool_dresses_length = false
			for _, k := range tag_dresses_length_split {
				if v.dresses_length == k || k == "" {
					tag_bool_dresses_length = true
				}
			}
			tag_bool_fabric = false
			for _, k := range tag_fabric_type_split {
				if v.fabric_type == k || k == "" {
					tag_bool_fabric = true
				}
			}
			tag_bool_quality_level = false
			for _, k := range tag_quality_level_split {
				if v.quality_level == k || k == "" {
					tag_bool_quality_level = true
				}
			}

			if tag_bool_size && tag_bool_neckline && tag_bool_dresses_length && tag_bool_fabric && tag_bool_quality_level {
				tag_color[v.color] = tag_color[v.color] + 1
			}

		}
		//dresses_length
		if v.dresses_length != "" {
			tag_bool_color = false
			for _, k := range tag_color_split {
				if v.color == k || k == "" {
					tag_bool_color = true
				}
			}
			tag_bool_neckline = false
			for _, k := range tag_neckline_split {
				if v.neckline == k || k == "" {
					tag_bool_neckline = true
				}
			}
			tag_bool_size = false
			for _, k := range tag_size_split {
				if v.color == k || k == "" {
					tag_bool_size = true
				}
			}
			tag_bool_fabric = false
			for _, k := range tag_fabric_type_split {
				if v.fabric_type == k || k == "" {
					tag_bool_fabric = true
				}
			}
			tag_bool_quality_level = false
			for _, k := range tag_quality_level_split {
				if v.quality_level == k || k == "" {
					tag_bool_quality_level = true
				}
			}
			if tag_bool_color && tag_bool_neckline && tag_bool_size && tag_bool_fabric && tag_bool_quality_level {
				tag_dresses_length[v.dresses_length] = tag_dresses_length[v.dresses_length] + 1
			}

		}
		//size
		if v.size != "" {
			tag_bool_color = false
			for _, k := range tag_color_split {
				if v.color == k || k == "" {
					tag_bool_color = true
				}
			}
			tag_bool_neckline = false
			for _, k := range tag_neckline_split {
				if v.neckline == k || k == "" {
					tag_bool_neckline = true
				}
			}
			tag_bool_dresses_length = false
			for _, k := range tag_dresses_length_split {
				if v.dresses_length == k || k == "" {
					tag_bool_dresses_length = true
				}
			}
			tag_bool_fabric = false
			for _, k := range tag_fabric_type_split {
				if v.fabric_type == k || k == "" {
					tag_bool_fabric = true
				}
			}
			tag_bool_quality_level = false
			for _, k := range tag_quality_level_split {
				if v.quality_level == k || k == "" {
					tag_bool_quality_level = true
				}
			}

			if tag_bool_color && tag_bool_neckline && tag_bool_dresses_length && tag_bool_fabric && tag_bool_quality_level {
				tag_size[v.size] = tag_size[v.size] + 1
			}

		}
		//neckline
		if v.neckline != "" {
			tag_bool_color = false
			for _, k := range tag_color_split {
				if v.color == k || k == "" {
					tag_bool_color = true
				}
			}
			tag_bool_size = false
			for _, k := range tag_size_split {
				if v.size == k || k == "" {
					tag_bool_size = true
				}
			}
			tag_bool_dresses_length = false
			for _, k := range tag_dresses_length_split {
				if v.dresses_length == k || k == "" {
					tag_bool_dresses_length = true
				}
			}
			tag_bool_fabric = false
			for _, k := range tag_fabric_type_split {
				if v.fabric_type == k || k == "" {
					tag_bool_fabric = true
				}
			}
			tag_bool_quality_level = false
			for _, k := range tag_quality_level_split {
				if v.quality_level == k || k == "" {
					tag_bool_quality_level = true
				}
			}
			if tag_bool_color && tag_bool_size && tag_bool_dresses_length && tag_bool_fabric && tag_bool_quality_level {
				tag_neckline[v.neckline] = tag_neckline[v.neckline] + 1
			}

		}
		//fabric
		if v.fabric_type != "" {
			tag_bool_color = false
			for _, k := range tag_color_split {
				if v.color == k || k == "" {
					tag_bool_color = true
				}
			}
			tag_bool_size = false
			for _, k := range tag_size_split {
				if v.size == k || k == "" {
					tag_bool_size = true
				}
			}
			tag_bool_dresses_length = false
			for _, k := range tag_dresses_length_split {
				if v.dresses_length == k || k == "" {
					tag_bool_dresses_length = true
				}
			}
			tag_bool_neckline = false
			for _, k := range tag_neckline_split {
				if v.neckline == k || k == "" {
					tag_bool_neckline = true
				}
			}
			tag_bool_quality_level = false
			for _, k := range tag_quality_level_split {
				if v.quality_level == k || k == "" {
					tag_bool_quality_level = true
				}
			}
			if tag_bool_color && tag_bool_size && tag_bool_dresses_length && tag_bool_neckline && tag_bool_quality_level {
				tag_fabric_type[v.fabric_type] = tag_fabric_type[v.fabric_type] + 1
			}

		}
		//fabric
		if v.quality_level != "" {
			tag_bool_color = false
			for _, k := range tag_color_split {
				if v.color == k || k == "" {
					tag_bool_color = true
				}
			}
			tag_bool_size = false
			for _, k := range tag_size_split {
				if v.size == k || k == "" {
					tag_bool_size = true
				}
			}
			tag_bool_dresses_length = false
			for _, k := range tag_dresses_length_split {
				if v.dresses_length == k || k == "" {
					tag_bool_dresses_length = true
				}
			}
			tag_bool_neckline = false
			for _, k := range tag_neckline_split {
				if v.neckline == k || k == "" {
					tag_bool_neckline = true
				}
			}
			tag_bool_fabric = false
			for _, k := range tag_fabric_type_split {
				if v.fabric_type == k || k == "" {
					tag_bool_fabric = true
				}
			}
			if tag_bool_color && tag_bool_size && tag_bool_dresses_length && tag_bool_neckline && tag_bool_fabric {
				tag_quality_level[v.quality_level] = tag_quality_level[v.quality_level] + 1
			}

		}

	}
	th["color"] = tag_color
	th["dresses_length"] = tag_dresses_length
	th["size"] = tag_size
	th["neckline"] = tag_neckline
	th["quality_level"] = tag_quality_level
	th["fabric_type"] = tag_fabric_type

	str, err := json.Marshal(th)

	if err != nil {
		fmt.Println(err)
	}
	return string(str)
	//_, _ := io.WriteString(w, str)
	//fmt.Println(th)

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
