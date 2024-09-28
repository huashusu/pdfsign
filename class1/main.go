package main

func main() {
	// 签名，输入的是源文件，输出是将源文件签名后保存的签名PDF文件
	Sign("./class1/input.pdf", "./class1/output.pdf")
	// 验证，验证输出的签名文件
	Valid("./class1/output.pdf")
}
