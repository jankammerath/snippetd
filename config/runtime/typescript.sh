#!/bin/sh
cd /app
tsc index.ts
node index.js > __output 2> __error