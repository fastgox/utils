package orm

import (
	"fmt"
	"sort"
	"time"
)

// MigrationManager 迁移管理器
type MigrationManager struct {
	orm        *ORM
	migrations []Migration
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(orm *ORM) *MigrationManager {
	return &MigrationManager{
		orm:        orm,
		migrations: make([]Migration, 0),
	}
}

// AddMigration 添加迁移
func (mm *MigrationManager) AddMigration(migration Migration) {
	mm.migrations = append(mm.migrations, migration)
}

// Run 运行迁移
func (mm *MigrationManager) Run() error {
	// 确保迁移表存在
	if err := mm.ensureMigrationTable(); err != nil {
		return err
	}

	// 获取已执行的迁移
	executed, err := mm.getExecutedMigrations()
	if err != nil {
		return err
	}

	// 排序迁移
	sort.Slice(mm.migrations, func(i, j int) bool {
		return mm.migrations[i].Version() < mm.migrations[j].Version()
	})

	// 执行未执行的迁移
	for _, migration := range mm.migrations {
		version := migration.Version()
		if _, exists := executed[version]; !exists {
			fmt.Printf("运行迁移: %s\n", version)

			if err := migration.Up(); err != nil {
				return fmt.Errorf("迁移 %s 失败: %w", version, err)
			}

			if err := mm.recordMigration(version); err != nil {
				return fmt.Errorf("记录迁移 %s 失败: %w", version, err)
			}

			fmt.Printf("迁移 %s 完成\n", version)
		}
	}

	return nil
}

// Rollback 回滚迁移
func (mm *MigrationManager) Rollback(steps int) error {
	// 获取已执行的迁移
	executed, err := mm.getExecutedMigrations()
	if err != nil {
		return err
	}

	// 获取要回滚的迁移版本
	var versions []string
	for version := range executed {
		versions = append(versions, version)
	}

	// 排序（降序）
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	// 限制回滚步数
	if steps > len(versions) {
		steps = len(versions)
	}

	// 执行回滚
	for i := 0; i < steps; i++ {
		version := versions[i]

		// 找到对应的迁移
		var migration Migration
		for _, m := range mm.migrations {
			if m.Version() == version {
				migration = m
				break
			}
		}

		if migration == nil {
			return fmt.Errorf("找不到迁移: %s", version)
		}

		fmt.Printf("回滚迁移: %s\n", version)

		if err := migration.Down(); err != nil {
			return fmt.Errorf("回滚迁移 %s 失败: %w", version, err)
		}

		if err := mm.removeMigrationRecord(version); err != nil {
			return fmt.Errorf("删除迁移记录 %s 失败: %w", version, err)
		}

		fmt.Printf("迁移 %s 回滚完成\n", version)
	}

	return nil
}

// ensureMigrationTable 确保迁移表存在
func (mm *MigrationManager) ensureMigrationTable() error {
	dialect := NewDatabaseManager(mm.orm).GetDialect()

	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s VARCHAR(255) PRIMARY KEY,
			%s TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`,
		dialect.Quote("migrations"),
		dialect.Quote("version"),
		dialect.Quote("executed_at"),
	)

	_, err := mm.orm.Exec(sql)
	return err
}

// getExecutedMigrations 获取已执行的迁移
func (mm *MigrationManager) getExecutedMigrations() (map[string]time.Time, error) {
	executed := make(map[string]time.Time)

	rows, err := mm.orm.Query("SELECT version, executed_at FROM migrations")
	if err != nil {
		return executed, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		var executedAt time.Time

		if err := rows.Scan(&version, &executedAt); err != nil {
			return executed, err
		}

		executed[version] = executedAt
	}

	return executed, rows.Err()
}

// recordMigration 记录迁移
func (mm *MigrationManager) recordMigration(version string) error {
	_, err := mm.orm.Exec("INSERT INTO migrations (version) VALUES (?)", version)
	return err
}

// removeMigrationRecord 删除迁移记录
func (mm *MigrationManager) removeMigrationRecord(version string) error {
	_, err := mm.orm.Exec("DELETE FROM migrations WHERE version = ?", version)
	return err
}

// Status 获取迁移状态
func (mm *MigrationManager) Status() error {
	executed, err := mm.getExecutedMigrations()
	if err != nil {
		return err
	}

	// 排序迁移
	sort.Slice(mm.migrations, func(i, j int) bool {
		return mm.migrations[i].Version() < mm.migrations[j].Version()
	})

	fmt.Println("迁移状态:")
	fmt.Println("版本\t\t状态\t\t执行时间")
	fmt.Println("----------------------------------------")

	for _, migration := range mm.migrations {
		version := migration.Version()
		if executedAt, exists := executed[version]; exists {
			fmt.Printf("%s\t已执行\t\t%s\n", version, executedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%s\t待执行\t\t-\n", version)
		}
	}

	return nil
}

// BaseMigration 基础迁移结构
type BaseMigration struct {
	version string
	orm     *ORM
}

// NewBaseMigration 创建基础迁移
func NewBaseMigration(version string, orm *ORM) *BaseMigration {
	return &BaseMigration{
		version: version,
		orm:     orm,
	}
}

// Version 获取版本
func (bm *BaseMigration) Version() string {
	return bm.version
}

// Up 默认上迁移实现
func (bm *BaseMigration) Up() error {
	return fmt.Errorf("Up方法未实现")
}

// Down 默认下迁移实现
func (bm *BaseMigration) Down() error {
	return fmt.Errorf("Down方法未实现")
}

// Exec 执行SQL
func (bm *BaseMigration) Exec(sql string, args ...interface{}) error {
	_, err := bm.orm.Exec(sql, args...)
	return err
}

// CreateTable 创建表
func (bm *BaseMigration) CreateTable(tableName string, callback func(TableInterface)) error {
	schema := NewSchema(bm.orm)
	return schema.CreateTable(tableName, callback)
}

// DropTable 删除表
func (bm *BaseMigration) DropTable(tableName string) error {
	schema := NewSchema(bm.orm)
	return schema.DropTable(tableName)
}

// AlterTable 修改表
func (bm *BaseMigration) AlterTable(tableName string, callback func(TableInterface)) error {
	schema := NewSchema(bm.orm)
	return schema.AlterTable(tableName, callback)
}

// AddColumn 添加列
func (bm *BaseMigration) AddColumn(tableName, columnName, columnType string) error {
	dialect := NewDatabaseManager(bm.orm).GetDialect()
	definition := ColumnDefinition{
		Name: columnName,
		Type: columnType,
	}
	sql := dialect.AddColumnSQL(tableName, columnName, definition)
	return bm.Exec(sql)
}

// DropColumn 删除列
func (bm *BaseMigration) DropColumn(tableName, columnName string) error {
	dialect := NewDatabaseManager(bm.orm).GetDialect()
	sql := dialect.DropColumnSQL(tableName, columnName)
	return bm.Exec(sql)
}

// AddIndex 添加索引
func (bm *BaseMigration) AddIndex(tableName, indexName string, columns []string, unique bool) error {
	dialect := NewDatabaseManager(bm.orm).GetDialect()
	sql := dialect.CreateIndexSQL(tableName, indexName, columns, unique)
	return bm.Exec(sql)
}

// DropIndex 删除索引
func (bm *BaseMigration) DropIndex(tableName, indexName string) error {
	dialect := NewDatabaseManager(bm.orm).GetDialect()
	sql := dialect.DropIndexSQL(tableName, indexName)
	return bm.Exec(sql)
}

// 全局便捷方法

// Migrate 运行迁移
func Migrate(migrations ...Migration) error {
	mm := NewMigrationManager(GetGlobalORM())
	for _, migration := range migrations {
		mm.AddMigration(migration)
	}
	return mm.Run()
}

// RollbackMigration 回滚迁移
func RollbackMigration(steps int) error {
	mm := NewMigrationManager(GetGlobalORM())
	return mm.Rollback(steps)
}

// MigrationStatus 获取迁移状态
func MigrationStatus() error {
	mm := NewMigrationManager(GetGlobalORM())
	return mm.Status()
}
