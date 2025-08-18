package orm

import (
	"fmt"
	"reflect"
	"time"
)

// 移除BaseModel，让用户自己定义模型结构

// ModelManager 模型管理器
type ModelManager struct {
	orm *ORM
}

// NewModelManager 创建模型管理器
func NewModelManager(orm *ORM) *ModelManager {
	return &ModelManager{orm: orm}
}

// GetTableInfo 获取表信息
func (mm *ModelManager) GetTableInfo(model interface{}) *TableInfo {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	tableName := mm.getTableName(model)
	columns := mm.getColumns(t)

	return &TableInfo{
		Name:    tableName,
		Columns: columns,
		Model:   model,
	}
}

// getTableName 获取表名
func (mm *ModelManager) getTableName(model interface{}) string {
	if m, ok := model.(ModelInterface); ok {
		tableName := m.TableName()
		if tableName != "" {
			return tableName
		}
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return camelToSnake(t.Name())
}

// getColumns 获取列信息
func (mm *ModelManager) getColumns(t reflect.Type) []ColumnInfo {
	var columns []ColumnInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}

		// 解析标签
		tag := field.Tag.Get("orm")
		if tag == "-" {
			continue
		}

		fieldTag := parseFieldTag(tag)
		if fieldTag.Column == "" {
			fieldTag.Column = camelToSnake(field.Name)
		}

		column := ColumnInfo{
			Name:          fieldTag.Column,
			GoName:        field.Name,
			Type:          mm.getColumnType(field.Type, fieldTag),
			GoType:        field.Type,
			Primary:       fieldTag.Primary,
			AutoIncrement: fieldTag.AutoIncrement,
			NotNull:       fieldTag.NotNull,
			Unique:        fieldTag.Unique,
			Default:       fieldTag.Default,
			Comment:       fieldTag.Comment,
			Index:         fieldTag.Index,
			Size:          fieldTag.Size,
		}

		columns = append(columns, column)
	}

	return columns
}

// getColumnType 获取列类型
func (mm *ModelManager) getColumnType(goType reflect.Type, tag FieldTag) string {
	// 如果标签中指定了类型，使用标签中的类型
	if tag.Type != "" {
		return tag.Type
	}

	// 获取数据库方言
	dialect := NewDatabaseManager(mm.orm).GetDialect()

	return dialect.DataType(goType, tag.Size)
}

// CreateTable 创建表
func (mm *ModelManager) CreateTable(model interface{}) error {
	tableInfo := mm.GetTableInfo(model)
	if tableInfo == nil {
		return fmt.Errorf("无法获取表信息")
	}

	// 构建列定义
	var columnDefs []ColumnDefinition
	for _, col := range tableInfo.Columns {
		colDef := ColumnDefinition{
			Name:          col.Name,
			Type:          col.Type,
			Size:          col.Size,
			NotNull:       col.NotNull,
			Primary:       col.Primary,
			AutoIncrement: col.AutoIncrement,
			Unique:        col.Unique,
			Default:       col.Default,
			Comment:       col.Comment,
		}
		columnDefs = append(columnDefs, colDef)
	}

	// 获取方言并生成SQL
	dialect := NewDatabaseManager(mm.orm).GetDialect()
	sql := dialect.CreateTableSQL(tableInfo.Name, columnDefs)

	// 执行SQL
	_, err := mm.orm.Exec(sql)
	return err
}

// DropTable 删除表
func (mm *ModelManager) DropTable(model interface{}) error {
	tableName := mm.getTableName(model)
	dialect := NewDatabaseManager(mm.orm).GetDialect()
	sql := dialect.DropTableSQL(tableName)

	_, err := mm.orm.Exec(sql)
	return err
}

// HasTable 检查表是否存在
func (mm *ModelManager) HasTable(model interface{}) (bool, error) {
	tableName := mm.getTableName(model)

	var sql string
	switch mm.orm.config.Type {
	case MySQL:
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	case PostgreSQL:
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	case SQLite:
		sql = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?"
	case SQLServer:
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	default:
		return false, fmt.Errorf("不支持的数据库类型")
	}

	var count int
	err := mm.orm.QueryRow(sql, tableName).Scan(&count)
	return count > 0, err
}

// AutoMigrate 自动迁移
func (mm *ModelManager) AutoMigrate(models ...interface{}) error {
	for _, model := range models {
		exists, err := mm.HasTable(model)
		if err != nil {
			return err
		}

		if !exists {
			if err := mm.CreateTable(model); err != nil {
				return err
			}
		} else {
			// TODO: 实现表结构更新逻辑
		}
	}

	return nil
}

// TableInfo 表信息
type TableInfo struct {
	Name    string       `json:"name"`
	Columns []ColumnInfo `json:"columns"`
	Model   interface{}  `json:"-"`
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name          string       `json:"name"`
	GoName        string       `json:"go_name"`
	Type          string       `json:"type"`
	GoType        reflect.Type `json:"-"`
	Primary       bool         `json:"primary"`
	AutoIncrement bool         `json:"auto_increment"`
	NotNull       bool         `json:"not_null"`
	Unique        bool         `json:"unique"`
	Default       interface{}  `json:"default"`
	Comment       string       `json:"comment"`
	Index         string       `json:"index"`
	Size          int          `json:"size"`
}

// GetPrimaryKey 获取主键列
func (ti *TableInfo) GetPrimaryKey() *ColumnInfo {
	for _, col := range ti.Columns {
		if col.Primary {
			return &col
		}
	}
	return nil
}

// GetColumnByName 根据名称获取列
func (ti *TableInfo) GetColumnByName(name string) *ColumnInfo {
	for _, col := range ti.Columns {
		if col.Name == name || col.GoName == name {
			return &col
		}
	}
	return nil
}

// GetColumnNames 获取所有列名
func (ti *TableInfo) GetColumnNames() []string {
	names := make([]string, len(ti.Columns))
	for i, col := range ti.Columns {
		names[i] = col.Name
	}
	return names
}

// ValidateModel 验证模型
func (mm *ModelManager) ValidateModel(model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("模型必须是结构体类型")
	}

	tableInfo := mm.GetTableInfo(model)
	if tableInfo == nil {
		return fmt.Errorf("无法获取表信息")
	}

	// 验证必填字段
	for _, col := range tableInfo.Columns {
		if col.NotNull && !col.AutoIncrement {
			field := v.FieldByName(col.GoName)
			if field.IsValid() && isZeroValue(field) {
				return fmt.Errorf("字段 %s 不能为空", col.GoName)
			}
		}
	}

	return nil
}

// SetTimestamps 设置时间戳
func (mm *ModelManager) SetTimestamps(model interface{}, isUpdate bool) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	now := time.Now()

	// 设置UpdatedAt
	if field := v.FieldByName("UpdatedAt"); field.IsValid() && field.CanSet() {
		if field.Type() == reflect.TypeOf(time.Time{}) {
			field.Set(reflect.ValueOf(now))
		}
	}

	// 如果是新建记录，设置CreatedAt
	if !isUpdate {
		if field := v.FieldByName("CreatedAt"); field.IsValid() && field.CanSet() {
			if field.Type() == reflect.TypeOf(time.Time{}) && field.Interface().(time.Time).IsZero() {
				field.Set(reflect.ValueOf(now))
			}
		}
	}
}

// 全局便捷方法

// AutoMigrate 自动迁移
func AutoMigrate(models ...interface{}) error {
	mm := NewModelManager(GetGlobalORM())
	return mm.AutoMigrate(models...)
}

// CreateTable 创建表
func CreateTable(model interface{}) error {
	mm := NewModelManager(GetGlobalORM())
	return mm.CreateTable(model)
}

// DropTable 删除表
func DropTable(model interface{}) error {
	mm := NewModelManager(GetGlobalORM())
	return mm.DropTable(model)
}

// HasTable 检查表是否存在
func HasTable(model interface{}) (bool, error) {
	mm := NewModelManager(GetGlobalORM())
	return mm.HasTable(model)
}
