CREATE TABLE "conversation"
(
    "conversation_id" bigserial        NOT NULL,
    "text"            TEXT             NULL,
    "audio_name"      TEXT             NOT NULL,
    "created_at"      DATE             NOT NULL,
    "good_percent"    DOUBLE PRECISION NOT NULL,
    "bad_percent"     DOUBLE PRECISION NOT NULL
);
ALTER TABLE
    "conversation"
    ADD PRIMARY KEY ("conversation_id");