package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
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
	header, err := reader.Read()
	if err != nil {
		fmt.Println("CSVのヘッダー行読み込みエラー:", err)
		return
	}

	fmt.Println(header)

	players := make(map[string]Player)

	for {
		line, err := reader.Read()

		if err == io.EOF {
			fmt.Println("CSVの読み込み完了")
			break
		}

		// 空の行をスキップ
		if len(line) == 0 {
			continue
		}

		if err != nil {
			fmt.Println("CSVの読み込みエラー:", err)
			return
		}

		_, player_id, score_str := line[0], line[1], line[2]

		score, err := strconv.Atoi(score_str)
		if err != nil {
			fmt.Println("スコアの変換エラー:", err)
			return
		}

		if player, exists := players[player_id]; exists {
			// 既存のプレイヤーがいる場合、スコアを追加し、プレイ回数を増やす
			player.AddScore(score)
			player.IncrementPlayingCount()
			players[player_id] = player
		} else {
			// 新しいプレイヤーを追加
			players[player_id] = Player{
				PlayerId:     player_id,
				TotalScore:   score,
				PlayingCount: 1,
			}
		}
	}
	fmt.Println(players)
}

// プレイヤーのスコアを記録する構造体
type Player struct {
	PlayerId     string
	TotalScore   int
	PlayingCount int
}

// TotalScoreにスコアを追加するメソッド
func (p *Player) AddScore(score int) {
	p.TotalScore += score
}

// PlayingCountを1増やすメソッド
func (p *Player) IncrementPlayingCount() {
	p.PlayingCount++
}

// 平均スコアの取得
func (p Player) AvarageScore() int {
	average := float64(p.TotalScore) / float64(p.PlayingCount)
	return int(math.Round(average))
}
