#!/bin/sh
cd /app
gcc -Wimplicit-int -o main main.c
./main > __output 2> __error