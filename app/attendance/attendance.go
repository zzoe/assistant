package attendance

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Attendance struct {
	mw *walk.MainWindow
	tl *walk.TextLabel
	tv *walk.TableView
	FilePath string
	rcModel *RecordModel
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

			//a.rcModel.items = []*Record{
			//	&Record{JobNum:1, Name:"张三", Times: []string{"08:00", "18:00","08:00", "18:00"}},
			//	&Record{JobNum:2, Name:"李四", Times: []string{"08:01", "18:01","08:00", "18:00"}},
			//	&Record{JobNum:3, Name:"王五", Times: []string{"08:02", "18:02","08:00", "18:00"}},
			//}
			if err := a.rcModel.ReadFromExcel(a.FilePath); err != nil{
				log.Error("读取考勤文件失败", zap.Error(err))
				return
			}
			a.refresh()
		}
	}
}

func (a *Attendance) onClicked()func(){
	return func() {
		if len(a.rcModel.items) < 1{
			return
		}

		xlsx, err := excelize.OpenFile(viper.GetString("excel.template"))
		if err != nil{
			log.Error("打开 excel 模板失败",zap.String("template", viper.GetString("excel.template")), zap.Error(err))
			return
		}

		sheetName := viper.GetString("excel.outsheet")
		for i, record := range a.rcModel.items{
			xlsx.SetCellValue(sheetName, "A" + strconv.Itoa(i+3), record.JobNum)
			xlsx.SetCellValue(sheetName, "B" + strconv.Itoa(i+3), record.Name)
			for j,timeStamp := range record.Times{
				col := ""
				if j > 23{
					j -= 26
					col = "A"
				}
				col += string('C' + j)
				xlsx.SetCellValue(sheetName, col+ strconv.Itoa(i+3), timeStamp)
			}
		}

		//xlsx.SetActiveSheet(xlsx.GetSheetIndex(sheetName))
		if err = xlsx.SaveAs(viper.GetString("excel.outpath")); err != nil{
			log.Error("另存考勤结果失败", zap.Error(err))
		}
	}
}

func (a *Attendance)refresh(){
	if len(a.rcModel.items) == 0{
		return
	}

	length := len(a.rcModel.items[0].Times)
	a.clearCol(length+2)
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

func (a *Attendance) resetCol(length int){
	a.addCol( &declarative.TableViewColumn{Title: "JobNum"})
	a.addCol( &declarative.TableViewColumn{Title: "Name"})

	for i:=1; i<= length; i++{
		a.addCol( &declarative.TableViewColumn{Title: strconv.Itoa(i)})
	}
}

func (a *Attendance) addCol(tvc *declarative.TableViewColumn){
	if err := tvc.Create(a.tv); err != nil{
		log.Error("添加列失败", zap.Error(err))
	}
}
