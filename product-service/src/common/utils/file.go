package utils

import (
	"fmt"
	"os"
	"product-service/src/common/log"
	"regexp"
	"time"
)

func CreateUniqueFileName(fileName string) string {
	re := regexp.MustCompile("[%?#&=]")
	encodedName := re.ReplaceAllString(fileName, "-")
	epochMillis := time.Now().UnixMilli()

	fileName = fmt.Sprintf("%d-%s", epochMillis, encodedName)
	return fileName
}

func CheckExisDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
func DeleteFile(path string) {
	if _, err := os.Stat(path); err == nil {
		err := os.Remove(path)
		if err != nil {
			log.Logger.Error(err.Error())
		}
	}
}
