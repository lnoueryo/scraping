package main

import (
	_ "embed"
	"net/http"
	// "encoding/csv"
	// "io"
	// "os"
	"database/sql"
	"fmt"
	"strings"
	"log"
	"github.com/PuerkitoBio/goquery"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	webPage := ("https://suumo.jp/jj/chintai/ichiran/FR301FC001/?ar=030&bs=040&ta=11&sc=11201&cb=0.0&ct=9999999&mb=0&mt=9999999&et=9999999&cn=9999999&shkr1=03&shkr2=03&shkr3=03&shkr4=03&sngz=&po1=25&pc=50")
	// webPage := ("https://suumo.jp/chintai/jnc_000074294413/?bc=100289520382")
	// webPage := ("https://suumo.jp/jj/chintai/ichiran/FR301FC001/?ar=030&bs=040&ta=11&sc=11201&cb=0.0&ct=9999999&mb=0&mt=9999999&et=9999999&cn=9999999&shkr1=03&shkr2=03&shkr3=03&shkr4=03&sngz=&po1=25&pc=50&page=2")
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

	// ここでタイトルを取得
	doc.Find("div.cassetteitem a").Each(func(i int, s *goquery.Selection) {
		band, ok := s.Attr("href")
		if ok {
			fmt.Printf(strings.Replace(band, "javascript:void(0);", "", -1))
		}
	})
	// fmt.Println(title)

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