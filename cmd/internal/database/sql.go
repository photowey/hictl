package database

const (
	DriverMysql = "mysql"
	DsnTemplate = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

const (
	TableInfoSql = `
SELECT 
TABLE_NAME AS tableName, 
TABLE_COMMENT AS tableComment 
FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = ? `
	ColumnInfoSql = `
SELECT
TABLE_NAME AS tableName,
COLUMN_NAME AS columnName,
COLUMN_COMMENT AS columnComment,
IS_NULLABLE AS notNull,
DATA_TYPE AS dataType,
CHARACTER_MAXIMUM_LENGTH AS dataLength,
COLUMN_KEY AS primaryKey,
NUMERIC_PRECISION AS maxBit,
NUMERIC_SCALE AS minBit
FROM
information_schema.COLUMNS
WHERE
TABLE_SCHEMA = ?
AND TABLE_NAME = ?
ORDER BY TABLE_NAME,
ORDINAL_POSITION`

	IndexInfoSql = `SELECT 
TABLE_NAME AS tableName, 
INDEX_NAME AS indexName, 
GROUP_CONCAT(COLUMN_NAME) AS indexColumn
FROM
information_schema.STATISTICS
WHERE
TABLE_SCHEMA = ?
AND TABLE_NAME = ?
GROUP BY TABLE_NAME, INDEX_NAME`
)
