-- Version: 1.1
-- Description: Create table workspace
CREATE TABLE workspaces (
    workspace_id                   uuid constraint workspaces_pk primary key,
    name                           text,
    profile                        integer          default 0,
    premium                        boolean          default false,
    admin                          boolean          default true,
    default_hourly_rate            double precision default 0,
    default_currency               text             default 'USD',
    only_admin_may_create_projects boolean          default false,
    only_admin_see_billable_rates  boolean          default false,
    only_admin_see_team_dashboard  boolean          default false,
    project_billable_by_default    boolean          default true,
    rounding                       integer,
    rounding_minutes               integer,
    date_updated                   timestamp,
    logo_url                       text,
    ical_url                       text,
    ical_enabled                   boolean          default true
);

-- Description: Create table client
CREATE TABLE clients(
    client_id uuid constraint client_pk primary key,
    name      text,
    wid       uuid,
    constraint client_wid_fk foreign key (wid) references workspaces(workspace_id),
    notes     text,
    date_updated timestamp
);
-- Description: Create table project
CREATE TABLE projects(
     project_id uuid constraint project_pk primary key,
     name       text,
     wid        uuid,
     constraint project_wid_fk foreign key (wid) references workspaces(workspace_id),
     cid        uuid,
     constraint project_cid_fk foreign key (cid) references clients(client_id),
     active     boolean default true,
     is_private boolean default true,
     template   boolean default false,
     template_id uuid,
     billable   boolean default true,
     auto_estimates boolean default false,
     estimated_hours double precision,
     date_updated timestamp,
     color text,
     rate double precision default 0,
     date_created timestamp
);

-- Description: Create table user
CREATE TABLE users (
   user_id                   uuid   constraint users_pk primary key,
   api_token                 text,
   default_wid               uuid,
   email                     text,
   password_hash             text,
   roles    text[],
   full_name                 text,
   jquery_time_of_day_format text,
   jquery_date_format        text,
   time_of_day_format        text,
   date_format               text,
   store_start_and_stop_time boolean default false,
   begining_of_week          integer default 1,
   language                  text    default 'en_US',
   image_url                 text,
   sidebar_piechart          boolean default false,
   date_created              timestamp,
   date_updated              timestamp,
   record_timeline           boolean default false,
   should_upgrade            boolean default false,
   send_product_emails       boolean default true,
   send_weekly_report        boolean default true,
   send_timer_notification   boolean default true,
   openid_enabled            boolean default false,
   timezone                  text,
   invitation                text[],
   duration_format           text
);

-- Description: Create table task
CREATE TABLE tasks(
  task_id uuid constraint task_pk primary key,
  name    text,
  pid     uuid,
  constraint task_pid_fk foreign key (pid) references projects(project_id),
  wid     uuid,
  constraint task_wid_fk foreign key (wid) references workspaces(workspace_id),
  uid    uuid,
  constraint task_uid_fk foreign key (uid) references users(user_id),
  estimated_seconds integer default 0,
  active  boolean default true,
  date_updated timestamp,
  tracked_seconds integer default 0
);

-- Description: Create table user
CREATE TABLE time_entries(
     time_entrie_id uuid constraint time_entries_pk primary key,
     description     text       default '',
     uid             uuid,
     constraint time_entries_user_id_fk foreign key (uid) references users (user_id),
     wid             uuid,
     constraint time_entries_wid_fk foreign key (wid) references workspaces(workspace_id),
     pid             uuid,
     constraint time_entries_pid_fk foreign key (pid) references projects(project_id),
     tid             uuid,
     constraint time_entries_tid_fk foreign key (tid) references tasks(task_id),
     billable        boolean    default false,
     start           timestamp,
     stop            timestamp,
     duration        integer    default -1,
     created_with    text,
     tags            text[],
     dur_only        boolean    default true,
     date_created    timestamp,
     date_updated    timestamp
);



-- Description: Create table group
CREATE TABLE groups(
    group_id uuid constraint group_pk primary key,
    name     text,
    wid      uuid,
    constraint group_wid_fk foreign key (wid) references workspaces(workspace_id),
    date_updated timestamp
);

-- Description: Create table tags
CREATE TABLE tags(
     tag_id uuid constraint tags_pk primary key,
     name   text,
     wid    uuid,
     constraint tags_wid_fk foreign key (wid) references workspaces(workspace_id)
);

-- Description: Create table project_user
CREATE TABLE project_users(
    project_user_id uuid constraint project_user_pk primary key,
    pid      uuid,
    constraint project_user_pid_fk foreign key (pid) references projects(project_id),
    uid         uuid,
    constraint project_user_uid_fk foreign key (uid) references users(user_id),
    wid         uuid,
    constraint project_user_wid_fk foreign key (wid) references workspaces(workspace_id),
    manager     boolean default false,
    rate    double precision default 0,
    date_updated    timestamp
);

-- Description: Create table workspace_user
CREATE TABLE workspace_users(
    workspace_user_id uuid constraint workspace_user_pk primary key,
    uid         uuid,
    constraint workspace_user_uid_fk foreign key (uid) references users(user_id),
    admin       boolean default false,
    active      boolean default true,
    email    text[],
    invite_url  text
);