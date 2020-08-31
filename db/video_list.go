package db

type VideoList struct {
	BasicModel

	Viewkey   string
	Image       string
	Author    string
	UpTime    string
	Title     string
	Vid       string
	Duration  string
	View      int
	Favorites int
	Comment   int
	Integral  int
	VideoUrl  string

	Timestamp
}

func (VideoList) TableName() string {
	return "video_list"
}
