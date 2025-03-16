#[derive(Debug, Clone)]

pub struct Player {
    pub id: String,
    pub total_score: usize,
    pub play_counts: usize,
}

impl Player {
    pub fn get_mean_score(&self) -> usize {
        (self.total_score as f32 / self.play_counts as f32).round() as usize
    }

    pub fn add_game_score(&mut self, score: usize) {
        self.total_score += score;
        self.play_counts += 1;
    }

    pub fn get_id_number(&self) -> usize {
        self.id
            .replace("player", "")
            .parse::<usize>()
            .expect("idを数字に変換できません")
    }
}
