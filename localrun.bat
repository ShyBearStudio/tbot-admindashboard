%ECHO OFF
ECHO Running T-Bot Admin Dashboard locally ...

ECHO Building GO application ... 
go build && (
	ECHO === [done]
) || (
	ECHO === [failed]
	GOTO :FAILURE 
)

ECHO Run the application ...
START tbot-admindashboard.exe -config configs/localconfig.json && (
	ECHO === [done]
) || (
	ECHO === [failed]
	GOTO :FAILURE 
)

ECHO Open up chrome ...
START chrome "http://localhost:8080/" && (
	ECHO === [done]
) || (
	ECHO === [failed]
	GOTO :FAILURE 
)

PAUSE
ECHO Killing server process ...
taskkill /FI "IMAGENAME eq tbot-admin*" /f && (
	ECHO === [done]
) || (
	ECHO === [failed]
	GOTO :FAILURE 
)

ECHO === SUCCESS ===
EXIT /b


:FAILURE
ECHO === FAILURE ===




