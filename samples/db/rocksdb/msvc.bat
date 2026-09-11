::@echo off
setlocal enabledelayedexpansion

set "VSWHERE_PATH=%ProgramFiles(x86)%\Microsoft Visual Studio\Installer\vswhere.exe"

for /f "usebackq tokens=*" %%i in (`"!VSWHERE_PATH!" -latest -products * -requires Microsoft.VisualStudio.Component.VC.Tools.x86.x64 -property installationPath 2^>nul`) do (
    set "VS_DIR=%%i"
)
set "VCVARS64_PATH=!VS_DIR!\VC\Auxiliary\Build\vcvars64.bat"
call "!VCVARS64_PATH!"
echo "CGO_LDFLAGS=%~2"
echo "CGO_CFLAGS=%~1"
where.exe cl.exe
where.exe link.exe
set "CGO_LDFLAGS=%~2"
set "CGO_CFLAGS=%~1"
call go build -o %~3 ./...
endlocal
