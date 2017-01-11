call heroku pg:psql --app testchitchat < create_tables.sql
call heroku pg:psql --app testchitchat < insert_testdata.sql
