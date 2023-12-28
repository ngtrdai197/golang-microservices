
#!/bin/sh

#* variables
PROTO_PATH=./proto
PROTO_OUT=./api
IDL_PATH=./api
DOC_OUT=./docs

rm -rf ${PROTO_OUT}

#! create folders if not exists
mkdir -p ${DOC_OUT}/html
mkdir -p ${DOC_OUT}/markdown
mkdir -p ${DOC_OUT}/swagger
mkdir -p ${IDL_PATH}

#* gen normal proto
protoc \
    ${PROTO_PATH}/*/*.proto \
    -I=/usr/local/include \
    --proto_path=${PROTO_PATH} \
    --go_out=:${IDL_PATH} \
    --validate_out=lang=go:${IDL_PATH} \
    --go-grpc_out=:${IDL_PATH} \
    --grpc-gateway_out=:${IDL_PATH} \
    --openapiv2_out=:${DOC_OUT}/swagger \
    --doc_out=:${DOC_OUT}/html --doc_opt=html,index.html

#! remove permission folders
chmod -R 777 ${PROTO_OUT}
chmod -R 777 ${DOC_OUT}