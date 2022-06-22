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
    wid          uuid,
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
    cid             uuid,
    uid             uuid,
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
    wid               uuid,
    uid               uuid,
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
    wid           UUID,
    pid           UUID,
    tid           UUID,
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
    uid          uuid,
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
    date_created timestamp,
    date_updated timestamp
);

-- Description: Create table team
CREATE TABLE teams
(
    team_id uuid
        constraint team_pk primary key,
    pid             uuid,
    uid             uuid,
    wid             uuid,
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
    wid               uuid,
    admin             boolean,
    active            boolean,
    invite_key       text,
    date_created      timestamp,
    date_updated      timestamp
);