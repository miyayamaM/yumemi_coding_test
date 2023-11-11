use Player;

#[derive(Debug, Clone)]
pub struct PlayerList {
    pub players: Vec<Player>,
}

impl PlayerList {
    fn contains_id(&self, player_id: &str) -> bool {
        self.players
            .iter()
            .find(|player| player.id == player_id.to_string())
            .is_some()
    }

    pub fn initialize_player(&mut self, player_id: &str) {
        if self.contains_id(&player_id) {
            return;
        }

        self.players.push(Player {
            id: player_id.to_string(),
            total_score: 0,
            play_counts: 0,
        })
    }

    pub fn add_player_score(&mut self, player_id: &str, score: usize) {
        if let Some(player) = self
            .players
            .iter_mut()
            .find(|player| player.id == player_id.to_string())
        {
            player.add_game_score(score);
        }
    }

    pub fn sort_by_player_id(&self) -> PlayerList {
        let mut player_list = self.clone();
        player_list
            .players
            .sort_by(|player1, player2| player2.get_id_number().cmp(&player1.get_id_number()));
        return player_list;
    }
}
