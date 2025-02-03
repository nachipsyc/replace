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
var hook_word string

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
	hook_word = *hookWordPtr

	// 対象ディレクトリの中の全ファイルを取得([]fs.DirEntry)
	files, err := getFiles(dir)

	if err != nil {
		log.Fatal(err)
	}

	// 取得したファイルからJPEGファイルのみを抽出([]os.FileInfo)
	jpeg_files := extractJpegFiles(files)

	// ファイルをリネーム
	renameFiles(jpeg_files, hook_word)
}

// used by main()
func getFiles(source_dir string) ([]fs.DirEntry, error) {
	// 対象ディレクトリの中のファイル全てを取得、格納
	files, err := os.ReadDir(source_dir)

	// エラーがあれば"err"を返す
	if err != nil {
		return nil, err
	}

	return files, nil
}

// used by main()
func extractJpegFiles(files []fs.DirEntry) []fs.DirEntry {
	var jpeg_images []fs.DirEntry
	for _, file := range files {
		switch filepath.Ext(file.Name()) {
		case ".jpeg", ".jpg", ".JPG":
			jpeg_images = append(jpeg_images, file)
		default:
		}
	}

	return jpeg_images
}

// used by main()
func renameFiles(jpeg_files []fs.DirEntry, hook_word string) {
	for _, jpeg_file := range jpeg_files {
		original_file_name := jpeg_file.Name()

		// フックワードを含む場合のみ処理
		if strings.Contains(original_file_name, hook_word) {
			new_file_name := strings.Replace(original_file_name, hook_word, "", 1)

			// ファイル名のみ変更 (画像データは変更しない)
			oldPath := filepath.Join(dir, original_file_name)
			newPath := filepath.Join(dir, new_file_name)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("failed: %s -> %s, エラー: %v\n", oldPath, newPath, err)
			} else {
				fmt.Printf("scraping: %s -> %s\n", original_file_name, new_file_name)
			}
		}
	}

	fmt.Println("scrape completed!!")
}
