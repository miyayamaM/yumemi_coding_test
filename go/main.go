package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	// コマンドライン引数を受ける
	if len(os.Args) != 2 {
		fmt.Println("コマンドライン引数が不正です")
		fmt.Println("Usage: go run main.go <file_path>")
	}

	filepath := os.Args[1]

	// ファイルを開く
	file, err := os.Open(filepath)

	if err != nil {
		fmt.Println("ファイルを開けませんでした")
		return
	}
	defer file.Close()

	// CSVリーダーを作成
	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()

		fmt.Println(line)

		if err == io.EOF {
			fmt.Println("CSVの読み込み完了")
			return
		}

		if err != nil {
			fmt.Println("CSVの読み込みエラー:", err)
			return
		}
	}
}
