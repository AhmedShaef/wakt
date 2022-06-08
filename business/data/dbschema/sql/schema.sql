-- Version: 1.0
-- Description: Create table users
CREATE TABLE users
(
    user_id            uuid
        constraint users_pk primary key,
    default_wid        uuid,
    email              text UNIQUE,
    password_hash      text,
    full_name          text,
    time_of_day_format text,
    date_format        text,
    beginning_of_week  integer,
    language           text,
    image_url          text,
    date_created       timestamp,
    date_updated       timestamp,
    timezone           text,
    invitation         text[],
    duration_format    text
);
-- Description: Create table workspaces
CREATE TABLE workspaces
(
    workspace_id                   uuid
        constraint workspaces_pk primary key,
    name                           text,
    uid                            uuid,
    constraint workspaces_uid_fk foreign key (uid) references users (user_id),
    default_hourly_rate            double precision,
    default_currency               text,
    only_admin_may_create_projects boolean,
    only_admin_see_billable_rates  boolean,
    only_admin_see_team_dashboard  boolean,
    rounding                       integer,
    rounding_minutes               integer,
    date_created                   timestamp,
    date_updated                   timestamp,
    logo_url                       text
);

-- Description: Create table clients
CREATE TABLE clients
(
    client_id    uuid
        constraint client_pk primary key,
    name         text,
    uid          uuid,
    constraint client_uid_fk foreign key (uid) references users (user_id),
    wid          uuid,
    constraint client_wid_fk foreign key (wid) references workspaces (workspace_id),
    notes        text,
    date_created timestamp,
    date_updated timestamp
);
-- Description: Create table projects
CREATE TABLE projects
(
    project_id      uuid
        constraint project_pk primary key,
    name            text,
    wid             uuid,
    constraint project_wid_fk foreign key (wid) references workspaces (workspace_id),
    cid             uuid,
    constraint project_cid_fk foreign key (cid) references clients (client_id),
    uid             uuid,
    constraint project_uid_fk foreign key (uid) references users (user_id),
    active          boolean,
    is_private      boolean,
    billable        boolean,
    auto_estimates  boolean,
    estimated_hours double precision,
    date_created    timestamp,
    date_updated    timestamp,
    rate            double precision,
    hex_color       text
);
-- Description: Create table tasks
CREATE TABLE tasks
(
    task_id           uuid
        constraint task_pk primary key,
    name              text,
    pid               uuid,
    constraint task_pid_fk foreign key (pid) references projects (project_id),
    wid               uuid,
    constraint task_wid_fk foreign key (wid) references workspaces (workspace_id),
    uid               uuid,
    constraint task_uid_fk foreign key (uid) references users (user_id),
    estimated_seconds integer,
    active            boolean,
    date_created      timestamp,
    date_updated      timestamp,
    tracked_seconds   integer
);

-- Description: Create table time_entries
CREATE TABLE time_entries
(
    time_entry_id UUID
        constraint time_entry_pk primary key,
    description   TEXT,
    uid           UUID,
    constraint time_entry_user_id_fk foreign key (uid) references users (user_id),
    wid           UUID,
    constraint time_entry_wid_fk foreign key (wid) references workspaces (workspace_id),
    pid           UUID,
    constraint time_entry_pid_fk foreign key (pid) references projects (project_id),
    tid           UUID,
    constraint time_entry_tid_fk foreign key (tid) references tasks (task_id),
    billable      BOOLEAN,
    start         TIMESTAMP,
    stop          TIMESTAMP,
    duration      INTEGER,
    created_with  TEXT,
    tags          TEXT[],
    dur_only      BOOLEAN,
    date_created  TIMESTAMP,
    date_updated  TIMESTAMP
);


-- Description: Create table groups
CREATE TABLE groups
(
    group_id     uuid
        constraint group_pk primary key,
    name         text,
    wid          uuid,
    constraint group_wid_fk foreign key (wid) references workspaces (workspace_id),
    date_created timestamp,
    date_updated timestamp
);

-- Description: Create table tags
CREATE TABLE tags
(
    tag_id       uuid
        constraint tags_pk primary key,
    name         text,
    wid          uuid,
    constraint tags_wid_fk foreign key (wid) references workspaces (workspace_id),
    date_created timestamp,
    date_updated timestamp
);

-- Description: Create table project_users
CREATE TABLE project_users
(
    project_user_id uuid
        constraint project_user_pk primary key,
    pid             uuid,
    constraint project_user_pid_fk foreign key (pid) references projects (project_id),
    uid             uuid,
    constraint project_user_uid_fk foreign key (uid) references users (user_id),
    wid             uuid,
    constraint project_user_wid_fk foreign key (wid) references workspaces (workspace_id),
    manager         boolean,
    date_created    timestamp,
    date_updated    timestamp
);

-- Description: Create table workspace_users
CREATE TABLE workspace_users
(
    workspace_user_id uuid
        constraint workspace_user_pk primary key,
    uid               uuid,
    constraint workspace_user_uid_fk foreign key (uid) references users (user_id),
    wid               uuid,
    constraint workspace_user_wid_fk foreign key (wid) references workspaces (workspace_id),
    admin             boolean,
    active            boolean,
    invite_key       text,
    date_created      timestamp,
    date_updated      timestamp
);