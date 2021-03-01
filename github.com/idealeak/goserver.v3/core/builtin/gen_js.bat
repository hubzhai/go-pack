@echo off

cd protocol

for %%i in (*.proto) do (
	protoc3 --js_out=import_style=commonjs,binary:. %%i
)

cd ..


if not exist protocol_js mkdir protocol_js
move /y protocol\*.js protocol_js\


ping -n 3 127.0.0.1>nul
