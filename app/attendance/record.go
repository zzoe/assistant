package attendance

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lxn/walk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Record struct {
	JobNum   int
	Name     string
	Times    []string
}

type RecordModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items []*Record
}

func (m *RecordModel) RowCount() int{
	return len(m.items)
}

func (m *RecordModel)Value(row, col int) interface{}{
	if len(m.items) < row +1 {
		return nil
	}

	switch col{
	case 0:
		return m.items[row].JobNum
	case 1:
		return m.items[row].Name
	}

	if len(m.items[row].Times) < col-1{
		return nil
	}

	return m.items[row].Times[col -2]
}

func (m *RecordModel)ReadFromExcel(filePath string) (err error){
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Error("读取 excel 文件失败", zap.Error(err))
		return
	}

	records := make([]*Record, 0)
	idx := 0
	rows := xlsx.GetRows(viper.GetString("excel.insheet"))
	for i, row := range rows {
		if i< 4{
			continue
		}

		if i%2 == 0{
			if len(records) < idx + 1{
				records = append(records, new(Record))
			}
			records[idx].JobNum, err = strconv.Atoi(row[2])
			if err != nil{
				log.Error("工号有误", zap.Int("i", i), zap.String("row", strings.Join(row, " ")), zap.Error(err))
				return
			}
			records[idx].Name = row[10]
			continue
		}

		if strings.TrimSpace(strings.Join(row, "")) == ""{
			continue
		}

		records[idx].Times = make([]string, len(row))
		copy(records[idx].Times, row)
		idx++
	}

	m.items = records
	return
}