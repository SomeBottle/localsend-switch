package utils

// 执行 proto 文件编译

// go generate 中命令的执行目录是该模块所在的目录，因此还要向上一层 ../protos/ 查找 proto 文件
//go:generate mkdir -p ../generated
//go:generate bash -c "protoc --proto_path=../protos/ --go_out=../generated/ $(find ../protos/ -name '*.proto')"
