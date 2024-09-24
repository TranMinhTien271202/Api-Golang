package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Method struct {
	DB    *sql.DB
	Table string
}

// NewMethod khởi tạo một instance mới của Method với DB và Table
func NewMethod(db *sql.DB, table string) *Method {
	return &Method{
		DB:    db,
		Table: table,
	}
}

// GetAllRecord
func (m *Method) GetAllRecord(columns []string, whereClause string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	columnsStr := ""
	for i, column := range columns {
		if i > 0 {
			columnsStr += ", "
		}
		columnsStr += column
	}
	query := fmt.Sprintf("SELECT %s FROM %s", columnsStr, m.Table)
	if whereClause != "" {
		query += " " + whereClause
	}
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		result := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				result[col] = nil
				continue
			}
			switch v := val.(type) {
			case []byte:
				result[col] = string(v) // chuyển đổi []byte sang string
			case int, int8, int16, int32, int64:
				result[col] = v // giữ nguyên kiểu int
			case float32, float64:
				result[col] = v // giữ nguyên kiểu float
			default:
				result[col] = val // giữ nguyên nếu không cần xử lý đặc biệt
			}
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// GetRecordByID lấy thông tin bản ghi theo ID
func (m *Method) GetRecordByID(id int) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	var idVal int
	var nameVal string
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id = ?", m.Table)
	row := m.DB.QueryRow(query, id)
	err := row.Scan(&idVal, &nameVal)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Record not found")
		}
		return nil, err
	}
	data["id"] = idVal
	data["name"] = nameVal
	return data, nil
}

// InsertRecord thêm bản ghi mới với trường và giá trị động
func (m *Method) InsertRecord(fields []string, values []interface{}) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", m.Table,
		strings.Join(fields, ", "),
		strings.Repeat("?,", len(values)-1)+"?")
	_, err := m.DB.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRecord cập nhật bản ghi với trường và giá trị động
func (m *Method) UpdateRecord(id int, fields []string, values []interface{}) error {
	sets := make([]string, len(fields))
	for i, field := range fields {
		sets[i] = fmt.Sprintf("%s = ?", field)
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.Table, strings.Join(sets, ", "))
	values = append(values, id) // Thêm ID vào cuối để phù hợp với query
	_, err := m.DB.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRecord xóa bản ghi theo ID
func (m *Method) DeleteRecord(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.Table)
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
