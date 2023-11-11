use entites::player::Player;
use entites::player_list::PlayerList;
use std::collections::HashMap;
use std::env;
use std::fs::File;
use std::io::{BufRead, BufReader, BufWriter, Write};

mod entites;

fn main() {
    //csvの読み込み
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        panic!("csvファイルを指定してください。 USAGE: $./main <example.csv>")
    };
    let file = File::open(&args[1]).expect("ファイルを開けませんでした");
    let mut reader = BufReader::new(file);

    //スコアをユーザーごとに集計
    let player_list = aggregate_score(&mut reader);

    //ユーザーをid順にソート
    let player_list_sorted_by_id = player_list.sort_by_player_id();
    let players = player_list_sorted_by_id.players;

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

fn aggregate_score(file: &mut dyn BufRead) -> PlayerList {
    let mut player_list = PlayerList {
        players: Vec::new(),
    };
    for line in file.lines().skip(1) {
        let line = line.expect("ファイルの読み取りに失敗しました");

        let score: Vec<&str> = line.split(",").collect();
        let player_id = score[1];
        let game_score: usize = score[2].parse().expect("数字でないスコアがあります");

        player_list.initialize_player(player_id);
        player_list.add_player_score(player_id, game_score)
    }
    player_list
}

fn group_by_mean_score(players: Vec<Player>) -> HashMap<usize, Vec<String>> {
    let mut mean_scores: HashMap<usize, Vec<String>> = HashMap::new();

    for player in players {
        match mean_scores.get_mut(&player.get_mean_score()) {
            None => {
                mean_scores.insert(player.get_mean_score(), vec![player.id]);
            }
            Some(same_score_players) => {
                same_score_players.push(player.id);
            }
        };
    }
    mean_scores
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
