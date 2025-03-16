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
	// コマンドライン引数を取得
	filepath, err := getFilePathFromArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

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

	// CSVファイルを読み込み、playerIdをkey、Playerをvalueに持つmapにまとめる
	players, err := readPlayersFromCSV(reader)
	if err != nil {
		fmt.Println("CSVの読み込みエラー:", err)
		return
	}

	// プレイヤーデータを平均スコアごとにグルーピング
	playersGroupedByAvgScore := groupPlayersByAverageScore(players)

	// CSVファイルに書き込み
	if err := writeCSV("output.csv", playersGroupedByAvgScore); err != nil {
		fmt.Println("CSVファイルの書き込みエラー:", err)
	}
}

// コマンドライン引数からファイルパスを取得する関数
func getFilePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("コマンドライン引数が不正です\nUsage: go run main.go <file_path>")
	}
	return os.Args[1], nil
}

// CSVからプレイヤーデータを読み込む関数
func readPlayersFromCSV(reader *csv.Reader) ([]Player, error) {
	// 同じPlayerIdのログは同じPlayerにまとめるので、
	// 検索性から、playerIdをキーとしたmapに格納する
	players := make(map[PlayerId]Player)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			fmt.Println("CSVの読み込み完了")
			break
		}
		if err != nil {
			return nil, fmt.Errorf("CSVの読み込みエラー: %w", err)
		}

		// 空の行をスキップ
		if len(line) == 0 {
			continue
		}

		_, playerIDStr, scoreStr := line[0], line[1], line[2]

		playerID, err := NewPlayerId(playerIDStr)
		if err != nil {
			return nil, fmt.Errorf("不正なplayer_id: %w", err)
		}

		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			return nil, fmt.Errorf("スコアの変換エラー: %w", err)
		}

		if player, exists := players[playerID]; !exists {
			// 新しいプレイヤーを追加
			players[playerID] = Player{
				PlayerId:     playerID,
				TotalScore:   score,
				PlayingCount: 1,
			}
		} else {
			// 既存のプレイヤーがいる場合、スコアを追加し、プレイ回数を増やす
			player.AddScore(score)
			player.IncrementPlayingCount()
			players[playerID] = player
		}
	}

	return slices.Collect(maps.Values(players)), nil
}

// プレイヤーデータを平均スコアごとにグルーピングする関数
// 返り値は平均スコアを key, （同じ平均スコアの）Playerの配列 value にもつ map
func groupPlayersByAverageScore(players []Player) map[int][]Player {
	playersGroupedByAvgScore := make(map[int][]Player)

	for _, player := range players {
		avgScore := player.AvarageScore()
		playersGroupedByAvgScore[avgScore] = append(playersGroupedByAvgScore[avgScore], player)
	}

	return playersGroupedByAvgScore
}

// CSVファイルに書き込む関数
func writeCSV(filename string, playersGroupedByAvgScore map[int][]Player) error {
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
	playersSortedKeys := slices.Collect(maps.Keys(playersGroupedByAvgScore))
	sort.Sort(sort.Reverse(sort.IntSlice(playersSortedKeys)))

	// 書き込み
	rank := 1
	maxPlayer := 10
	currentPlayer := 0
	for _, avgScore := range playersSortedKeys {
		playerGroup := playersGroupedByAvgScore[avgScore]
		for _, player := range playerGroup {
			record := []string{
				strconv.Itoa(rank),
				string(player.PlayerId),
				strconv.Itoa(player.AvarageScore()),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("CSVレコードの書き込みエラー: %w", err)
			}
		}
		rank += len(playerGroup)
		currentPlayer += len(playerGroup)
		if currentPlayer > maxPlayer {
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
