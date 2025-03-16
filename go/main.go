package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"regexp"
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

	players := make(map[PlayerId]Player)

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

		_, player_id_str, score_str := line[0], line[1], line[2]

		player_id, err := NewPlayerId(player_id_str)
		if err != nil {
			fmt.Println("不正なplayer_id:", err)
			return
		}

		score, err := strconv.Atoi(score_str)
		if err != nil {
			fmt.Println("スコアの変換エラー:", err)
			return
		}

		if player, exists := players[player_id]; !exists {
			// 新しいプレイヤーを追加
			players[player_id] = Player{
				PlayerId:     player_id,
				TotalScore:   score,
				PlayingCount: 1,
			}
		} else {
			// 既存のプレイヤーがいる場合、スコアを追加し、プレイ回数を増やす
			player.AddScore(score)
			player.IncrementPlayingCount()
			players[player_id] = player
		}
	}

	// プレイヤーデータを平均スコアごとにグルーピング
	players_grouped_by_avg_score := groupPlayersByAverageScore(players)

	// 書き込み
	if err := writeCSV("output.csv", players_grouped_by_avg_score); err != nil {
		fmt.Println("CSVファイルの書き込みエラー:", err)
	}
}

// プレイヤーデータを平均スコアごとにグルーピングする関数
// 返り値は平均スコアを key, （同じ平均スコアの）Playerの配列 value にもつ map
func groupPlayersByAverageScore(players map[PlayerId]Player) map[int][]Player {
	players_grouped_by_avg_score := make(map[int][]Player)

	for _, player := range players {
		avgScore := player.AvarageScore()
		players_grouped_by_avg_score[avgScore] = append(players_grouped_by_avg_score[avgScore], player)
	}

	return players_grouped_by_avg_score
}

// CSVファイルに書き込む関数
func writeCSV(filename string, players_grouped_by_avg_score map[int][]Player) error {
	outputFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("CSVファイルの作成エラー: %w", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// ヘッダーを書き込む
	header := []string{"rank", "player_id", "mean_score"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("CSVヘッダーの書き込みエラー: %w", err)
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
				string(player.PlayerId),
				strconv.Itoa(player.AvarageScore()),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("CSVレコードの書き込みエラー: %w", err)
			}
		}
		rank += len(player_group)
		current_player += len(player_group)
		if current_player > max_player {
			break
		}
	}

	return nil
}

// プレイヤーのスコアを記録する構造体
type Player struct {
	PlayerId     PlayerId
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

type PlayerId string

func NewPlayerId(v string) (PlayerId, error) {
	pattern := `player\d{4}$`

	// 正規表現をコンパイル
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return "", err
	}

	// マッチングを行う
	matches := re.FindString(v)

	if matches == "" {
		fmt.Printf("不正なplayer_id: %s\n", v)
		return "", err
	} else {
		return PlayerId(matches), nil
	}
}
