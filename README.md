# Go 文件上传管理

```
go run main.go
```

打开浏览器访问 `http://127.0.0.1:8888/files/uploads`，尝试上传一张图片，图片会被保存在 `/tmp/go-test-files/`下。


- Mac生成1GB大小文件

```
dd if=/dev/zero of=fileName bs=1048576 count=1000
```
> `of` 文件名称
> `bs` 生成文件的大小：1024 * 1024 = 1GB

- Mac查看文件sha1值

```
shasum main.go
```
- `main.go` 文件名