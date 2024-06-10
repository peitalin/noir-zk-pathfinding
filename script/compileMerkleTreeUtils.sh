#!/bin/bash

cd ./frontend && pnpm install && npx tsx ./compileHasher.ts
cd ..