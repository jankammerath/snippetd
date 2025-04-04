#!/bin/sh
export DOTNET_CLI_TELEMETRY_OPTOUT=1
export DOTNET_NOLOGO=1
export DOTNET_CONSOLE_UNBUFFERED=1
cd /app
dotnet new console -o tempProject > /dev/null
cp Program.cs tempProject/
cd tempProject
dotnet run > ../__output 2> ../__error