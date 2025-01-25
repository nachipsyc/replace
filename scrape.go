package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dir string = ""
var hook_word string = ""

func main() {
	// 入力をパース
	flag.Parse()

	// 対象ディレクトリのパスを引数から取得(string)
	dir = flag.Arg(0)

	// 削ぎ落としたいファイル名をセット
	hook_word = flag.Arg(1)

	// 対象ディレクトリの中の全ファイルを取得([]fs.DirEntry)
	files, err := getFiles(dir)

	if err != nil {
		log.Fatal(err)
		panic(err)
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
		// オリジナルのファイル名を保管
		original_file_name := jpeg_file.Name()
		// フックワードを含む場合のみ
		if strings.Contains(original_file_name, hook_word) {
			// ファイルを画像として読み込み
			decoded_image, _ := decodeImage(jpeg_file)

			if decoded_image != nil {
				// リサイズした画像を書き込み
				encodeImage(decoded_image, original_file_name, hook_word)
			}
		}

	}

	fmt.Println("scrape completed!!")
}

// used by renameFiles()
func decodeImage(jpeg_file fs.DirEntry) (image.Image, error) {
	io_file, err := conversionToReader(jpeg_file)
	if err != nil {
		return nil, err
	}

	decoded_image, _, err := image.Decode(io_file)
	if err != nil {
		return nil, err
	}
	return decoded_image, nil
}

// used by renameFiles()
func encodeImage(decoded_image image.Image, original_file_name string, hook_word string) error {
	if strings.Contains(original_file_name, hook_word) {

	}
	new_file_name := strings.Replace(original_file_name, hook_word, "", 1)
	output, err := os.Create(dir + "/" + new_file_name)
	if err != nil {
		return err
	}

	defer output.Close()

	opts := &jpeg.Options{Quality: 100}
	if err := jpeg.Encode(output, decoded_image, opts); err != nil {
		return err
	}

	return nil
}

// used by decodeImage()
func conversionToReader(jpeg_file fs.DirEntry) (io.Reader, error) {
	io_file, err := os.Open(dir + "/" + jpeg_file.Name())
	if err != nil {
		return nil, err
	}
	return io_file, nil
}
