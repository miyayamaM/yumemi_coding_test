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
    let file = File::open(&args[1]).expect("ファイルを開けませんでした");
    let mut reader = BufReader::new(file);

    //スコアをユーザーごとに集計
    let mut players = aggregate_score(&mut reader);

    //ユーザーをid順にソート
    sort_players(&mut players);

    //平均スコアでユーザーをグループ分け
    let mean_scores = group_by_mean_score(players);

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

impl Player {
    fn mean_score(&self) -> usize {
        (self.total_score as f32 / self.play_counts as f32).round() as usize
    }
}

fn aggregate_score(file: &mut dyn BufRead) -> Vec<Player> {
    let mut players: Vec<Player> = Vec::new();
    for line in file.lines().skip(1) {
        let line = line.expect("ファイルの読み取りに失敗しました");

        let score: Vec<&str> = line.split(",").collect();
        let player_id = score[1].to_string();
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
    players
}

fn group_by_mean_score(players: Vec<Player>) -> HashMap<usize, Vec<String>> {
    let mut mean_scores: HashMap<usize, Vec<String>> = HashMap::new();

    for player in players {
        match mean_scores.get_mut(&player.mean_score()) {
            None => {
                mean_scores.insert(player.mean_score(), vec![player.id]);
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

    let mut rank = 1;
    for (score, player_ids) in scores.iter() {
        if rank > limit {
            break;
        }
        for player_id in player_ids {
            let line = format!("{},{},{}\n", rank, player_id, score);
            writer
                .write(line.as_bytes())
                .expect("ファイルへの書き込みに失敗しました");
        }
        rank += player_ids.len();
    }
}
