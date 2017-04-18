package main

import (
  "time"
  "strconv"
  "encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
  "./store"
)

type ApiError struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Detail      string `json:"detail"`
}

type Article struct {
	Id          int `json:"id"`
	Text        string `json:"text"`
	Created_by  string `json:"created_by"`
	Created_at  string `json:"created_at"`
}

type Comment struct {
	Id          int `json:"id"`
	Article_id  int `json:"article_id"`
	Text        string `json:"text"`
	Created_by  string `json:"created_by"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
}

type ArticleList struct {
  Articles  []Article   `json:"articles"`
}

type CommentList struct {
  Comments  []Comment   `json:"comments"`
}

func getTS() (ts string) {
  t:= time.Now()
  ts = t.Format( time.RFC3339)
return
}

func postArticle(c echo.Context) error {

	article := new(Article)
	if err := c.Bind(article); err != nil {
    ae:=ApiError{ Code: "1", Title:" Article bind failed", Detail: "Invalid bind" }
	  return c.JSON(http.StatusUnprocessableEntity, ae )
	}

  id,err:= store.GetNextId("articles")
  article.Id = id
  article.Created_at = getTS()
  buf, err := json.Marshal(article)
  err= store.Create("articles",buf,id)

  if err != nil {
    ae:=ApiError{ Code: "2", Title:" Article store failed", Detail: err.Error() }
	  return c.JSON(http.StatusUnprocessableEntity, ae )
  }

  log.Printf ( " POST article : %d \n", id)

return c.JSON(http.StatusCreated, article)
}

func postComment(c echo.Context) error {

	comment := new(Comment)
	if err := c.Bind(comment); err != nil {
    ae:=ApiError{ Code: "3", Title:" Comment bind failed", Detail: "Invalid bind" }
	  return c.JSON(http.StatusUnprocessableEntity, ae )
	}

  id,err:= store.GetNextId("comments")
  comment.Id = id
  comment.Created_at = getTS()
  comment.Updated_at = comment.Created_at
  buf, err := json.Marshal(comment)
  err= store.Create("comments",buf,id)
  err= store.Create("article:"+strconv.Itoa(comment.Article_id)+":comments",store.Itob(id),id)

  if err != nil {
    ae:=ApiError{ Code: "4", Title:"Comment store failed", Detail: err.Error() }
	  return c.JSON(http.StatusUnprocessableEntity, ae )
  }

  log.Printf ( " POST comment : %d \n", id)

return c.JSON(http.StatusCreated, comment)
}

func updateComment(c echo.Context) error {

  var comment Comment
	newComment := new(Comment)
  id,_ := strconv.Atoi(c.Param("comment_id"))

  log.Printf ( " PUT comment : %d \n", id)

	if err := c.Bind(newComment); err != nil {
    ae:=ApiError{ Code: "1", Title:"Bind failed", Detail: "Invalid bind" }
	  return c.JSON(http.StatusUnprocessableEntity, ae )
	}

  buf,_:= store.ReadUnique("comments",id)
  json.Unmarshal(buf,&comment)

  comment.Text = newComment.Text
  comment.Updated_at = getTS()

  b, _ := json.Marshal(comment)
  store.Create("comments",b,id)

return c.JSON(http.StatusCreated, comment)
}

func deleteArticle(c echo.Context) error {
  id,_ := strconv.Atoi(c.Param("article_id"))
  store.Remove("articles",id)
  store.Remove("articles:"+strconv.Itoa(id)+":comments",id)
  log.Printf ( " DELETE article : %d \n", id)
return c.String(http.StatusCreated, "DELETED")
}

func deleteComment(c echo.Context) error {
  var comment Comment
  id,_ := strconv.Atoi(c.Param("comment_id"))
  buf,_:= store.ReadUnique("comments",id)
  json.Unmarshal(buf,&comment)
  store.Remove("comments",id)
  store.Remove("articles:"+strconv.Itoa(comment.Article_id)+":comments",id)
  log.Printf ( " DELETE comment : %d \n", id)
return c.String(http.StatusCreated, "DELETED")
}

func getCommentById(c echo.Context) error {
  var comment Comment
  id,_ := strconv.Atoi(c.Param("comment_id"))
  log.Printf ( " GET comment : %d \n", id)
  buf,_:= store.ReadUnique("comments",id)
  json.Unmarshal(buf,&comment)
return c.JSON(http.StatusOK, comment)
}

func getCommentsByArticle(c echo.Context) error {
  var comment Comment
	commentList := new(CommentList)
  log.Println ( " GET comments, article : "+ c.Param("article_id"))
  buf,_:= store.ReadRangeFromSet("article:"+c.Param("article_id")+":comments","comments")
  for _,v := range buf {
    json.Unmarshal(v,&comment)
    commentList.Comments=append(commentList.Comments,comment)
  }
return c.JSON(http.StatusOK, commentList)
}

func getArticles(c echo.Context) error {
  var article Article
	articleList := new(ArticleList)
  log.Println (" GET articles ")
  buf,_:= store.ReadRange("articles")
  for _,v := range buf {
    json.Unmarshal(v,&article)
    articleList.Articles=append(articleList.Articles,article)
  }
return c.JSON(http.StatusOK, articleList)
}


func main() {
  const port = ":8080"

  store.Db = store.InitDB()
  defer store.Db.Close()

	e := echo.New()
  e.Use(middleware.CORS())

	e.POST("/api/articles", postArticle)
	e.POST("/api/comments", postComment)
  e.PUT("/api/comments/:comment_id", updateComment)
  e.GET("/api/articles", getArticles)
  e.GET("/api/articles/:article_id/comments", getCommentsByArticle)
  e.GET("/api/comments/:comment_id", getCommentById)
  e.DELETE("/api/comments/:comment_id", deleteComment)
  e.DELETE("/api/articles/:article_id", deleteArticle)

	e.Logger.Fatal(e.Start(port))
}

