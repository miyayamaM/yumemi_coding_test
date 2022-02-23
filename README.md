## ファイルの説明
- main.rs

  処理のプログラム本体

- game_score_log.csv

  生データのcsvファイル。これを読み込んで処理するプログラムを書く。

- expected_output.csv

  出力の正解サンプル

## 使い方

```
$rustc main.rs && ./main game_score_log.csv
```

## 出力

```
rank,player_id,mean_score
1,player0001,10000
1,player0002,10000
3,player0003,9000
4,player0004,7000
5,player0005,1000
6,player0006,999
7,player0007,998
8,player0008,997
9,player0009,990
9,player0010,990
9,player0011,990
9,player0012,990
```
