package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dir string
var hookWord string

func main() {
	// コマンドライン引数を明示
	dirPtr := flag.String("dir", "", "対象ディレクトリのパス")
	hookWordPtr := flag.String("word", "", "削除するフックワード")

	// 入力をパース
	flag.Parse()

	// 引数が正しくない場合は実行方法を明示
	if *dirPtr == "" || *hookWordPtr == "" {
		log.Fatal("使用方法: go run scrape.go -dir=<ディレクトリ> -word=<フックワード>")
	}

	dir = *dirPtr
	hookWord = *hookWordPtr

	// 対象ディレクトリの中の全ファイルを取得([]fs.DirEntry)
	files, err := getFiles(dir)

	if err != nil {
		log.Fatal(err)
	}

	// 取得したファイルからJPEGファイルのみを抽出([]fs.DirEntry)
	jpegFiles := extractJpegFiles(files)

	// ファイルをリネーム
	renameFiles(jpegFiles, hookWord)
}

// used by main()
func getFiles(sourceDir string) ([]fs.DirEntry, error) {
	// 対象ディレクトリの中のファイル全てを取得、格納
	files, err := os.ReadDir(sourceDir)

	// エラーがあれば"err"を返す
	if err != nil {
		return nil, err
	}

	return files, nil
}

// used by main()
func extractJpegFiles(files []fs.DirEntry) []fs.DirEntry {
	var jpegImages []fs.DirEntry
	for _, file := range files {
		switch filepath.Ext(file.Name()) {
		case ".jpeg", ".jpg", ".JPG":
			jpegImages = append(jpegImages, file)
		default:
		}
	}

	return jpegImages
}

// used by main()
func renameFiles(jpegFiles []fs.DirEntry, hookWord string) {
	for _, jpegFile := range jpegFiles {
		originalFileName := jpegFile.Name()

		// フックワードを含む場合のみ処理
		if strings.Contains(originalFileName, hookWord) {
			newFileName := strings.Replace(originalFileName, hookWord, "", 1)

			// ファイル名のみ変更 (画像データは変更しない)
			oldPath := filepath.Join(dir, originalFileName)
			newPath := filepath.Join(dir, newFileName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("failed: %s -> %s, エラー: %v\n", oldPath, newPath, err)
			} else {
				fmt.Printf("success: %s -> %s\n", originalFileName, newFileName)
			}
		}
	}

	fmt.Println("scrape completed!!")
}
