create unique index users_email_key on users(email) where archived_at is null;
create unique index users_username_key on users(username) where archived_at is null;
