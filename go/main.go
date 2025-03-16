package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"slices"
	"sort"
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
	_, err = reader.Read()
	if err != nil {
		fmt.Println("CSVのヘッダー行読み込みエラー:", err)
		return
	}

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

	// 書き込み
	outputFile, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("CSVファイルの作成エラー:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// ヘッダーを書き込む
	header := []string{"rank", "player_id", "mean_score"}
	if err := writer.Write(header); err != nil {
		fmt.Println("CSVヘッダーの書き込みエラー:", err)
		return
	}

	// プレイヤーデータを書き込む
	players_grouped_by_avg_score := make(map[int][]Player)

	for _, player := range players {
		if player_group, exists := players_grouped_by_avg_score[player.AvarageScore()]; exists {
			// 同じ平均スコアのプレイヤーは同一グループにまとめる
			players_grouped_by_avg_score[player.AvarageScore()] = append(player_group, player)
		} else {
			// 新しいプレイヤーを追加
			players_grouped_by_avg_score[player.AvarageScore()] = append(players_grouped_by_avg_score[player.AvarageScore()], player)
		}
	}
	// 平均スコア順にソート
	players_sorted_keys := slices.Collect(maps.Keys(players_grouped_by_avg_score))
	sort.Sort(sort.Reverse(sort.IntSlice(players_sorted_keys)))

	// 書き込み
	rank := 1
	max_player := 10
	current_player := 0
	for _, avg_score := range players_sorted_keys {
		player_group := players_grouped_by_avg_score[avg_score]
		for _, player := range player_group {
			record := []string{
				strconv.Itoa(rank),
				player.PlayerId,
				strconv.Itoa(player.AvarageScore()),
			}
			if err := writer.Write(record); err != nil {
				fmt.Println("CSVレコードの書き込みエラー:", err)
				return
			}
		}
		rank += len(player_group)
		current_player += len(player_group)
		if current_player > max_player {
			break
		}
	}
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
