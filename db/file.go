package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"database/sql"
)

// 文件上传完成，保存mysql
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`" +
	", `file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statemen, err:"+ err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}


//从 mysql 获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statemen, err:"+ err.Error())
		return nil, err
	}
	defer stmt.Close()

	var fFile = TableFile{}
	err = stmt.QueryRow(filehash).Scan(&fFile.FileHash, &fFile.FileName, &fFile.FileSize, &fFile.FileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &fFile, nil
}