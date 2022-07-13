package main

import (
	_ "embed"
	"encoding/csv"
	"net/http"
	"os"
	"strconv"

	// "encoding/csv"
	// "io"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	prefectureCode := 13
	areaCode := 101
	for {
		index := 1
		cityCode := strconv.Itoa(prefectureCode) + strconv.Itoa(areaCode)
		fmt.Println(fmt.Sprintf("city code %s", cityCode))
		for {
			fmt.Println(fmt.Sprintf("page %s start", strconv.Itoa(index)))
			webPage := fmt.Sprintf("https://suumo.jp/jj/chintai/ichiran/FR301FC001/?ar=030&bs=040&ta=%s&sc=%s&cb=0.0&ct=9999999&mb=0&mt=9999999&et=9999999&cn=9999999&shkr1=03&shkr2=03&shkr3=03&shkr4=03&sngz=&po1=25&pc=50&page=%s", strconv.Itoa(prefectureCode), cityCode, strconv.Itoa(index))
			resp, err := http.Get(webPage)
			if err != nil {
				log.Printf("failed to get html: %s", err)
			}
			defer resp.Body.Close()
	
			if resp.StatusCode != 200 {
				log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
			}
	
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("failed to load html: %s", err)
			}
	
			detailURLs := []string{}
			log.Println("start to get URLs")
			doc.Find("div.cassetteitem a").Each(func(i int, s *goquery.Selection) {
				band, ok := s.Attr("href")
				if ok {
					url := strings.Replace(band, "javascript:void(0);", "", -1)
					if url != "" {
						detailURLs = append(detailURLs, url)
					}
				}
			})
	
			if len(detailURLs) == 0 {
				log.Println("nothing more")
				break
			}
	
			log.Println("start to get records")
			records := [][]string{}
			for _, url := range detailURLs {
				texts := getDetail(url, cityCode)
				if len(texts) == 0 {
					continue
				}
				records = append(records, texts)
				log.Println(texts)
			}
	
			log.Println("start to write records")
			filename := cityCode + ".csv"
			f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal(err)
			}
	
			w := csv.NewWriter(f)
	
			if index == 1 {
				head := []string{
					"prefecture_code",
					"city_code",
					"price",
					"management_fee",
					"security_deposit",
					"key_money",
					"prefecture",
					"city",
					"station",
					"address",
					"layout",
					"area",
					"age",
					"floor",
				}
				w.Write(head)
			}
			w.WriteAll(records)
	
			w.Flush()
			index += 1
		}
		if areaCode == 700 {
			areaCode = 101
			break
		}
		areaCode += 1
	}
	// f, _ := os.Open("./users.csv")
	// defer f.Close()
	// r := csv.NewReader(f)
	// keys := make(map[string]string)
	// for {
	// 	row, err := r.Read() // csvを1行ずつ読み込む
	// 	if err == io.EOF {
	// 		break // 読み込むべきデータが残っていない場合，Readはnil, io.EOFを返すため、ここでループを抜ける
	// 	}
	// 	v, ok := keys[row[4]];if ok {
	// 		fmt.Println(v)
	// 	}
	// 	keys[row[4]] = row[4]
	//   }
}

func getDetail(path string, cityCode string) ([]string) {
	tds := []string{cityCode[:2], cityCode}
	url := "https://suumo.jp" + path
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to get html: %s", err)
		return tds
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
		return tds
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("failed to load html: %s", err)
		return tds
	}

	doc.Find("div.property_view_note-info").Each(func(i int, s *goquery.Selection) {
		s.Find("span").Each(func(i int, s *goquery.Selection) {
			if i < 4 {
				conditions := []string{
					"管理費・共益費:",
					"万円",
					"円",
					"敷金:",
					"礼金:",
				}
				text := s.Text()
				for _, condition := range conditions {
					text = strings.Replace(text, condition, "", -1)
				}
				text = strings.TrimSpace(text)
				if text == "-" {
					text = "0"
				}
				tds = append(tds, text)
			}
		})
	})

	doc.Find("div.left_column ul li").Each(func(i int, s *goquery.Selection) {
		if i > 1 && i < 5 {
			a := s.Find("a")
			tds = append(tds, a.Text())
		}
	})

	doc.Find("table.property_view_table").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if i != 1 && i < 7 {
				text := s.Text()
				if text == "新築" {
					text = "0"
				}
				conditions := []string{
					"m2",
					"築",
					"年",
					"階",
				}
				for _, condition := range conditions {
					text = strings.Replace(text, condition, "", -1)
				}
				tds = append(tds, text)
			}
		})
	})
	return tds
}

func connectDB() {
    dbconf := "root:Popo@6252@tcp(127.0.0.1:3306)/stocklocator-unicharm?charset=utf8mb4"

    db, err := sql.Open("mysql", dbconf)
	
    // 接続が終了したらクローズする
    defer db.Close()

    if err != nil {
		fmt.Println(err.Error())
    }
	
    err = db.Ping()
	
    if err != nil {
		fmt.Println("データベース接続失敗")
		return
	} else {
		fmt.Println("データベース接続成功")
	}
	createPoint(db)
}

func createCSV() {
	// head := []string{
	// 	"店舗コード",
	// 	"店舗",
	// 	"住所",
	// 	"電話番号",
	// 	"納品日",
	// 	"納品数",
	// 	"最終納品日",
	// 	"最終納品数",
	// 	"商品コード",
	// 	"商品名",
	// 	"商品並び順",
	// 	"カテゴリーコード",
	// 	"カテゴリー名",
	// 	"カテゴリー並び順",
	// 	"サブカテゴリーコード",
	// 	"サブカテゴリー名",
	// 	"サブカテゴリー並び順",
	// 	"ブランドコード",
	// 	"ブランド名",
	// 	"ブランド並び順",
	// 	"サブブランドコード",
	// 	"サブブランド名",
	// 	"サブブランド並び順",
	// 	"画像",
	// }
}

func createPoint(db *sql.DB) {
	// head := []string{
	// 	"id",
	// 	"name",
	// 	"address",
	// 	"phone",
	// 	"lat",
	// 	"lng",
	// 	"location_type",
	// }
	query := `INSERT INTO point (id, name, address, phone, lat, lng, location_type) VALUES `
	values := `(1, "井上領", "松原", 123456789, 123.000, 40.000, "ROOFTOP"),(1, "井上領", "松原", 123456789, 123.000, 40.000, "ROOFTOP")`
	_, err := db.Query(query + values)
	if err != nil {
		fmt.Println(err)
	}
}