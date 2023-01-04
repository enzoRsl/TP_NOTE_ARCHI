@echo off

if not exist "%CD%\build" mkdir "-p" "%CD%\build"
go "build" "-o" "./build/" "./..."