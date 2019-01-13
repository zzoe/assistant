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
	JobNum int
	Name   string
	Times  []string
}

type RecordModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Record
}

func (m *RecordModel) RowCount() int {
	return len(m.items)
}

func (m *RecordModel) Value(row, col int) interface{} {
	if len(m.items) < row+1 {
		return nil
	}

	switch col {
	case 0:
		return m.items[row].JobNum
	case 1:
		return m.items[row].Name
	}

	if len(m.items[row].Times) < col-1 {
		return nil
	}

	return m.items[row].Times[col-2]
}

func (m *RecordModel) ReadFromExcel(filePath string) (err error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Error("读取 excel 文件失败", zap.Error(err))
		return
	}

	order := viper.GetStringMap("excel.order")
	idx := len(order)
	records := make([]*Record, idx)
	for i := 0; i < idx; i++ {
		records[i] = new(Record)
	}

	rows := xlsx.GetRows(viper.GetString("excel.insheet"))
	ok := false
	skip := false
	cur := 4
	var curObj interface{}
	for i, row := range rows {
		if i < 4 || skip || strings.TrimSpace(strings.Join(row, "")) == "" {
			skip = false
			continue
		}

		if i%2 == 0 {
			name := row[10]
			log.Info("[" + name + "]")
			if name == "" {
				skip = true
				continue
			}

			curObj, ok = order[name]
			if ok {
				cur = int(curObj.(int64)) - 1
			} else {
				cur = idx
				if len(records) <= idx {
					records = append(records, new(Record))
				}
			}

			records[cur].JobNum, err = strconv.Atoi(row[2])
			if err != nil {
				log.Error("工号有误", zap.Int("i", i), zap.String("row", strings.Join(row, " ")), zap.Error(err))
				return
			}
			records[cur].Name = name
			continue
		}

		records[cur].Times = make([]string, len(row))
		copy(records[cur].Times, row)
		if !ok {
			idx++
		}
	}

	m.items = records
	return
}
