-- 1. ユーザーデータ（2件）
INSERT INTO users (name, email, password_hash) VALUES 
('田中太郎', 'tanaka@example.com', '$2a$10$dummy.hash.for.password123'),
('佐藤花子', 'sato@example.com', '$2a$10$dummy.hash.for.password456');

-- 2. カテゴリデータ（8件）
INSERT INTO categories (slug, name, description) VALUES 
('live', 'ライブ・コンサート', 'ライブ、コンサートなどのイベント情報'),
('goods', 'グッズ・商品', 'グッズ発売、商品情報、コラボ商品など'),
('media', 'メディア出演', 'テレビ、ラジオ、雑誌、インタビューなど'),
('release', 'リリース', '新曲、アルバム、MV、配信など'),
('event', 'イベント・ファンミーミング', 'ファンミーティング、サイン会、トークショーなど'),
('movie', 'ドラマ・映画', 'ドラマ出演、映画公開、舞台などの出演情報など'),
('sns', 'SNS・配信', 'インスタライブ、YouTube、TikTokなど'),
('news', 'ニュース・発表', '重大発表、ニュース、お知らせなど');

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

-- 6. イベントデータ（拡張版：通知機能とURL付き）
INSERT INTO events (oshi_id, category_id, title, description, url, has_alarm, notification_timing, has_notification_sent, starts_at, ends_at) VALUES 
-- 山田美咲のイベント
(1, 1, '山田美咲ソロライブ', '待望のソロライブ開催！', 'https://tickets.example.com/yamada-solo-live', 1, '1h', 0, '2024-12-25 18:00:00', '2024-12-25 21:00:00'),
(1, 3, 'テレビ出演「音楽番組SP」', '年末音楽番組にゲスト出演', 'https://tv.example.com/music-sp', 1, '30m', 0, '2024-12-23 20:00:00', '2024-12-23 21:00:00'),
(1, 2, 'ABCグループ 限定グッズ発売', 'ライブ会場限定グッズ', 'https://abc-goods.com/limited-2024', 1, '15m', 0, '2024-12-20 10:00:00', NULL),

-- 鈴木愛のイベント
(2, 2, '鈴木愛 写真集発売', '初写真集「Memories」発売開始', 'https://books.example.com/suzuki-ai-photobook', 1, '0', 0, '2024-11-15 10:00:00', NULL),
(2, 4, '新シングル「冬の歌」配信開始', 'デジタル配信限定リリース', 'https://music.example.com/suzuki-ai/winter-song', 1, '5m', 0, '2024-12-10 00:00:00', NULL),
(2, 7, 'YouTube Live配信', 'ファンとの質問コーナー', 'https://youtube.com/suzuki-ai-official', 1, '10m', 0, '2024-12-15 19:00:00', '2024-12-15 20:00:00'),
(2, 1, '鈴木愛 アコースティックライブ', 'カフェでの小規模ライブ', 'https://cafe-live.example.com/suzuki-ai', 1, '2h', 0, '2024-12-28 15:00:00', '2024-12-28 17:00:00');

-- 7. 追加ユーザーとその推しデータ
INSERT INTO users (name, email, password_hash) VALUES 
('山田花子', 'yamada@example.com', '$2a$10$dummy.hash.for.password789');

INSERT INTO oshis (user_id, name, description, theme_color) VALUES 
(3, '田中美咲', 'アイドルグループXYZのメンバー', '#FF1493'),
(3, '佐藤愛美', 'ソロアーティスト、シンガーソングライター', '#9370DB');

-- 追加推しのURL
INSERT INTO oshi_accounts (oshi_id, url) VALUES 
(3, 'https://twitter.com/tanaka_misaki_xyz'),
(3, 'https://instagram.com/tanaka_misaki_official'),
(4, 'https://twitter.com/sato_manami_music'),
(4, 'https://youtube.com/@sato_manami_channel');

-- 追加推しのカテゴリー
INSERT INTO oshi_categories (oshi_id, category_id) VALUES 
(3, 1), -- 田中美咲 → ライブ・コンサート
(3, 3), -- 田中美咲 → メディア出演
(4, 4), -- 佐藤愛美 → リリース
(4, 7); -- 佐藤愛美 → SNS・配信

-- 追加推しのイベント
INSERT INTO events (oshi_id, category_id, title, description, url, has_alarm, notification_timing, has_notification_sent, starts_at, ends_at) VALUES 
-- 田中美咲のイベント
(3, 1, '田中美咲 ソロライブ「SPARKLE」', '待望のソロライブ開催！新曲も披露予定', 'https://tickets.example.com/tanaka-sparkle-live', 1, '1h', 0, '2024-12-25 18:00:00', '2024-12-25 21:00:00'),
(3, 3, 'ラジオ出演「夜のトーク番組」', 'FM局の人気番組にゲスト出演', 'https://radio.example.com/night-talk', 1, '15m', 0, '2024-12-20 22:00:00', '2024-12-20 23:00:00'),
(3, 1, 'XYZグループ 年末ライブ', 'グループ全体での年末スペシャルライブ', 'https://xyz-group.com/year-end-live', 1, '2h', 0, '2024-12-31 19:00:00', '2024-12-31 22:00:00'),

-- 佐藤愛美のイベント
(4, 4, '新シングル「星空の約束」リリース', 'デジタル配信開始', 'https://music.example.com/sato-manami/hoshizora', 1, '0', 0, '2024-12-15 00:00:00', NULL),
(4, 7, 'Instagram Live配信', 'ファンとの交流配信', 'https://instagram.com/sato_manami_official', 1, '5m', 0, '2024-12-18 20:00:00', '2024-12-18 21:00:00'),
(4, 1, '佐藤愛美 アコースティックライブ', 'intimate なアコースティックライブ', 'https://tickets.example.com/sato-acoustic', 1, '1d', 0, '2024-12-22 15:00:00', '2024-12-22 17:00:00');