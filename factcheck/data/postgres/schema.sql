BEGIN;

-- Topics table
CREATE TABLE topics (
    id            UUID NOT NULL PRIMARY KEY,
    name          text NOT NULL,
    description   text NOT NULL,
    status        text NOT NULL,
    result        text,
    result_status text,
    created_at    timestamptz NOT NULL,
    updated_at    timestamptz
);

-- User messages table (generic table for UserMessage struct)
CREATE TABLE user_messages (
    id         UUID NOT NULL PRIMARY KEY,
    type       text NOT NULL,
    replied_at timestamptz,
    metadata   jsonb, -- Using jsonb for generic metadata
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);

-- Messages table
CREATE TABLE messages (
    id             UUID NOT NULL PRIMARY KEY,
    user_message_id UUID NOT NULL REFERENCES user_messages(id) ON DELETE CASCADE,
    type           text NOT NULL,
    status         text NOT NULL,
    topic_id       UUID REFERENCES topics(id) ON DELETE SET NULL,
    text           text NOT NULL,
    language       text,
    created_at     timestamptz NOT NULL,
    updated_at     timestamptz
);

-- MessageGroup table (groups messages with identical text)
CREATE TABLE message_groups (
    id         UUID NOT NULL PRIMARY KEY,
    topic_id   UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    name       text,
    text       text,
    text_sha1  text,
    created_at timestamptz,
    updated_at timestamptz NOT NULL
);

-- MessagesV2 table (replaces messages + user_messages relationship)
CREATE TABLE messages_v2 (
    id         UUID NOT NULL PRIMARY KEY,
    user_id    text NOT NULL,
    topic_id   UUID REFERENCES topics(id) ON DELETE SET NULL,
    group_id   UUID REFERENCES message_groups(id) ON DELETE SET NULL,
    type_user  text NOT NULL,
    type       text NOT NULL,
    text       text NOT NULL,
    language   text,
    metadata   jsonb,
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);

CREATE INDEX idx_topics_status ON topics(status);
CREATE INDEX idx_topics_created_at ON topics(created_at);
CREATE INDEX idx_user_messages_type ON user_messages(type);
CREATE INDEX idx_user_messages_created_at ON user_messages(created_at);
CREATE INDEX idx_messages_user_message_id ON messages(user_message_id);
CREATE INDEX idx_messages_topic_id ON messages(topic_id);
CREATE INDEX idx_messages_status ON messages(status);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_language ON messages(language);
CREATE INDEX idx_messages_v2_user_id ON messages_v2(user_id);
CREATE INDEX idx_messages_v2_topic_id ON messages_v2(topic_id);
CREATE INDEX idx_messages_v2_group_id ON messages_v2(group_id);
CREATE INDEX idx_messages_v2_type_user ON messages_v2(type_user);
CREATE INDEX idx_messages_v2_type ON messages_v2(type);
CREATE INDEX idx_messages_v2_created_at ON messages_v2(created_at);
CREATE INDEX idx_messages_v2_language ON messages_v2(language);
CREATE INDEX idx_message_groups_topic_id ON message_groups(topic_id);
CREATE INDEX idx_message_groups_text_sha1 ON message_groups(text_sha1);
CREATE INDEX idx_message_groups_created_at ON message_groups(created_at);

COMMIT; 
