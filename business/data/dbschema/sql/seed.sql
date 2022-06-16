INSERT INTO users(user_id, default_wid, email, password_hash, full_name, time_of_day_format, date_format,
                  beginning_of_week, language, image_url, date_created, date_updated, timezone, invitation,
                  duration_format)
VALUES ('5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8', 'admin@example.com',
        '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', 'Admin Gopher', 'hh:mm', 'DD/MM/YYYY', '1',
        'en_US', '', '2019-03-24 00:00:00', '2019-03-24 00:00:00', 'UTC', '{}', ''),
       ('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '6fa2132c-9bdd-428a-b025-5f1a4d6ee683', 'user@example.com',
        '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', 'User Gopher', 'hh:mm', 'DD/MM/YYYY', '1',
        'en_US', '', '2019-03-24 00:00:00', '2019-03-24 00:00:00', 'UTC', '{}', '')
ON CONFLICT DO NOTHING;
INSERT INTO workspaces(workspace_id, name, uid, default_hourly_rate, default_currency, only_admin_may_create_projects,
                       only_admin_see_billable_rates, only_admin_see_team_dashboard, rounding, rounding_minutes,
                       date_created, date_updated, logo_url)
VALUES ('7da3ca14-6366-47cf-b953-f706226567d8', 'Default Workspace', '5cf37266-3473-4006-984f-9325122678b7', '50.0',
        'USD', 'false', 'false', 'false', '1', '60', '2019-03-24 00:00:00', '2019-03-24 00:00:00', ''),
       ('6fa2132c-9bdd-428a-b025-5f1a4d6ee683', 'User Workspace', '5cf37266-3473-4006-984f-9325122678b7', '50.0', 'USD',
        'false', 'false', 'false', '1', '60', '2019-03-24 00:00:00', '2019-03-24 00:00:00', '')
ON CONFLICT DO NOTHING;
INSERT INTO clients(client_id, name, uid, wid, notes, date_created, date_updated)
VALUES ('c78db68e-e004-44f5-895b-ba562dc53d9d', 'Default Client', '5cf37266-3473-4006-984f-9325122678b7',
        '7da3ca14-6366-47cf-b953-f706226567d8', 'note1', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
       ('a9c8488a-5df2-40c5-8e76-4ac1670e7ac7', 'User Client', '5cf37266-3473-4006-984f-9325122678b7',
        '7da3ca14-6366-47cf-b953-f706226567d8', 'note2', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
INSERT INTO projects (project_id, name, wid, cid, uid, active, is_private, billable, auto_estimates, estimated_hours,
                      date_created, date_updated, rate, hex_color)
VALUES ('45cf87a3-5915-4079-a9af-6c559239ddbf', 'Default Project', '7da3ca14-6366-47cf-b953-f706226567d8',
        'c78db68e-e004-44f5-895b-ba562dc53d9d','5cf37266-3473-4006-984f-9325122678b7', 'false', 'false', 'true', 'false', '30.0', '2019-03-24 00:00:00',
        '2019-03-24 00:00:00', '30', '#ffffff'),
       ('d774cc57-e4a6-4be2-bca1-cb50610fb3f5', 'User Project', '7da3ca14-6366-47cf-b953-f706226567d8',
        'a9c8488a-5df2-40c5-8e76-4ac1670e7ac7','5cf37266-3473-4006-984f-9325122678b7', 'false', 'false', 'true', 'true', '0.0', '2019-03-24 00:00:00',
        '2019-03-24 00:00:00', '30', '#ffffff')
ON CONFLICT DO NOTHING;
INSERT INTO tasks (task_id, name, pid, wid, uid, estimated_seconds, active, date_created, date_updated, tracked_seconds)
VALUES ('346efd40-6d6e-46d5-b60b-5db9fc171779', 'Default Task', '45cf87a3-5915-4079-a9af-6c559239ddbf',
        '7da3ca14-6366-47cf-b953-f706226567d8', '5cf37266-3473-4006-984f-9325122678b7', '0', 'true',
        '2019-03-24 00:00:00', '2019-03-24 00:00:00', '0'),
       ('4ea20d73-a11e-4e83-b95c-ba8b4b5ff6c1', 'User Task', '45cf87a3-5915-4079-a9af-6c559239ddbf',
        '7da3ca14-6366-47cf-b953-f706226567d8', '5cf37266-3473-4006-984f-9325122678b7', '0', 'true',
        '2019-03-24 00:00:00', '2019-03-24 00:00:00', '0')
ON CONFLICT DO NOTHING;
INSERT INTO time_entries (time_entry_id, description, uid, wid, pid, tid, billable, start, stop, duration, created_with,
                          tags, dur_only, date_created, date_updated)
VALUES ('57a785f7-aff5-40a6-8b98-fc28e0f0465c', 'Default Time Entry', '5cf37266-3473-4006-984f-9325122678b7',
        '7da3ca14-6366-47cf-b953-f706226567d8', '45cf87a3-5915-4079-a9af-6c559239ddbf',
        '346efd40-6d6e-46d5-b60b-5db9fc171779', 'true', '2019-03-24 00:00:00', '2019-03-24 00:00:00', '-1', 'curl',
        '{tag1,tag2}', 'false', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
       ('3d4d8f5e-b776-4481-8664-265de2a07669', 'User Time Entry', '5cf37266-3473-4006-984f-9325122678b7',
        '6fa2132c-9bdd-428a-b025-5f1a4d6ee683', 'd774cc57-e4a6-4be2-bca1-cb50610fb3f5',
        '4ea20d73-a11e-4e83-b95c-ba8b4b5ff6c1', 'true', '2019-03-24 00:00:00', '2019-03-24 00:00:30', '-1', 'curl',
        '{tags1,tags2}', 'true', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
INSERT INTO groups (group_id, name, wid, uid, date_created, date_updated)
values ('ee8a891a-6e2e-4fa3-8f01-d4e559dd5a72', 'Default Group', '7da3ca14-6366-47cf-b953-f706226567d8',
        '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
       ('8e95bc44-17dc-4006-961e-1a0bec9ea943', 'User Group', '7da3ca14-6366-47cf-b953-f706226567d8',
        '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
INSERT INTO tags (tag_id, name, wid, date_created, date_updated)
VALUES ('82cf01da-4a6c-40fc-98cf-c9987aca40b2', 'Default tag', '7da3ca14-6366-47cf-b953-f706226567d8',
        '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
       ('4c698e19-8cb9-47de-a496-1f8ef724218d', 'User tag', '7da3ca14-6366-47cf-b953-f706226567d8',
        '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
INSERT INTO project_users (project_user_id, pid, uid, wid, manager, date_created, date_updated)
values ('efcc74aa-86d2-4e11-80f9-3ca912af8269', '45cf87a3-5915-4079-a9af-6c559239ddbf',
        '5cf37266-3473-4006-984f-9325122678b7', '7da3ca14-6366-47cf-b953-f706226567d8', 'true', '2019-03-24 00:00:00',
        '2019-03-24 00:00:00'),
       ('c7142720-91d3-4d1e-841d-680042b6500c', 'd774cc57-e4a6-4be2-bca1-cb50610fb3f5',
        '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '7da3ca14-6366-47cf-b953-f706226567d8', 'false', '2019-03-24 00:00:00',
        '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
INSERT INTO workspace_users (workspace_user_id, uid, wid, admin, active, invite_key, date_created, date_updated)
values ('32c1494f-1c1f-4981-857f-b0526cb654ec', '5cf37266-3473-4006-984f-9325122678b7',
        '7da3ca14-6366-47cf-b953-f706226567d8', 'true', 'false', '', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
       ('604125e7-f368-4ff0-8170-dfd2f428510a', '5cf37266-3473-4006-984f-9325122678b7',
        '7da3ca14-6366-47cf-b953-f706226567d8', 'false', 'true', '', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;