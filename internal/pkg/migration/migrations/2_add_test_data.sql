-- +goose Up
-- Test Users
INSERT INTO u_user (user_id, user_name, nickname, status, created_at, updated_at) VALUES
('11111111-1111-1111-1111-111111111111', 'user_a', 'Alice', 'Online', NOW(), NOW()),
('22222222-2222-2222-2222-222222222222', 'user_b', 'Bob', 'Offline', NOW(), NOW()),
('33333333-3333-3333-3333-333333333333', 'user_c', 'Charlie', 'Online', NOW(), NOW());

-- Test Channels
INSERT INTO u_channel (channel_id, channel_name, display_name, description, created_at, updated_at) VALUES
('44444444-4444-4444-4444-444444444444', 'general', 'General', 'For general announcements and discussions.', NOW(), NOW()),
('55555555-5555-5555-5555-555555555555', 'random', 'Random', 'A place for non-work-related chit-chat.', NOW(), NOW()),
('66666666-6666-6666-6666-666666666666', 'tech-talk', 'Tech Talk', 'Discussing technology, code, and everything in between.', NOW(), NOW());

-- Test Messages
-- Conversation in 'general' channel
INSERT INTO u_message (message_id, channel_id, user_id, content, created_at, updated_at) VALUES
('77777777-7777-7777-7777-777777777777', '44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'Hello everyone!', NOW() - INTERVAL 5 MINUTE, NOW() - INTERVAL 5 MINUTE),
('88888888-8888-8888-8888-888888888888', '44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 'Hi Alice, how are you?', NOW() - INTERVAL 4 MINUTE, NOW() - INTERVAL 4 MINUTE),
('99999999-9999-9999-9999-999999999999', '44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'I''m doing great, thanks! Just wanted to share the good news about our latest release.', NOW() - INTERVAL 3 MINUTE, NOW() - INTERVAL 3 MINUTE);

-- Conversation in 'random' channel
INSERT INTO u_message (message_id, channel_id, user_id, content, created_at, updated_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', '33333333-3333-3333-3333-333333333333', 'Does anyone have plans for the weekend?', NOW() - INTERVAL 10 MINUTE, NOW() - INTERVAL 10 MINUTE),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '55555555-5555-5555-5555-555555555555', '22222222-2222-2222-2222-222222222222', 'I''m thinking of going for a hike. The weather is supposed to be great.', NOW() - INTERVAL 9 MINUTE, NOW() - INTERVAL 9 MINUTE);

-- Conversation in 'tech-talk' channel
INSERT INTO u_message (message_id, channel_id, user_id, content, created_at, updated_at) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', '66666666-6666-6666-6666-666666666666', '11111111-1111-1111-1111-111111111111', 'I''ve been playing around with Go generics. They are pretty cool!', NOW() - INTERVAL 2 MINUTE, NOW() - INTERVAL 2 MINUTE),
('dddddddd-dddd-dddd-dddd-dddddddddddd', '66666666-6666-6666-6666-666666666666', '33333333-3333-3333-3333-333333333333', 'Oh nice! I haven''t had a chance to look at them yet. Any interesting findings?', NOW() - INTERVAL 1 MINUTE, NOW() - INTERVAL 1 MINUTE);

-- +goose Down
DELETE FROM u_message WHERE message_id IN (
'77777777-7777-7777-7777-777777777777',
'88888888-8888-8888-8888-888888888888',
'99999999-9999-9999-9999-999999999999',
'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
'cccccccc-cccc-cccc-cccc-cccccccccccc',
'dddddddd-dddd-dddd-dddd-dddddddddddd'
);

DELETE FROM u_channel WHERE channel_id IN (
'44444444-4444-4444-4444-444444444444',
'55555555-5555-5555-5555-555555555555',
'66666666-6666-6666-6666-666666666666'
);

DELETE FROM u_user WHERE user_id IN (
'11111111-1111-1111-1111-111111111111',
'22222222-2222-2222-2222-222222222222',
'33333333-3333-3333-3333-333333333333'
);
