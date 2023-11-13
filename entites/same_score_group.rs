use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct SameScoreGroups {
    pub groups: HashMap<usize, Vec<String>>,
}

impl SameScoreGroups {
    pub fn initialize_group(&mut self, score: usize) {
        if self.groups.contains_key(&score) {
            return;
        }
        self.groups.insert(score, Vec::new());
    }

    pub fn add_to_group(&mut self, score: usize, player_id: String) {
        if let Some(group) = self.groups.get_mut(&score) {
            group.push(player_id);
        }
    }

    pub fn sort_by_score(&self) -> SameScoreGroups {
        let score_groups = self.groups.clone();

        // 比較するためにHashMapからタプルにする
        let mut score_group_tupples: Vec<(usize, Vec<String>)> = score_groups.into_iter().collect();
        score_group_tupples.sort_by(|x, y| y.0.cmp(&x.0));
        let groups_hashmap: HashMap<usize, Vec<String>> = score_group_tupples.into_iter().collect();
        SameScoreGroups {
            groups: groups_hashmap,
        }
    }
}
