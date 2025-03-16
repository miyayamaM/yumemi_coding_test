use std::collections::HashSet;

#[derive(Debug, Clone)]
pub struct SameScoreGroup {
    pub score: usize,
    pub player_ids: HashSet<String>,
}

impl SameScoreGroup {
    pub fn add_player(&mut self, player_id: String) {
        self.player_ids.insert(player_id);
    }

    pub fn players_counts(&self) -> usize {
        self.player_ids.len()
    }
}
