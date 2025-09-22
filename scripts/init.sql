-- データベース作成
CREATE DATABASE IF NOT EXISTS lovender;

USE lovender;

-- 1) users
CREATE TABLE users (
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name            VARCHAR(191)     NOT NULL,
  email           VARCHAR(255)     NOT NULL,
  password_hash   VARCHAR(255)     NOT NULL,
  created_at      DATETIME(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at      DATETIME(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 2) oshii
CREATE TABLE oshis (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id       BIGINT UNSIGNED NOT NULL,
  name          VARCHAR(191)     NOT NULL,
  description   TEXT,
  theme_color   CHAR(7)          NOT NULL DEFAULT '#FFFFFF',  -- '#RRGGBB'
  created_at    DATETIME(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at    DATETIME(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_oshis_user_name (user_id, name),
  KEY idx_oshis_user (user_id),
  CONSTRAINT fk_oshis_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 2-1) 推しの外部アカウント（URL複数）
CREATE TABLE oshi_accounts (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  oshi_id    BIGINT UNSIGNED NOT NULL,
  url        VARCHAR(2048)   NOT NULL,
  created_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_oshi_accounts (oshi_id, url(191)),
  KEY idx_oshi_accounts_oshi (oshi_id),
  CONSTRAINT fk_oshi_accounts_oshi FOREIGN KEY (oshi_id) REFERENCES oshis(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 3) カテゴリ（共通マスタ）
CREATE TABLE categories (
  id          SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name        VARCHAR(100) NOT NULL,
  description TEXT,
  created_at  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_categories_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 3-1) 推しごとの重要視カテゴリ
CREATE TABLE oshi_categories (
  oshi_id     BIGINT UNSIGNED    NOT NULL,
  category_id SMALLINT UNSIGNED  NOT NULL,
  created_at  DATETIME(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (oshi_id, category_id),
  KEY idx_oshi_categories_oshi (oshi_id),
  CONSTRAINT fk_oc_oshi FOREIGN KEY (oshi_id) REFERENCES oshis(id) ON DELETE CASCADE,
  CONSTRAINT fk_oc_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 4) イベント（カレンダー）
CREATE TABLE events (
  id           BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT,
  oshi_id      BIGINT UNSIGNED   NOT NULL,
  category_id  SMALLINT UNSIGNED          DEFAULT NULL,
  title        VARCHAR(255)      NOT NULL,
  description  TEXT,
  has_alarm    TINYINT(1)        NOT NULL DEFAULT 1,  -- 通知ON/OFF
  starts_at    DATETIME(3)       NOT NULL,
  ends_at      DATETIME(3)                DEFAULT NULL,
  created_at   DATETIME(3)       NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at   DATETIME(3)       NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  KEY idx_events_oshi (oshi_id),
  KEY idx_events_category (category_id),
  KEY idx_events_oshi_starts (oshi_id, starts_at),
  CONSTRAINT fk_events_oshi FOREIGN KEY (oshi_id) REFERENCES oshis(id) ON DELETE CASCADE,
  CONSTRAINT fk_events_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
  CONSTRAINT chk_events_time CHECK (ends_at IS NULL OR ends_at >= starts_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- 1. ユーザーデータ（2件）
INSERT INTO users (name, email, password_hash) VALUES 
('田中太郎', 'tanaka@example.com', '$2a$10$dummy.hash.for.password123'),
('佐藤花子', 'sato@example.com', '$2a$10$dummy.hash.for.password456');

-- 2. カテゴリデータ（8件）
INSERT INTO categories (name, description) VALUES 
('ライブ・コンサート', 'ライブ、コンサートなどのイベント情報'),
('グッズ・商品', 'グッズ発売、商品情報、コラボ商品など'),
('メディア出演', 'テレビ、ラジオ、雑誌、インタビューなど'),
('リリース', '新曲、アルバム、MV、配信など'),
('イベント・ファンミ', 'ファンミーティング、サイン会、トークショーなど'),
('ドラマ・映画', 'ドラマ出演、映画公開、舞台などの出演情報など'),
('SNS・配信', 'インスタライブ、YouTube、TikTokなど'),
('ニュース・発表', '重大発表、ニュース、お知らせなど');

-- 3. 推しデータ（2件）
INSERT INTO oshis (user_id, name, description, theme_color) VALUES 
(1, '山田美咲', 'アイドルグループABCのメンバー', '#FF69B4'),
(2, '鈴木愛', 'ソロシンガー', '#87CEEB');

-- 4. 推しの外部アカウント（2件）
INSERT INTO oshi_accounts (oshi_id, url) VALUES 
(1, 'https://twitter.com/yamada_misaki'),
(2, 'https://instagram.com/suzuki_ai_official');

-- 5. 推しごとの重要視カテゴリ（2件）
INSERT INTO oshi_categories (oshi_id, category_id) VALUES 
(1, 1), -- 山田美咲 → ライブ
(2, 2); -- 鈴木愛 → グッズ

-- 6. イベントデータ（2件）
INSERT INTO events (oshi_id, category_id, title, description, has_alarm, starts_at, ends_at) VALUES 
(1, 1, '山田美咲ソロライブ', '待望のソロライブ開催！', 1, '2024-12-25 18:00:00', '2024-12-25 21:00:00'),
(2, 2, 'もちどる', 'もちどる発売開始', 1, '2024-11-15 10:00:00', NULL);
