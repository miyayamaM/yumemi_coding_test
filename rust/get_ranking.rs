use entites::player::Player;
use entites::player_list::PlayerList;
use entites::same_score_group::SameScoreGroup;
use entites::same_score_group_list::SameScoreGroupList;
use std::env;
use std::fs::File;
use std::io::{BufRead, BufReader, BufWriter, Write};

mod entites;

const OUTPUT_FILE_NAME: &str = "output.csv";
const OUTPUT_FILE_HEADER: [&str; 3] = ["rank", "player_id", "mean_score"];
const OUTPUT_LINES_LIMIT: usize = 10;

fn main() {
    //csvの読み込み
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        panic!("csvファイルを指定してください。 USAGE: $./main <example.csv>")
    };
    let file = File::open(&args[1]).expect("ファイルを開けませんでした");
    let mut reader = BufReader::new(file);

    //スコアをユーザーごとに集計
    let player_list: PlayerList = aggregate_score(&mut reader);

    //ユーザーをid順にソート
    let player_list_sorted_by_id = player_list.sort_by_player_id();

    //平均スコアでユーザーをグループ分け
    let group_by_mean_score: SameScoreGroupList = player_list_sorted_by_id.group_by_mean_score();

    //平均スコアでsort
    let sorted_groups = group_by_mean_score.sort_by_score();

    //fileの書き込み
    let file = File::create(OUTPUT_FILE_NAME).expect("ファイルの生成に失敗しました");
    let mut writer = BufWriter::new(file);
    output_ranking_as_csv(
        &mut writer,
        OUTPUT_FILE_HEADER,
        sorted_groups,
        OUTPUT_LINES_LIMIT,
    );
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

fn output_ranking_as_csv<W: Write>(
    writer: &mut BufWriter<W>,
    column_names: [&str; 3],
    score_groups: SameScoreGroupList,
    limit: usize,
) {
    let header = column_names.join(",") + "\n";
    writer
        .write(header.as_bytes())
        .expect("ヘッダーの書き込みに失敗しました");

    let mut rank = 1;
    for group in score_groups.groups.iter() {
        if rank > limit {
            break;
        }
        for player_id in group.player_ids.iter() {
            let line = format!("{},{},{}\n", rank, player_id, group.score);
            writer
                .write(line.as_bytes())
                .expect("ファイルへの書き込みに失敗しました");
        }
        rank += group.players_counts();
    }
}
