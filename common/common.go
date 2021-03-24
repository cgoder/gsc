package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

//JsonFormat Json outupt.
func JsonFormat(v interface{}) string {
	// if msg == "" {
	// 	return ""
	// }
	var out bytes.Buffer

	bs, _ := json.Marshal(v)
	json.Indent(&out, bs, "", "\t")

	return out.String()
}

func GetFileMd5(filePath string) string {
	pFile, err := os.Open(filePath)
	defer pFile.Close()
	if err != nil {
		log.Errorln("open file fail! ", filePath, err)
		return ""
	}

	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}
