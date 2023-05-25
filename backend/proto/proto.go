package proto

type Dir struct {
	RealPath string `json:"real_path"`
}

type VideoInfo struct {
	ID         int    `json:"id"`
	SourcePath string `json:"source_path"`
	RealPath   string `json:"real_path"`
	Size       int64  `json:"size"`
	Duration   int64  `json:"duration"`
	ModTime    int64  `json:"mod_time"`
}

type ScrapeReq struct {
	SourcePath   string `json:"root_path"`
	IncludeChild bool   `json:"include_child"`
}

type GetScrapeResultReq struct {
	SourcePath string `json:"root_path"`
	Token      int    `json:"token,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type ScrapeResult struct {
	Videos    []*VideoInfo `json:"videos"`
	NextToken int          `json:"next_token"`
}

type HandleScrapeResultReq struct {
	SourcePath string `json:"source_path"`
	VideoID    int    `json:"video_id"`
	Action     string `json:"action"`
	DestDir    string `json:"dest_dir"`
}
