package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hantdev/chorus-controller/internal/domain"
)

func main() {
	// Export schema to SQL file
	schemaFile := "schema.sql"
	if len(os.Args) > 1 {
		schemaFile = os.Args[1]
	}

	// Write schema to file
	file, err := os.Create(schemaFile)
	if err != nil {
		log.Fatal("Failed to create schema file:", err)
	}
	defer file.Close()

	// Generate PostgreSQL schema from GORM models
	generateSchemaFromModels(file)

	fmt.Printf("Schema exported to %s\n", schemaFile)
}

func generateSchemaFromModels(file *os.File) {
	// Write header
	file.WriteString(`-- Auto-generated schema from GORM models
-- Generated at: ` + time.Now().Format("2006-01-02 15:04:05") + `

-- CreateExtension creates the "uuid-ossp" extension if it does not exist.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

`)

	// Define all models
	models := []interface{}{
		&domain.Storage{},
		&domain.ReplicateJob{},
		&domain.TokenInfo{},
	}

	// Generate schema for each model
	for _, model := range models {
		generateTableSchema(file, model)
	}
}

func generateTableSchema(file *os.File, model interface{}) {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

	// Get table name from GORM tag or use struct name
	tableName := getTableName(modelType)

	file.WriteString(fmt.Sprintf(`-- CreateTable creates the "%s" table.
CREATE TABLE "%s" (
`, tableName, tableName))

	// Generate columns
	var columns []string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		columnDef := generateColumnDefinition(field, modelValue.Field(i))
		if columnDef != "" {
			columns = append(columns, "    "+columnDef)
		}
	}

	file.WriteString(strings.Join(columns, ",\n"))
	file.WriteString(fmt.Sprintf(`,

    CONSTRAINT "%s_pkey" PRIMARY KEY ("id")
);

`, tableName))

	// Generate indexes
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		generateIndexes(file, field, tableName)
	}

	file.WriteString("\n")
}

func getTableName(modelType reflect.Type) string {
	// Default to snake_case of struct name
	return toSnakeCase(modelType.Name())
}

func generateColumnDefinition(field reflect.StructField, fieldValue reflect.Value) string {
	gormTag := field.Tag.Get("gorm")

	// Skip fields with "-" tag
	if strings.Contains(gormTag, "-") {
		return ""
	}

	// Get column name
	columnName := getColumnName(field, gormTag)
	if columnName == "" {
		return ""
	}

	var parts []string
	parts = append(parts, fmt.Sprintf(`"%s"`, columnName))

	// Data type
	dataType := getPostgreSQLType(field, gormTag)
	parts = append(parts, dataType)

	// NOT NULL - check if field is pointer (nullable) or has explicit NOT NULL
	isPointer := field.Type.Kind() == reflect.Ptr
	hasNotNull := strings.Contains(gormTag, "not null") || strings.Contains(gormTag, "NOT NULL")

	if !isPointer && !hasNotNull {
		// Non-pointer fields are NOT NULL by default
		parts = append(parts, "NOT NULL")
	} else if hasNotNull {
		parts = append(parts, "NOT NULL")
	}

	// DEFAULT value
	if strings.Contains(gormTag, "default:") {
		defaultValue := extractTagValue(gormTag, "default:")
		// Keep gen_random_uuid() for now
		if defaultValue == "gen_random_uuid()" {
			defaultValue = "uuid_generate_v4()"
		}
		parts = append(parts, fmt.Sprintf("DEFAULT %s", defaultValue))
	} else if field.Type == reflect.TypeOf(uuid.UUID{}) {
		parts = append(parts, "DEFAULT uuid_generate_v4()")
	} else if field.Type.Kind() == reflect.Struct && field.Type.Name() == "UUID" {
		parts = append(parts, "DEFAULT uuid_generate_v4()")
	} else if field.Type == reflect.TypeOf(time.Time{}) {
		if strings.Contains(gormTag, "autoCreateTime") {
			parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
		}
	}

	return strings.Join(parts, " ")
}

func getColumnName(field reflect.StructField, gormTag string) string {
	// Check for column name in GORM tag
	if strings.Contains(gormTag, "column:") {
		return extractTagValue(gormTag, "column:")
	}

	// Special handling for common field names
	if field.Name == "ID" {
		return "id"
	}

	// Special handling for acronyms
	if field.Name == "RPM" {
		return "rpm"
	}

	// Handle fields containing RPM
	if strings.Contains(field.Name, "RPM") {
		return strings.ReplaceAll(toSnakeCase(field.Name), "_r_p_m", "_rpm")
	}

	// Handle fields ending with ID
	if strings.HasSuffix(field.Name, "ID") {
		baseName := strings.TrimSuffix(field.Name, "ID")
		return toSnakeCase(baseName) + "_id"
	}

	// Default to snake_case of field name
	return toSnakeCase(field.Name)
}

func getPostgreSQLType(field reflect.StructField, gormTag string) string {
	// Check for explicit type in GORM tag
	if strings.Contains(gormTag, "type:") {
		typeValue := extractTagValue(gormTag, "type:")
		// Replace uuid with UUID for PostgreSQL compatibility
		if typeValue == "uuid" {
			typeValue = "UUID"
		}
		return typeValue
	}

	// Map Go types to PostgreSQL types
	switch field.Type {
	case reflect.TypeOf(uuid.UUID{}):
		return "UUID"
	case reflect.TypeOf(time.Time{}):
		if strings.Contains(gormTag, "type:timestamp") {
			return "TIMESTAMP(3)"
		}
		return "TIMESTAMP WITH TIME ZONE"
	case reflect.TypeOf((*time.Time)(nil)):
		return "TIMESTAMP WITH TIME ZONE"
	}

	switch field.Type.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "INTEGER"
	case reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.String:
		// Check for size constraint
		if strings.Contains(gormTag, "size:") {
			size := extractTagValue(gormTag, "size:")
			return fmt.Sprintf("VARCHAR(%s)", size)
		}
		return "TEXT"
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			return "BYTEA"
		}
		return "TEXT"
	default:
		return "TEXT"
	}
}

func generateIndexes(file *os.File, field reflect.StructField, tableName string) {
	gormTag := field.Tag.Get("gorm")
	if gormTag == "" {
		return
	}

	columnName := getColumnName(field, gormTag)
	if columnName == "" {
		return
	}

	tags := strings.Split(gormTag, ";")
	for _, tag := range tags {
		if strings.Contains(tag, "uniqueIndex") {
			indexName := fmt.Sprintf("idx_%s_%s", tableName, columnName)
			file.WriteString(fmt.Sprintf(`-- CreateIndex creates the "%s" index on the "%s" table.
CREATE UNIQUE INDEX "%s" ON "%s"("%s");

`, indexName, tableName, indexName, tableName, columnName))
		} else if strings.Contains(tag, "index") {
			indexName := fmt.Sprintf("idx_%s_%s", tableName, columnName)
			file.WriteString(fmt.Sprintf(`-- CreateIndex creates the "%s" index on the "%s" table.
CREATE INDEX "%s" ON "%s"("%s");

`, indexName, tableName, indexName, tableName, columnName))
		}
	}
}

func extractTagValue(tag, prefix string) string {
	start := strings.Index(tag, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.Index(tag[start:], ";")
	if end == -1 {
		return tag[start:]
	}
	return tag[start : start+end]
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
