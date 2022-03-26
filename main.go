package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	// "fmt"

	"github.com/flosch/pongo2"
	_ "github.com/go-sql-driver/mysql" // Using MySQL driver
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

const tmplPath = "src/template/"

var db *sqlx.DB
var e = createMux()

type info struct {
	Name    string
	Address string
	Pincode int
}

func main() {
	db = connectDB()
	repository.SetDB(db)

	// TOP ページに記事の一覧を表示します。
	e.GET("/", handler.ArticleIndex)

	// 記事に関するページは "/articles" で開始するようにします。
	// 記事一覧画面には "/" と "/articles" の両方でアクセスできるようにします。
	// パスパラメータの ":id" も ":articleID" と明確にしています。
	e.GET("/articles", handler.ArticleIndex)         // 一覧画面
	e.GET("/articles/new", handler.ArticleNew)       // 新規作成画面
	e.GET("/articles/:articleID", handler.ArticleShow)      // 詳細画面
	e.GET("/articles/:articleID/edit", handler.ArticleEdit) // 編集画面

	// HTML ではなく JSON を返却する処理は "/api" で開始するようにします。
	// 記事に関する処理なので "/articles" を続けます。
	e.GET("/api/articles", handler.ArticleList)          // 一覧
	e.POST("/api/articles", handler.ArticleCreate)       // 作成
	e.DELETE("/api/articles/:articleID", handler.ArticleDelete) // 削除

	e.Logger.Fatal(e.Start(":8080"))
}

func createMux() *echo.Echo {
	// アプリケーションインスタンスを生成
	e := echo.New()

	// アプリケーションに各種ミドルウェアを設定
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CSRF())

	// `src/css` ディレクトリ配下のファイルに `/css` のパスでアクセスできるようにする
	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	e.Validator = &CustomValidator{validator: validator.New()}

	// アプリケーションインスタンスを返却
	return e
}

func connectDB() *sqlx.DB {
	dsn := os.Getenv("DSN")
	driver := os.Getenv("DRIVER")
	
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
			e.Logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
			e.Logger.Fatal(err)
	}
	log.Println("db connection succeeded")
	return db
}

func articleIndex(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article Index",
		"Now":     time.Now(),
		}
		return render(c, "article/index.html", data)
}

func articleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}

	return render(c, "article/new.html", data)
}

func articleShow(c echo.Context) error {
	// パスパラメータを抽出
	id, _ := strconv.Atoi(c.Param("id"))

	data := map[string]interface{}{
		"Message": "Article Show",
		"Now":     time.Now(),
		"ID":      id,
	}

	return render(c, "article/show.html", data)
}

func articleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	data := map[string]interface{}{
		"Message": "Article Edit",
		"Now":     time.Now(),
		"ID":      id,
	}

	return render(c, "article/edit.html", data)
}

func htmlBlob(file string, data map[string]interface{}) ([]byte, error) {
	return pongo2.Must(pongo2.FromCache(tmplPath + file)).ExecuteBytes(data)
}

func render(c echo.Context, file string, data map[string]interface{}) error {
	// 定義した htmlBlob() 関数を呼び出し、生成された HTML をバイトデータとして受け取る
	b, err := htmlBlob(file, data)

	// エラーチェック
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// ステータスコード 200 で HTML データをレスポンス
	return c.HTMLBlob(http.StatusOK, b)
}

// CustomValidator ...
type CustomValidator struct {
    validator *validator.Validate
}
	
// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}
