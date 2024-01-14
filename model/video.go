package model

// FrameInfo 帧信息和base64编码
type FrameInfo struct {
	SegmentIndex int    `json:"segment_index"`
	FrameIndex   int    `json:"frame_index"`
	Base64Data   string `json:"base64_data"`
}

// VideoInfo 视频的基本信息
type VideoInfo struct {
	Format struct {
		Duration string `json:"duration"` // 视频时长
	} `json:"format"`
}

type DYVideoInfo struct {
	ItemList []struct {
		Video struct {
			PlayAddr struct {
				Uri string `json:"uri"`
			} `json:"play_addr"`
		} `json:"video"`
		Desc string `json:"desc"`
	} `json:"item_list"`
}
