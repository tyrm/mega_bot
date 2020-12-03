-- +migrate Up
CREATE TABLE "public"."connected_accounts" (
    id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL,
    provider character varying NOT NULL,
    provider_user_id character varying NOT NULL,
    provider_username character varying,
    provider_avatar character varying,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
)
;

-- +migrate Down
DROP TABLE "public"."connected_accounts";