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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3) カテゴリ（共通マスタ）
CREATE TABLE categories (
  id          SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  slug        VARCHAR(50)  NOT NULL,
  name        VARCHAR(100) NOT NULL,
  description TEXT,
  created_at  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_categories_slug (slug),
  KEY idx_categories_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3-1) 推しごとの重要視カテゴリ
CREATE TABLE oshi_categories (
  oshi_id     BIGINT UNSIGNED    NOT NULL,
  category_id SMALLINT UNSIGNED  NOT NULL,
  created_at  DATETIME(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (oshi_id, category_id),
  KEY idx_oshi_categories_oshi (oshi_id),
  CONSTRAINT fk_oc_oshi FOREIGN KEY (oshi_id) REFERENCES oshis(id) ON DELETE CASCADE,
  CONSTRAINT fk_oc_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4) イベント（カレンダー）
CREATE TABLE events (
  id           BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT,
  oshi_id      BIGINT UNSIGNED   NOT NULL,
  category_id  SMALLINT UNSIGNED          DEFAULT NULL,
  post_id      BIGINT UNSIGNED            DEFAULT NULL,
  title        VARCHAR(255)      NOT NULL,
  description  TEXT,
  url          VARCHAR(2048)              DEFAULT NULL, -- イベントURL
  has_alarm    TINYINT(1)        NOT NULL DEFAULT 1,    -- 通知ON/OFF
  notification_timing ENUM('0', '5m', '10m', '15m', '30m', '1h', '2h', '1d', '2d', '1w') DEFAULT '15m',
  has_notification_sent   TINYINT(1)    DEFAULT 0,
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
