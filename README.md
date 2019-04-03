# lottery_client_server
## 编译优化
~~~bash
go build -ldflags '-w -s'
~~~
-w 禁止生成debug信息(使用后，无法使用 gdb 进行调试）

-s 禁用符号表

这可以缩小可执行文件的体积
