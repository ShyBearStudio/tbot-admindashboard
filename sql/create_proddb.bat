call heroku pg:psql --app tbotadmin < create_tables.sql
call heroku pg:psql --app tbotadmin < insert_testdata.sql
