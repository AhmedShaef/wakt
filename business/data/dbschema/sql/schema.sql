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