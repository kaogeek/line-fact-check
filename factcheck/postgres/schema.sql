-- Topics table
CREATE TABLE topics (
    id            UUID NOT NULL PRIMARY KEY,
    name          text NOT NULL,
    status        text NOT NULL,
    result        text,
    result_status text,
    created_at    timestamptz NOT NULL,
    updated_at    timestamptz
);

-- Messages table
CREATE TABLE messages (
    id         UUID NOT NULL PRIMARY KEY,
    topic_id   UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    text       text NOT NULL,
    type       text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);

-- User messages table (generic table for UserMessage struct)
CREATE TABLE user_messages (
    id         UUID NOT NULL PRIMARY KEY,
    replied_at timestamptz,
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    metadata   jsonb, -- Using jsonb for generic metadata
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);

-- Indexes for better performance
CREATE INDEX idx_topics_status ON topics(status);
CREATE INDEX idx_topics_created_at ON topics(created_at);
CREATE INDEX idx_messages_topic_id ON messages(topic_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_user_messages_message_id ON user_messages(message_id);
CREATE INDEX idx_user_messages_created_at ON user_messages(created_at); 