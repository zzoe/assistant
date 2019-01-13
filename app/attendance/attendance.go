package attendance

import (
	"context"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/spf13/viper"
	"github.com/zzoe/assistant/format"
	"github.com/zzoe/assistant/util"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"strings"
)

type Attendance struct {
	mw       *walk.MainWindow
	tl       *walk.TextLabel
	de       *walk.DateEdit
	tv       *walk.TableView
	FilePath string
	rcModel  *RecordModel
}

func New() (a *Attendance) {
	a = new(Attendance)
	a.rcModel = new(RecordModel)
	a.rcModel.items = make([]*Record, 0)
	return
}

func (a *Attendance) onDropFiles() func(filePaths []string) {
	return func(filePaths []string) {
		if len(filePaths) > 0 {
			log.Info(strings.Join(filePaths, " "))
			a.FilePath = filePaths[0]

			if err := a.tl.SetText(a.FilePath); err != nil {
				log.Error("a.fp.SetText(a.FilePath)", zap.String("a.FilePath", a.FilePath), zap.Error(err))
			}

			viper.GetStringMapString("")
			if err := a.rcModel.ReadFromExcel(a.FilePath); err != nil {
				log.Error("读取考勤文件失败", zap.Error(err))
				return
			}
			a.refresh()
		}
	}
}

func (a *Attendance) onClicked() func() {
	return func() {
		if len(a.rcModel.items) < 1 {
			return
		}

		outPath := make(chan string, 1)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer close(outPath)

		go a.saveAs(ctx, outPath)

		dlg := new(walk.FileDialog)
		dlg.InitialDirPath = filepath.Dir(a.FilePath)
		dlg.Title = "另存为考勤表"
		dlg.Filter = "*.xlsx"

		accepted, err := dlg.ShowSave(a.mw)
		if err != nil {
			cancel()
			log.Error("选择导出文件失败", zap.Error(err))
			return
		}
		if !accepted {
			cancel()
			return
		}

		outFile := dlg.FilePath
		if filepath.Ext(outFile) != ".xlsx" {
			outFile += ".xlsx"
		}
		outPath <- outFile
	}
}

func (a *Attendance) refresh() {
	if len(a.rcModel.items) == 0 {
		return
	}

	length := len(a.rcModel.items[0].Times)
	a.clearCol(length + 2)
	a.resetCol(length)

	a.rcModel.PublishRowsReset()
}

func (a *Attendance) clearCol(length int) {
	cols := a.tv.Columns()
	if cols.Len() != length {
		if err := cols.Clear(); err != nil {
			log.Error("cols clear", zap.Error(err))
			return
		}
	}
}

func (a *Attendance) resetCol(length int) {
	a.addCol(&declarative.TableViewColumn{Title: "JobNum"})
	a.addCol(&declarative.TableViewColumn{Title: "Name"})

	for i := 1; i <= length; i++ {
		a.addCol(&declarative.TableViewColumn{Title: strconv.Itoa(i)})
	}
}

func (a *Attendance) addCol(tvc *declarative.TableViewColumn) {
	if err := tvc.Create(a.tv); err != nil {
		log.Error("添加列失败", zap.Error(err))
	}
}

func (a *Attendance) saveAs(ctx context.Context, outPath chan string) {
	xlsx, err := excelize.OpenFile(viper.GetString("excel.template"))
	if err != nil {
		log.Error("打开 excel 模板失败", zap.String("template", viper.GetString("excel.template")), zap.Error(err))
		return
	}

	date := a.de.Date()
	lastweekday := date.AddDate(0, 0, 1-date.Day()).Weekday() - 1

	sheetName := viper.GetString("excel.outsheet")
	for i, record := range a.rcModel.items {
		select {
		case <-ctx.Done():
			return
		default:
		}
		row := strconv.Itoa(i + 3)
		xlsx.SetCellValue(sheetName, "A"+row, record.JobNum)
		xlsx.SetCellValue(sheetName, "B"+row, record.Name)
		weekday := lastweekday
		for j, timeStamp := range record.Times {
			weekday = (weekday + 1) % 7
			col := ""
			if j > 23 {
				j -= 26
				col = "A"
			}
			col += string('C' + j)

			yellowFill, err := xlsx.NewStyle(format.YellowFill)
			util.Warn(err)
			redFont, err := xlsx.NewStyle(format.RedFont)
			util.Warn(err)
			normal, err := xlsx.NewStyle(format.Normal)
			util.Warn(err)

			times := strings.Fields(timeStamp)
			switch {
			case len(times) == 0:
			case len(times) == 1:
				xlsx.SetCellStyle(sheetName, col+row, col+row, yellowFill)
			case weekday != 6 && (times[0] > "08:30" || times[len(times)-1] < "17:30"),
				weekday == 6 && (times[0] > "09:30" || times[len(times)-1] < "17:00"):
				//log.Info(col+row, zap.Any("weekday", weekday), zap.String(times[0], times[len(times)-1]))
				xlsx.SetCellStyle(sheetName, col+row, col+row, redFont)
			default:
				xlsx.SetCellStyle(sheetName, col+row, col+row, normal)
			}

			xlsx.SetCellValue(sheetName, col+row, timeStamp)
		}
	}

	outFile, ok := <-outPath
	if ok && outFile != "" {
		//xlsx.SetActiveSheet(xlsx.GetSheetIndex(sheetName))
		if err = xlsx.SaveAs(outFile); err != nil {
			log.Error("另存考勤结果失败", zap.Error(err))
		}
	}
}
