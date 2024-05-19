CREATE TABLE "conversation"
(
    "conversation_id" bigserial NOT NULL,
    "audio_text"      TEXT      NULL,
    "audio_name"      TEXT      NOT NULL,
    "created_at"      DATE      NOT NULL,
    "good_percent"    INTEGER   NULL,
    "bad_percent"     INTEGER   NULL
);
ALTER TABLE
    "conversation"
    ADD PRIMARY KEY ("conversation_id");