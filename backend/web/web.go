package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/reqresp"
	"github.com/logxxx/utils/runutil"
	"github.com/logxxx/video_selector/base"
	"github.com/logxxx/video_selector/proto"
	"github.com/logxxx/video_selector/scrape"
	"os"
	"path/filepath"
)

func InitWeb() {
	runutil.GoRunSafe(func() {
		g := gin.Default()
		g.Use(reqresp.Cors())
		g.GET("/", handleHome)
		g.HEAD("/", handleHome)
		g.GET("/dist/*filepath", handleDist)
		g.HEAD("/dist/*filepath", handleDist)
		g.GET("/dirs", func(c *gin.Context) {
			parentPath := c.Query("parent_path")
			resp, err := getChildDirs(parentPath)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			reqresp.MakeResp(c, resp)
		})
		g.POST("/do_scrape", func(c *gin.Context) {
			req := &proto.ScrapeReq{}
			err := reqresp.ParseReq(c, req)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			err = scrape.DoScrape(req)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}

			reqresp.MakeRespOk(c)

		})
		g.POST("/get_scrape_result", func(c *gin.Context) {
			req := &proto.GetScrapeResultReq{}
			err := reqresp.ParseReq(c, req)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			resp, err := scrape.GetScrapeResult(req)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			reqresp.MakeResp(c, resp)
		})
		g.POST("/handle_scrape_result", func(c *gin.Context) {
			req := &proto.HandleScrapeResultReq{}
			err := reqresp.ParseReq(c, req)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			err = scrape.HandleScrapeResult(req)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			reqresp.MakeRespOk(c)
		})
		err := g.Run(fmt.Sprintf("0.0.0.0:%v", *base.Port))
		if err != nil {
			panic(err)
		}
	})

}

func getChildDirs(dir string) ([]proto.Dir, error) {

	if dir == "" {
		dir = "/"
	}

	childs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	resp := make([]proto.Dir, 0)
	for _, child := range childs {
		resp = append(resp, proto.Dir{RealPath: filepath.Join(dir, child.Name())})
	}

	return resp, nil
}
