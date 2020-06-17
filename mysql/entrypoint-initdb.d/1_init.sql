-- phpMyAdmin SQL Dump
-- version 5.0.2
-- https://www.phpmyadmin.net/
--
-- ホスト: mysql:3306
-- 生成日時: 2020 年 6 月 17 日 12:46
-- サーバのバージョン： 8.0.20
-- PHP のバージョン: 7.4.5

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- データベース: `procon31`
--

-- --------------------------------------------------------

--
-- テーブルの構造 `matches`
--

CREATE TABLE `matches` (
  `id` int NOT NULL,
  `start_at` int NOT NULL,
  `turn_ms` int NOT NULL,
  `interval_ms` int NOT NULL,
  `turn_num` int NOT NULL,
  `field` json NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- テーブルの構造 `match_teams`
--

CREATE TABLE `match_teams` (
  `match_id` int NOT NULL,
  `local_team_id` int NOT NULL,
  `global_team_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `update_actions` json NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- テーブルの構造 `teams`
--

CREATE TABLE `teams` (
  `global_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- ダンプしたテーブルのインデックス
--

--
-- テーブルのインデックス `matches`
--
ALTER TABLE `matches`
  ADD PRIMARY KEY (`id`);

--
-- テーブルのインデックス `match_teams`
--
ALTER TABLE `match_teams`
  ADD PRIMARY KEY (`local_team_id`);

--
-- テーブルのインデックス `teams`
--
ALTER TABLE `teams`
  ADD PRIMARY KEY (`global_id`);

--
-- ダンプしたテーブルのAUTO_INCREMENT
--

--
-- テーブルのAUTO_INCREMENT `matches`
--
ALTER TABLE `matches`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- テーブルのAUTO_INCREMENT `match_teams`
--
ALTER TABLE `match_teams`
  MODIFY `local_team_id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
