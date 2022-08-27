use std::collections::HashMap;
use std::env;
use std::fs::File;
use std::io::{BufRead, BufReader, BufWriter, Write};

fn main() {
    //csvの読み込み
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        panic!("csvファイルを指定してください。 USAGE: $./main <example.csv>")
    };
    let file_name = File::open(&args[1]).expect("ファイルを開けませんでした");
    let mut file = BufReader::new(file_name);

    //スコアをユーザーごとに集計
    let mut players: Vec<Player> = aggregate_score(&mut file);

    //ユーザーをid順にソート
    sort_players(&mut players);

    //平均スコアでユーザーをグループ分け
    let mean_scores: HashMap<usize, Vec<String>> = group_by_mean_score(players);

    //平均スコアでsort
    let mut sorted_mean_scores: Vec<(usize, Vec<String>)> = mean_scores.into_iter().collect();
    sorted_mean_scores.sort_by(|x, y| y.0.cmp(&x.0));

    //fileの書き込み
    let output_file_name = "output.csv";
    let file = File::create(output_file_name).expect("ファイルの生成に失敗しました");
    let mut writer = BufWriter::new(file);

    let column_names = vec!["rank", "player_id", "mean_score"];
    let limit = 10;
    output_ranking_as_csv(&mut writer, column_names, sorted_mean_scores, limit);
}

struct Player {
    id: String,
    total_score: usize,
    play_counts: usize,
}

fn aggregate_score(file: &mut dyn BufRead) -> Vec<Player> {
    let lines = file.lines();
    let mut players: Vec<Player> = Vec::new();
    for line in lines.skip(1) {
        let line = line.expect("ファイルの読み取りに失敗しました");

        let score: Vec<&str> = line.split(",").collect();
        let player_id: String = score[1].to_string();
        let game_score: usize = score[2].parse().expect("数字でないスコアがあります");

        match players.iter_mut().find(|player| player.id == player_id) {
            None => {
                let new_player = Player {
                    id: player_id,
                    total_score: game_score,
                    play_counts: 1,
                };
                players.push(new_player);
            }
            Some(player) => {
                player.total_score += game_score;
                player.play_counts += 1;
            }
        };
    }
    return players;
}

fn group_by_mean_score(players: Vec<Player>) -> HashMap<usize, Vec<String>> {
    let mut mean_scores: HashMap<usize, Vec<String>> = HashMap::new();

    for player in players {
        let mean_score = player.total_score / player.play_counts;
        match mean_scores.get_mut(&mean_score) {
            None => {
                let new_mean_score_player = vec![player.id];
                mean_scores.insert(mean_score, new_mean_score_player);
            }
            Some(same_score_players) => {
                same_score_players.push(player.id);
            }
        };
    }
    mean_scores
}

fn sort_players(players: &mut Vec<Player>) {
    players.sort_by_key(|player| {
        player
            .id
            .replace("player", "")
            .parse::<usize>()
            .expect("idを数字に変換できません")
    });
}

fn output_ranking_as_csv<W: Write>(
    writer: &mut BufWriter<W>,
    column_names: Vec<&str>,
    scores: Vec<(usize, Vec<String>)>,
    limit: usize,
) {
    let header = column_names.join(",") + "\n";
    writer
        .write(header.as_bytes())
        .expect("ヘッダーの書き込みに失敗しました");

    let mut index = 0;
    let mut counts = 0;
    let mut rank = 1;
    while counts < limit {
        //(score, [player_id])のタプル
        let score_with_players = &scores[index];
        for player_name in score_with_players.1.iter() {
            let line = format!("{},{},{}\n", rank, player_name, score_with_players.0);
            writer
                .write(line.as_bytes())
                .expect("ファイルへの書き込みに失敗しました");
            counts += 1;
        }
        rank += score_with_players.1.len();
        index += 1;
    }
}
