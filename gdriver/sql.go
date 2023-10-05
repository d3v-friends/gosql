package gdriver

// mysql5
const (
	mysql5CreateTable       = "CREATE TABLE %s.%s (id char(36) primary key)"
	mysql5CreateSchema      = "CREATE SCHEMA '%s' CHARACTER SET 'utf8mb4'"
	mysql5CountTable        = "SELECT COUNT(NAME) AS cnt FROM information_schema.INNODB_SYS_TABLES WHERE NAME = '%s'"
	mysql5CountSchema       = "SELECT COUNT(*) AS cnt FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = '%s'"
	mysql5ReadColumnInfo    = "SELECT COLUMN_NAME, IS_NULLABLE, DATA_TYPE, COLUMN_KEY, COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_NAME = '%s' AND TABLE_SCHEMA = '%s'"
	mysql5CreateColumn      = "ALTER TABLE %s ADD %s %s"
	mysql5AlterColumn       = "ALTER TABLE %s MODIFY %s %s"
	mysql5ReadIndexes       = "SHOW INDEX FROM %s.%s WHERE KEY_NAME = '%s'"
	mysql5CreateIndex       = "CREATE INDEX %s ON %s.%s (%s)"
	mysql5CreateUniqueIndex = "CREATE UNIQUE INDEX %s ON %s.%s (%s)"
	mysql5DropIndex         = "ALTER TABLE %s.%s DROP INDEX %s"
)
