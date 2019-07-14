package handler

import (
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
	"os"
	"filestore-server/meta"
	"time"
	"filestore-server/util"
	"encoding/json"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html 页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			return 
		}

		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data, err: %s", err.Error())
			return 
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			// #Todo: company
			// Location: "D:/goproject/src/filestore-server/tmp/" + head.Filename,
			// #Todo: private
			Location: "G:/goproject/src/filestore-server/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2016-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("failed to create file, err: %s", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file, err: %s", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Printf(fileMeta.FileSha1)
		// meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetaDB(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 文件上传成功提示
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// 获取文件元信息
func GetFileMateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["filehash"][0]
	// fMeta := meta.GetFileMeta(fileHash)
	fMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsSha1 := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(fsSha1)

	f, err := os.Open(fMeta.Location)
	if(err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if(err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/octect-stream")
	w.Header().Set("Content-disposition", "attachment;filename=\""+ fMeta.FileName +"\"")
	w.Write(data)
}

// 文件名元信息更改
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fsSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

    if(opType != "0") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if(r.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fMeta := meta.GetFileMeta(fsSha1)
	fMeta.FileName = newFileName
	// meta.UpdateFileMeta(fMeta)
	_ = meta.UpdateFileMetaDB(fMeta)

	data, err := json.Marshal(fMeta)
	if(err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fsSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fsSha1)
	w.WriteHeader(http.StatusOK)
}