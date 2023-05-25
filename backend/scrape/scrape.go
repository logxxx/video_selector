package scrape

import (
	"errors"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/log"
	"github.com/logxxx/utils/media"
	"github.com/logxxx/video_selector/proto"
	"os"
	"path/filepath"
	"sync"
)

var (
	allScrapeResult     = make(map[string][]*proto.VideoInfo, 0)
	allScrapeResultLock sync.Mutex
)

func GetAllScrapeResult(sourcePath string) []*proto.VideoInfo {
	allScrapeResultLock.Lock()
	defer allScrapeResultLock.Unlock()
	return allScrapeResult[sourcePath]
}

func GetVideoFromScrapeReuslt(sourcePath string, videoID int) *proto.VideoInfo {
	allScrapeResultLock.Lock()
	defer allScrapeResultLock.Unlock()
	videos, ok := allScrapeResult[sourcePath]
	if !ok {
		return nil
	}
	for _, video := range videos {
		if video.ID == videoID {
			return video
		}
	}
	return nil
}

func SetScrapeResult(sourcePath string, videos []*proto.VideoInfo) {
	allScrapeResultLock.Lock()
	defer allScrapeResultLock.Unlock()
	allScrapeResult[sourcePath] = videos
}

func RemoveVideoFromScrapeResult(sourcePath string, videoID int) {
	allScrapeResultLock.Lock()
	defer allScrapeResultLock.Unlock()
	videos, ok := allScrapeResult[sourcePath]
	if !ok {
		return
	}
	idx := -1
	for i, video := range videos {
		if video.ID == videoID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	if len(videos) == 1 {
		delete(allScrapeResult, sourcePath)
	}
	allScrapeResult[sourcePath] = append(videos[:idx], videos[idx+1:]...)
}

func HandleScrapeResult(req *proto.HandleScrapeResultReq) error {

	videoInfo := GetVideoFromScrapeReuslt(req.SourcePath, req.VideoID)
	if videoInfo == nil {
		return errors.New("video not found")
	}

	if !utils.HasFile(videoInfo.RealPath) {
		return errors.New("video file not exist")
	}

	var err error
	switch req.Action {
	case "DELETE":
		err = os.RemoveAll(videoInfo.RealPath)
	case "MOVE":
		err = fileutil.CopyFile(req.SourcePath, filepath.Join(req.DestDir, filepath.Base(videoInfo.RealPath)), 0777)
	default:
		err = errors.New("unknown action")
	}
	if err != nil {
		return err
	}

	RemoveVideoFromScrapeResult(req.SourcePath, videoInfo.ID)

	return nil

}

func GetScrapeResult(req *proto.GetScrapeResultReq) (*proto.ScrapeResult, error) {

	result := GetAllScrapeResult(req.SourcePath)
	if len(result) == 0 {
		return nil, nil
	}

	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 50
	}
	resp := &proto.ScrapeResult{Videos: result}
	if req.Token > 0 {
		for i := range resp.Videos {
			if resp.Videos[i].ID == req.Token {
				resp.Videos = resp.Videos[i:]
				break
			}
		}
	}
	if len(resp.Videos) > req.Limit {
		resp.NextToken = resp.Videos[req.Limit].ID
		resp.Videos = resp.Videos[:req.Limit]
	}
	return resp, nil
}

func DoScrape(req *proto.ScrapeReq) error {
	result, err := scrapeIter(req.SourcePath, req.IncludeChild)
	if err != nil {
		log.Errorf("DoScrape err:%v req:%+v", err, req)
		return err
	}
	for i, elem := range result {
		elem.ID = i + 1
		elem.SourcePath = req.SourcePath
	}
	SetScrapeResult(req.SourcePath, result)
	return nil
}

func scrapeIter(rootPath string, includeChild bool) ([]*proto.VideoInfo, error) {

	childs, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}
	videos := make([]*proto.VideoInfo, 0)
	childDirs := make([]string, 0)
	for _, child := range childs {
		if child.IsDir() {
			childDirs = append(childDirs, child.Name())
		} else {
			if !media.IsVideo(child.Name()) {
				continue
			}
			childInfo, err := child.Info()
			if err != nil {
				continue
			}
			videos = append(videos, &proto.VideoInfo{
				RealPath: filepath.Join(rootPath, child.Name()),
				Size:     childInfo.Size(),
				ModTime:  childInfo.ModTime().Unix(),
			})
		}
	}

	if includeChild {
		for _, childDir := range childDirs {
			childVideos, err := scrapeIter(filepath.Join(rootPath, childDir), includeChild)
			if err != nil {
				log.Errorf("DoScrape DoScrape err:%v", err)
			} else {
				videos = append(videos, childVideos...)
			}
		}
	}

	return videos, nil

}
