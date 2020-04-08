package meta

// 文件的原始信息结构
type FileMeta struct {
    FileSha1  string
    FileName  string
    FileSize  int64
    Location  string
    UpdatedAt string
}

var fileMetas map[string]FileMeta

func init() {
    fileMetas = make(map[string]FileMeta) // 初始化map
}

// 新增或者更新文件原始信息
func UpdateFileMeta(fileMeta FileMeta) {
    fileMetas[fileMeta.FileSha1] = fileMeta
}

// 通过Sha1获取文件原始信息
func GetFileMeta(fileSha1 string) FileMeta {
    return fileMetas[fileSha1]
}

// 删除文件
func RemoveFileMeta(fileSha1 string) {
    delete(fileMetas, fileSha1) // TODO 多线程的时候这样删除会不安全
}