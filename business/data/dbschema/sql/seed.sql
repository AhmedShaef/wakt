INSERT INTO users (user_id, default_wid, full_name, email, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7','7da3ca14-6366-47cf-b953-f706226567d8', 'Admin Gopher', 'admin@example.com', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f','6fa2132c-9bdd-428a-b025-5f1a4d6ee683', 'User Gopher', 'user@example.com', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;
INSERT INTO workspaces (workspace_id, name, uid, date_updated) VALUES
    ('7da3ca14-6366-47cf-b953-f706226567d8', 'Default Workspace', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00'),
    ('6fa2132c-9bdd-428a-b025-5f1a4d6ee683', 'User Workspace', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO clients (client_id, name, uid, wid, notes, date_updated) VALUES
    ('c78db68e-e004-44f5-895b-ba562dc53d9d', 'Default Client', '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8','note1', '2019-03-24 00:00:00'),
    ('a9c8488a-5df2-40c5-8e76-4ac1670e7ac7', 'User Client', '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8','note2', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO projects (project_id, name, uid, wid, cid, date_created,date_updated) VALUES
    ('45cf87a3-5915-4079-a9af-6c559239ddbf', 'Default Project', '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8', 'c78db68e-e004-44f5-895b-ba562dc53d9d', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
    ('d774cc57-e4a6-4be2-bca1-cb50610fb3f5', 'User Project', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '7da3ca14-6366-47cf-b953-f706226567d8', 'a9c8488a-5df2-40c5-8e76-4ac1670e7ac7', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO tasks (task_id, name, pid, wid, uid, date_updated) VALUES
    ('346efd40-6d6e-46d5-b60b-5db9fc171779', 'Default Task', '45cf87a3-5915-4079-a9af-6c559239ddbf', '7da3ca14-6366-47cf-b953-f706226567d8', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00'),
    ('4ea20d73-a11e-4e83-b95c-ba8b4b5ff6c1', 'User Task', 'd774cc57-e4a6-4be2-bca1-cb50610fb3f5', '6fa2132c-9bdd-428a-b025-5f1a4d6ee683', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO time_entries(time_entry_id, description, uid,wid, pid, tid, start,stop,duration,created_with,tags,date_created,date_updated) VALUES
    ('57a785f7-aff5-40a6-8b98-fc28e0f0465c', 'Default Time Entry', '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8', '45cf87a3-5915-4079-a9af-6c559239ddbf', '346efd40-6d6e-46d5-b60b-5db9fc171779', '2019-03-24 00:00:00', '2019-03-25 00:00:00', '86400', 'curl', '{tagsy}', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
    ('3d4d8f5e-b776-4481-8664-265de2a07669', 'User Time Entry', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '6fa2132c-9bdd-428a-b025-5f1a4d6ee683', 'd774cc57-e4a6-4be2-bca1-cb50610fb3f5', '4ea20d73-a11e-4e83-b95c-ba8b4b5ff6c1', '2019-03-24 00:00:00', '2019-03-26 00:00:00', '172800', 'curl', '{tagso}', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO groups(group_id, name, wid, uid, date_updated) values
    ('ee8a891a-6e2e-4fa3-8f01-d4e559dd5a72', 'Default Group', '7da3ca14-6366-47cf-b953-f706226567d8', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00'),
    ('8e95bc44-17dc-4006-961e-1a0bec9ea943', 'User Group', '7da3ca14-6366-47cf-b953-f706226567d8', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO tags(tag_id, name, wid, uid) VALUES
    ('82cf01da-4a6c-40fc-98cf-c9987aca40b2', 'Default tag', '7da3ca14-6366-47cf-b953-f706226567d8', '5cf37266-3473-4006-984f-9325122678b7'),
    ('4c698e19-8cb9-47de-a496-1f8ef724218d','User tag', '7da3ca14-6366-47cf-b953-f706226567d8', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f')
    ON CONFLICT DO NOTHING;
INSERT INTO project_users(project_user_id, pid, uid, wid, date_updated) values
    ('efcc74aa-86d2-4e11-80f9-3ca912af8269', '45cf87a3-5915-4079-a9af-6c559239ddbf', '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8', '2019-03-24 00:00:00'),
    ('c7142720-91d3-4d1e-841d-680042b6500c', 'd774cc57-e4a6-4be2-bca1-cb50610fb3f5', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '6fa2132c-9bdd-428a-b025-5f1a4d6ee683', '2019-03-24 00:00:00')
    ON CONFLICT DO NOTHING;
INSERT INTO workspace_users(workspace_user_id, uid, wid) values
    ('32c1494f-1c1f-4981-857f-b0526cb654ec', '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8'),
    ('604125e7-f368-4ff0-8170-dfd2f428510a', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '6fa2132c-9bdd-428a-b025-5f1a4d6ee683')
    ON CONFLICT DO NOTHING;