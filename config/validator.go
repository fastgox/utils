package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Validator 配置验证器
type Validator struct {
	config *Config
}

// NewValidator 创建配置验证器
func NewValidator(config *Config) *Validator {
	return &Validator{
		config: config,
	}
}

// Validate 验证当前配置
func (v *Validator) Validate() error {
	// 这里可以添加全局配置验证逻辑
	return nil
}

// ValidateStruct 验证结构体
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validateValue(reflect.ValueOf(s), "")
}

// validateValue 验证反射值
func (v *Validator) validateValue(val reflect.Value, path string) error {
	// 处理指针
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		return v.validateStruct(val, path)
	case reflect.Slice, reflect.Array:
		return v.validateSlice(val, path)
	case reflect.Map:
		return v.validateMap(val, path)
	default:
		return nil
	}
}

// validateStruct 验证结构体
func (v *Validator) validateStruct(val reflect.Value, path string) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过未导出的字段
		if !field.CanInterface() {
			continue
		}

		// 构建字段路径
		fieldPath := path
		if configTag := fieldType.Tag.Get("config"); configTag != "" && configTag != "-" {
			if fieldPath != "" {
				fieldPath += "." + configTag
			} else {
				fieldPath = configTag
			}
		} else {
			fieldName := strings.ToLower(fieldType.Name)
			if fieldPath != "" {
				fieldPath += "." + fieldName
			} else {
				fieldPath = fieldName
			}
		}

		// 验证字段
		if err := v.validateField(field, fieldType, fieldPath); err != nil {
			return err
		}

		// 递归验证嵌套结构
		if err := v.validateValue(field, fieldPath); err != nil {
			return err
		}
	}

	return nil
}

// validateField 验证单个字段
func (v *Validator) validateField(field reflect.Value, fieldType reflect.StructField, path string) error {
	validateTag := fieldType.Tag.Get("validate")
	if validateTag == "" || validateTag == "-" {
		return nil
	}

	// 解析验证标签
	rules := strings.Split(validateTag, ",")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if err := v.validateRule(field, rule, path); err != nil {
			return err
		}
	}

	return nil
}

// validateRule 验证单个规则
func (v *Validator) validateRule(field reflect.Value, rule, path string) error {
	parts := strings.SplitN(rule, "=", 2)
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		return v.validateRequired(field, path)
	case "min":
		return v.validateMin(field, ruleValue, path)
	case "max":
		return v.validateMax(field, ruleValue, path)
	case "len":
		return v.validateLen(field, ruleValue, path)
	case "email":
		return v.validateEmail(field, path)
	case "url":
		return v.validateURL(field, path)
	case "oneof":
		return v.validateOneOf(field, ruleValue, path)
	default:
		return fmt.Errorf("未知的验证规则: %s (字段: %s)", ruleName, path)
	}
}

// validateRequired 验证必填字段
func (v *Validator) validateRequired(field reflect.Value, path string) error {
	if v.isZeroValue(field) {
		return &ValidationError{
			Field:   path,
			Tag:     "required",
			Value:   field.Interface(),
			Message: fmt.Sprintf("字段 %s 是必填的", path),
		}
	}
	return nil
}

// validateMin 验证最小值
func (v *Validator) validateMin(field reflect.Value, ruleValue, path string) error {
	min, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		return fmt.Errorf("无效的min规则值: %s", ruleValue)
	}

	var fieldValue float64
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue = float64(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue = float64(field.Uint())
	case reflect.Float32, reflect.Float64:
		fieldValue = field.Float()
	case reflect.String:
		fieldValue = float64(len(field.String()))
	case reflect.Slice, reflect.Array, reflect.Map:
		fieldValue = float64(field.Len())
	default:
		return fmt.Errorf("min规则不支持类型: %s (字段: %s)", field.Kind(), path)
	}

	if fieldValue < min {
		return &ValidationError{
			Field:   path,
			Tag:     "min",
			Value:   field.Interface(),
			Message: fmt.Sprintf("字段 %s 的值 %v 小于最小值 %v", path, fieldValue, min),
		}
	}

	return nil
}

// validateMax 验证最大值
func (v *Validator) validateMax(field reflect.Value, ruleValue, path string) error {
	max, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		return fmt.Errorf("无效的max规则值: %s", ruleValue)
	}

	var fieldValue float64
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue = float64(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue = float64(field.Uint())
	case reflect.Float32, reflect.Float64:
		fieldValue = field.Float()
	case reflect.String:
		fieldValue = float64(len(field.String()))
	case reflect.Slice, reflect.Array, reflect.Map:
		fieldValue = float64(field.Len())
	default:
		return fmt.Errorf("max规则不支持类型: %s (字段: %s)", field.Kind(), path)
	}

	if fieldValue > max {
		return &ValidationError{
			Field:   path,
			Tag:     "max",
			Value:   field.Interface(),
			Message: fmt.Sprintf("字段 %s 的值 %v 大于最大值 %v", path, fieldValue, max),
		}
	}

	return nil
}

// validateLen 验证长度
func (v *Validator) validateLen(field reflect.Value, ruleValue, path string) error {
	expectedLen, err := strconv.Atoi(ruleValue)
	if err != nil {
		return fmt.Errorf("无效的len规则值: %s", ruleValue)
	}

	var actualLen int
	switch field.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		actualLen = field.Len()
	default:
		return fmt.Errorf("len规则不支持类型: %s (字段: %s)", field.Kind(), path)
	}

	if actualLen != expectedLen {
		return &ValidationError{
			Field:   path,
			Tag:     "len",
			Value:   field.Interface(),
			Message: fmt.Sprintf("字段 %s 的长度 %d 不等于期望长度 %d", path, actualLen, expectedLen),
		}
	}

	return nil
}

// validateEmail 验证邮箱格式
func (v *Validator) validateEmail(field reflect.Value, path string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("email规则只支持字符串类型 (字段: %s)", path)
	}

	email := field.String()
	if email == "" {
		return nil // 空值跳过验证，使用required规则验证必填
	}

	// 简单的邮箱格式验证
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return &ValidationError{
			Field:   path,
			Tag:     "email",
			Value:   email,
			Message: fmt.Sprintf("字段 %s 的值 %s 不是有效的邮箱格式", path, email),
		}
	}

	return nil
}

// validateURL 验证URL格式
func (v *Validator) validateURL(field reflect.Value, path string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("url规则只支持字符串类型 (字段: %s)", path)
	}

	url := field.String()
	if url == "" {
		return nil // 空值跳过验证
	}

	// 简单的URL格式验证
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return &ValidationError{
			Field:   path,
			Tag:     "url",
			Value:   url,
			Message: fmt.Sprintf("字段 %s 的值 %s 不是有效的URL格式", path, url),
		}
	}

	return nil
}

// validateOneOf 验证枚举值
func (v *Validator) validateOneOf(field reflect.Value, ruleValue, path string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("oneof规则只支持字符串类型 (字段: %s)", path)
	}

	value := field.String()
	if value == "" {
		return nil // 空值跳过验证
	}

	validValues := strings.Split(ruleValue, " ")
	for _, validValue := range validValues {
		if value == validValue {
			return nil
		}
	}

	return &ValidationError{
		Field:   path,
		Tag:     "oneof",
		Value:   value,
		Message: fmt.Sprintf("字段 %s 的值 %s 不在允许的值列表中: %s", path, value, ruleValue),
	}
}

// validateSlice 验证切片
func (v *Validator) validateSlice(val reflect.Value, path string) error {
	for i := 0; i < val.Len(); i++ {
		itemPath := fmt.Sprintf("%s[%d]", path, i)
		if err := v.validateValue(val.Index(i), itemPath); err != nil {
			return err
		}
	}
	return nil
}

// validateMap 验证映射
func (v *Validator) validateMap(val reflect.Value, path string) error {
	for _, key := range val.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		itemPath := fmt.Sprintf("%s[%s]", path, keyStr)
		if err := v.validateValue(val.MapIndex(key), itemPath); err != nil {
			return err
		}
	}
	return nil
}

// isZeroValue 检查是否为零值
func (v *Validator) isZeroValue(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Bool:
		return !field.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return field.IsNil()
	default:
		return false
	}
}
