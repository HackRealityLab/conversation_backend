CREATE TABLE "conversation"
(
    "conversation_id" bigserial NOT NULL,
    "text"            TEXT      NULL,
    "audio_name"      TEXT      NOT NULL,
    "created_at"      DATE      NOT NULL,
    "errors_cnt"      BIGINT    NOT NULL,
    "is_ok"           BOOLEAN   NOT NULL
);
ALTER TABLE
    "conversation"
    ADD PRIMARY KEY ("conversation_id");