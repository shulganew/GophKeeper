CREATE USER keeper WITH ENCRYPTED PASSWORD '1';
CREATE DATABASE keeper;

GRANT ALL PRIVILEGES ON DATABASE keeper TO keeper;

-- need for migrations (issue https://github.com/golang-migrate/migrate/issues/826)
ALTER DATABASE keeper OWNER TO keeper;