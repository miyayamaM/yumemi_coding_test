use std::collections::HashSet;
use SameScoreGroup;

#[derive(Debug, Clone)]
pub struct SameScoreGroupList {
    pub groups: Vec<SameScoreGroup>,
}

impl SameScoreGroupList {
    pub fn contains_score_group(&self, score: usize) -> bool {
        self.groups
            .iter()
            .find(|group| group.score == score)
            .is_some()
    }
    pub fn initialize_group(&mut self, score: usize) {
        if self.contains_score_group(score) {
            return;
        }
        self.groups.push(SameScoreGroup {
            score: score,
            player_ids: HashSet::new(),
        });
    }

    pub fn add_to_group(&mut self, score: usize, player_id: String) {
        if let Some(group) = self.groups.iter_mut().find(|group| group.score == score) {
            group.add_player(player_id);
        }
    }

    pub fn sort_by_score(&self) -> SameScoreGroupList {
        let mut score_groups = self.groups.clone();

        score_groups.sort_by(|group1, group2| group2.score.cmp(&group1.score));
        SameScoreGroupList {
            groups: score_groups,
        }
    }
}
