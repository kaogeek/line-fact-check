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

-- Indexes for better performance
CREATE INDEX idx_topics_status ON topics(status);
CREATE INDEX idx_topics_created_at ON topics(created_at);
CREATE INDEX idx_user_messages_type ON user_messages(type);
CREATE INDEX idx_user_messages_created_at ON user_messages(created_at);
CREATE INDEX idx_messages_user_message_id ON messages(user_message_id);
CREATE INDEX idx_messages_topic_id ON messages(topic_id);
CREATE INDEX idx_messages_status ON messages(status);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_language ON messages(language); 