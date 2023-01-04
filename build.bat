@echo off

if not exist "%CD%\build" mkdir "%CD%\build"
go "build" "-o" "./build/" "./..."

echo Build folder created