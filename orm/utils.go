package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// extractColumnsAndValues 从结构体中提取列名和值
func (qb *queryBuilder) extractColumnsAndValues(data interface{}) ([]string, []interface{}) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	
	// 如果是指针，获取其指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		return nil, nil
	}
	
	var columns []string
	var values []interface{}
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		// 跳过未导出的字段
		if !fieldValue.CanInterface() {
			continue
		}
		
		// 获取字段标签
		tag := field.Tag.Get("orm")
		if tag == "-" {
			continue
		}
		
		columnName := field.Name
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				columnName = parts[0]
			}
		}
		
		// 转换为下划线命名
		columnName = camelToSnake(columnName)
		
		columns = append(columns, columnName)
		values = append(values, fieldValue.Interface())
	}
	
	return columns, values
}

// scanRows 扫描多行结果到切片
func scanRows(rows *sql.Rows, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest必须是指针类型")
	}
	
	destValue = destValue.Elem()
	if destValue.Kind() != reflect.Slice {
		return fmt.Errorf("dest必须是切片指针")
	}
	
	// 获取切片元素类型
	elemType := destValue.Type().Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}
	
	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	
	for rows.Next() {
		// 创建新的元素
		elem := reflect.New(elemType).Elem()
		
		// 准备扫描目标
		scanDest := make([]interface{}, len(columns))
		for i, col := range columns {
			field := findFieldByColumn(elem, col)
			if field.IsValid() && field.CanSet() {
				scanDest[i] = field.Addr().Interface()
			} else {
				var dummy interface{}
				scanDest[i] = &dummy
			}
		}
		
		// 扫描行数据
		if err := rows.Scan(scanDest...); err != nil {
			return err
		}
		
		// 添加到切片
		if isPtr {
			destValue.Set(reflect.Append(destValue, elem.Addr()))
		} else {
			destValue.Set(reflect.Append(destValue, elem))
		}
	}
	
	return rows.Err()
}

// scanRow 扫描单行结果到结构体
func scanRow(row *sql.Row, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest必须是指针类型")
	}
	
	destValue = destValue.Elem()
	if destValue.Kind() != reflect.Struct {
		return fmt.Errorf("dest必须是结构体指针")
	}
	
	// 这里需要更复杂的实现来处理单行扫描
	// 暂时返回一个简单的错误，实际使用中需要根据具体需求实现
	return fmt.Errorf("scanRow方法需要进一步实现")
}

// findFieldByColumn 根据列名查找结构体字段
func findFieldByColumn(structValue reflect.Value, columnName string) reflect.Value {
	structType := structValue.Type()
	
	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		
		// 检查orm标签
		tag := field.Tag.Get("orm")
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] == columnName {
				return fieldValue
			}
		}
		
		// 检查字段名转换
		if camelToSnake(field.Name) == columnName {
			return fieldValue
		}
		
		// 直接匹配字段名
		if strings.EqualFold(field.Name, columnName) {
			return fieldValue
		}
	}
	
	return reflect.Value{}
}

// getStructName 获取结构体名称
func getStructName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// parseFieldTag 解析字段标签
func parseFieldTag(tag string) FieldTag {
	fieldTag := FieldTag{}
	
	if tag == "" {
		return fieldTag
	}
	
	parts := strings.Split(tag, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		if i == 0 {
			// 第一部分是列名
			fieldTag.Column = part
			continue
		}
		
		// 解析其他属性
		switch part {
		case "primary":
			fieldTag.Primary = true
		case "auto_increment":
			fieldTag.AutoIncrement = true
		case "not_null":
			fieldTag.NotNull = true
		case "unique":
			fieldTag.Unique = true
		default:
			if strings.HasPrefix(part, "type:") {
				fieldTag.Type = strings.TrimPrefix(part, "type:")
			} else if strings.HasPrefix(part, "size:") {
				// 解析size，这里简化处理
				fieldTag.Size = 255 // 默认值
			} else if strings.HasPrefix(part, "default:") {
				fieldTag.Default = strings.TrimPrefix(part, "default:")
			} else if strings.HasPrefix(part, "comment:") {
				fieldTag.Comment = strings.TrimPrefix(part, "comment:")
			} else if strings.HasPrefix(part, "index:") {
				fieldTag.Index = strings.TrimPrefix(part, "index:")
			}
		}
	}
	
	return fieldTag
}

// convertValue 转换值类型
func convertValue(value interface{}, targetType reflect.Type) interface{} {
	if value == nil {
		return nil
	}
	
	sourceValue := reflect.ValueOf(value)
	if sourceValue.Type().ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType).Interface()
	}
	
	// 特殊类型转换
	switch targetType.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if str, ok := value.(string); ok {
			// 字符串转整数的逻辑
			_ = str
		}
	case reflect.Float32, reflect.Float64:
		if str, ok := value.(string); ok {
			// 字符串转浮点数的逻辑
			_ = str
		}
	case reflect.Bool:
		if str, ok := value.(string); ok {
			return str == "true" || str == "1" || str == "yes"
		}
	}
	
	// 时间类型特殊处理
	if targetType == reflect.TypeOf(time.Time{}) {
		if str, ok := value.(string); ok {
			if t, err := time.Parse("2006-01-02 15:04:05", str); err == nil {
				return t
			}
		}
	}
	
	return value
}

// isZeroValue 检查是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return v.Interface().(time.Time).IsZero()
		}
		// 对于其他结构体，检查所有字段是否都是零值
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
