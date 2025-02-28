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
var newWord string

func main() {
	// コマンドライン引数を明示
	dirPtr := flag.String("dir", "", "対象ディレクトリのパス")
	hookWordPtr := flag.String("word", "", "置換対象のフックワード")
	newWordPtr := flag.String("new", "", "新しいキーワード")

	// 入力をパース
	flag.Parse()

	// 引数が正しくない場合は実行方法を明示
	if *dirPtr == "" || *hookWordPtr == "" || *newWordPtr == "" {
		log.Fatal("使用方法: go run scrape.go -dir=<ディレクトリ> -word=<フックワード> -new=<新しいキーワード>")
	}

	dir = *dirPtr
	hookWord = *hookWordPtr
	newWord = *newWordPtr

	// 対象ディレクトリの中の全ファイルを取得([]fs.DirEntry)
	files, err := getFiles(dir)

	if err != nil {
		log.Fatal(err)
	}

	// ファイルをリネーム
	renameFiles(files, hookWord, newWord)
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
func renameFiles(jpegFiles []fs.DirEntry, hookWord string, newWord string) {
	for _, jpegFile := range jpegFiles {
		originalFileName := jpegFile.Name()

		// フックワードを含む場合のみ処理
		if strings.Contains(originalFileName, hookWord) {
			newFileName := strings.Replace(originalFileName, hookWord, newWord, 1)

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
