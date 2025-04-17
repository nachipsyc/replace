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

// グローバル変数の宣言
var (
	targetDir string
	hookWord  string
	altWord   string
)

func main() {
	// コマンドの引数が正しく入力されているかの検証を先に行う
	if !isValidCommandArgs() {
		return // エラーがあれば処理を中断
	}

	// 直接変数に値を設定
	flag.StringVar(&targetDir, "targetDir", "", "対象ディレクトリのパス")
	flag.StringVar(&hookWord, "hookWord", "", "置換対象のフックワード")
	flag.StringVar(&altWord, "altWord", "", "新しいキーワード")

	// flag.Parseの前にフラグの区切り文字を空文字に設定
	flag.CommandLine.SetOutput(os.Stdout)
	oldUsage := flag.CommandLine.Usage
	flag.CommandLine.Usage = func() {
		oldUsage()
		os.Exit(1)
	}

	// 入力をパース
	flag.Parse()

	files, err := os.ReadDir(targetDir)
	if err != nil {
		fmt.Printf("Read Error: 入力ディレクトリを読み込めません: %v\n", err)
		printUsage()
		return
	}

	// ディレクトリが空の場合のチェック
	if len(files) == 0 {
		fmt.Println("Read Error: 入力ディレクトリが空です。")
		printUsage()
		return
	}

	// ファイルをリネーム
	renameFiles(files, hookWord, altWord)
}

// used by main()
func renameFiles(jpegFiles []fs.DirEntry, hookWord string, altWord string) {
	for _, jpegFile := range jpegFiles {
		originalFileName := jpegFile.Name()

		// フックワードを含む場合のみ処理
		if strings.Contains(originalFileName, hookWord) {
			newFileName := strings.Replace(originalFileName, hookWord, altWord, 1)

			// ファイル名のみ変更 (画像データは変更しない)
			oldPath := filepath.Join(targetDir, originalFileName)
			newPath := filepath.Join(targetDir, newFileName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("failed: %s -> %s, エラー: %v\n", oldPath, newPath, err)
			} else {
				fmt.Printf("success: %s -> %s\n", originalFileName, newFileName)
			}
		}
	}

	fmt.Println("replace completed!!")
}

func isValidCommandArgs() bool {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		// 引数がフラグ形式（-で始まる）かチェック
		if strings.HasPrefix(arg, "-") {
			// 引数名が有効かチェック
			validFlags := map[string]bool{
				"-targetDir": true,
				"-hookWord":  true,
				"-altWord":   true,
			}

			// 引数と値が=で結合されているかチェック
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				fmt.Printf("Command Error: 引数 '%s' は'='で値と結合する必要があります。\n", arg)
				printUsage()
				return false
			}

			argName := parts[0]
			value := parts[1]

			// targetDirの場合はダブルクオーテーションのチェックをスキップ
			if argName == "-targetDir" {
				if !isValidPath(value) {
					fmt.Printf("Command Error: 無効なディレクトリパス '%s' が指定されました。\n", value)
					printUsage()
					return false
				}
				targetDir = value
				continue // 次の引数へ
			}

			// hookWordと altWordの場合はダブルクオーテーションのチェックを削除
			if argName == "-hookWord" || argName == "-altWord" {
				// 特殊文字のチェック
				if !isValidWord(value) {
					fmt.Printf("Command Error: 無効なワード '%s' が指定されました。\n", value)
					printUsage()
					return false
				}
			}

			if !validFlags[argName] {
				fmt.Printf("Command Error: 不正な引数名です: '%s'\n", argName)
				printUsage()
				return false
			}

			// 各引数の値の検証
			switch argName {
			case "-hookWord":
				if value == "" {
					fmt.Println("Command Error: フックワードが指定されていません。")
					printUsage()
					return false
				}
				hookWord = value

			case "-altWord":
				if value == "" {
					fmt.Println("Command Error: 代わりのワードが指定されていません。")
					printUsage()
					return false
				}
				altWord = value
			}
		} else {
			fmt.Printf("Command Error: 引数 '%s' は不正な形式です。引数名は'-'で始まる必要があります。\n", arg)
			printUsage()
			return false
		}
	}

	// inputDirまたはoutputDirが空の場合のチェック
	if targetDir == "" {
		fmt.Println("Command Error: 入力ディレクトリが指定されていません。")
		printUsage()
		return false
	}

	return true // エラーがなければtrueを返す
}

// 有効なパスかどうかをチェックする関数（既存の関数を拡張）
func isValidPath(path string) bool {
	// パスが空でないことを確認
	if path == "" {
		return false
	}

	// パスに使用できない文字が含まれていないか確認
	invalidChars := []rune{'<', '>', ':', '"', '|', '?', '*'}
	for _, char := range path {
		for _, invalid := range invalidChars {
			if char == invalid {
				return false
			}
		}
	}

	return true
}

// 有効なワードかどうかをチェックする関数
func isValidWord(word string) bool {
	// ワードが空でないことを確認
	if word == "" {
		return false
	}

	// ワードに使用できない文字が含まれていないか確認
	invalidChars := []rune{'<', '>', ':', '"', '|', '?', '*'}
	for _, char := range word {
		for _, invalid := range invalidChars {
			if char == invalid {
				return false
			}
		}
	}

	return true
}

// 使い方を表示する関数
func printUsage() {
	fmt.Println("使い方:")
	fmt.Println("  -targetDir=<path>  : 入力画像ディレクトリのパス")
	fmt.Println("  -hookWord=<word> : 置き換えるワード")
	fmt.Println("  -altWord=<word> : 代わりのワード")
}
