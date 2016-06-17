package crawler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

const(
	Pretty		  = "pretty"
	StartSuccess  = "Start Successfully"
	StopSuccess   = "Stop Successfully"
	AuthFailed    = "You do not have the authorization!"
	AddSuccess    = "Add Successfully"
	GetStatFailed = "Fail to get Statistic Info"
)

type HttpServer struct {
	node *Node
	router *gin.Engine
}

func NewHttpServer(node *Node) *HttpServer {
	res := &HttpServer{
		node : node,
		router: gin.Default(),
	}
	return res
}

func (h *HttpServer)Serve() {
	h.Register()
	err := http.ListenAndServe(fmt.Sprintf(":%d", Configs.HttpPort), h.router)
	if err != nil {
		Error.Println(err)
	}
}

func (h *HttpServer)Register(){
	h.router.GET("/", GetProfile)
	h.router.GET("/help", Help)
	h.router.GET("/start", StartCrawler)
	h.router.GET("/stop", StopCrawler)
	h.router.GET("/restart", RestartCrawler)
	h.router.POST("/add", AddUrlToCrawl)
}

func GetProfile(c *gin.Context) {
	stat, err := GetNodeInstance().GetStatistic()
	if err != nil {
		Error.Println(err)
		Response(c, GetStatFailed)
		return
	}
	Response(c, stat)
}

//
func AddUrlToCrawl(c *gin.Context) {
	//todo
	Response(c, AddSuccess)
}
// verify the operator's Authorize
func Authorize(c *gin.Context) bool {
	if Configs.Auth == false || Configs.UserName == "" || Configs.Token == ""{
		return true
	}
	userName, ok := c.GetQuery("UserName")
	if !ok || userName == "" {
		return false
	}
	token, ok := c.GetQuery("Token")
	if !ok || token == "" {
		return false
	}
	if userName == Configs.UserName && token == Configs.Token {
		return true
	}
	return false
}

func StopCrawler(c *gin.Context) {
	if !Authorize(c){
		Response(c, AuthFailed)
		return
	}
	GetNodeInstance().crawl.Stop()
	Response(c, StopSuccess)
}

func StartCrawler(c *gin.Context) {
	if !Authorize(c){
		Response(c, AuthFailed)
		return
	}
	GetNodeInstance().crawl.Start()
	Response(c, StartSuccess)
}

func RestartCrawler(c *gin.Context) {
	if !Authorize(c){
		Response(c, AuthFailed)
		return
	}
	GetNodeInstance().crawl.Restart()
	Response(c, StartSuccess)
}

func Help(c *gin.Context) {

	usage := map[string]interface{}{
		"Usage" : "",
	}

	Response(c, usage)
}

// response the json if url's query string contains pretty param
func Response(c *gin.Context, v interface{}) {
	if _, ok := c.GetQuery(Pretty); ok {
		// response the pretty json
		c.IndentedJSON(http.StatusOK, v)
		return
	}
	c.JSON(http.StatusOK, v)
}
