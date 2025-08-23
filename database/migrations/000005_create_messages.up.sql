CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  -- default must be text, but we will also support image, audio, video, etc
  -- the restriction will be done in the code go side
  message_type VARCHAR(255) NOT NULL,
  -- must be SENT, DELIVERED, READ, FAILED
  status VARCHAR(255) NOT NULL,
  reply_to_message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
  media_url VARCHAR(500), -- for file/image/audio/video messages
  media_filename VARCHAR(255), -- original filename for files
  media_size BIGINT, -- size in bytes for files
  media_mime_type VARCHAR(255), -- MIME type for files
  location_latitude DECIMAL(10, 8), -- for location messages
  location_longitude DECIMAL(11, 8), -- for location messages
  location_address TEXT, -- human readable address
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  edited_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE,
  delivered_at TIMESTAMP WITH TIME ZONE,
  read_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_reply_to ON messages(reply_to_message_id);
CREATE INDEX IF NOT EXISTS idx_messages_type ON messages(message_type);
CREATE INDEX IF NOT EXISTS idx_messages_not_deleted ON messages(conversation_id, is_deleted, created_at DESC);