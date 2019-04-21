# lottery_client_server
## 编译信息
* `golang`版本：`1.12`
* 依赖的外部库：
    * `github.com/gorilla/mux`
    * `github.com/gorilla/websocket`
    * `github.com/jinzhu/gorm`
    * `github.com/jinzhu/gorm/dialects/sqlite`
    * `github.com/mattn/go-sqlite3`
* 设置信息储存在 `config/configures.go`中

## 编译优化
~~~bash
go build -ldflags '-w -s'
~~~
-w 禁止生成debug信息(使用后，无法使用 gdb 进行调试）

-s 禁用符号表

这可以缩小可执行文件的体积
