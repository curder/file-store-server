package meta

import (
    "github.com/curder/file-store-server/db"
)

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

// 新增或者更新文件原始信息到数据库
func UpdateFileMetaDB(fileMeta FileMeta) bool {
    return db.OnFileUploadFinished(fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize, fileMeta.Location)
}

// 通过Sha1获取文件原始信息
func GetFileMeta(fileSha1 string) FileMeta {
    return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
    file, err := db.GetFileMeta(fileSha1)
    if err != nil {
        return FileMeta{}, err
    }
    fileMeta := FileMeta{
        FileSha1: file.FileSha1,
        FileName: file.FileName.String,
        FileSize: file.FileSize.Int64,
        Location: file.Location.String,
    }
    return fileMeta, nil
}

// 删除文件
func RemoveFileMeta(fileSha1 string) {
    delete(fileMetas, fileSha1) // TODO 多线程的时候这样删除会不安全
}
