package attendance

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/zzoe/assistant/cfg"
	"go.uber.org/zap"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var (
	log = cfg.Log
)

func Run() (int, error) {
	a := New()
	mw := create(a)

	return mw.Run()
}

func icon() (icon *walk.Icon) {
	iconFile, err := os.Open("vim-go.png")
	if err != nil {
		log.Warn("icon 图片解析失败", zap.Error(err))
		return
	}

	iconImage, _, err := image.Decode(iconFile)
	if err != nil {
		log.Warn("icon 图片解析失败", zap.Error(err))
		return
	}

	icon, err = walk.NewIconFromImage(iconImage)
	if err != nil {
		log.Error("NewIconFromFile", zap.Error(err))
	}

	return
}

func create(a *Attendance) *declarative.MainWindow {
	return &declarative.MainWindow{
		Title:    "Attendance statistics",
		Icon:     icon(),
		Size:     declarative.Size{Width: 800, Height: 600},
		AssignTo: &a.mw,

		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.Label{
						Text: "文件路径:",
					},
					declarative.TextLabel{
						AssignTo: &a.tl,
						Text:     a.FilePath,
					},
					declarative.PushButton{
						Text:      "导出考勤",
						OnClicked: a.onClicked(),
					},
				},
			},
			declarative.TableView{
				AssignTo:              &a.tv,
				Model:                 a.rcModel,
				AlternatingRowBGColor: walk.RGB(239, 239, 239),
			},
		},

		OnDropFiles: a.onDropFiles(),
	}
}
