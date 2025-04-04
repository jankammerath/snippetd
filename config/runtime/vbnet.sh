#!/bin/sh
export DOTNET_CLI_TELEMETRY_OPTOUT=1
export DOTNET_NOLOGO=1
export DOTNET_CONSOLE_UNBUFFERED=1
cd /app
dotnet new console -o tempProject -lang vb > /dev/null
cp Program.vb tempProject/
cd tempProject
dotnet run