-- +migrate Up
CREATE TABLE "public"."connected_accounts" (
    id uuid DEFAULT uuid_generate_v4 (),
    provider character varying NOT NULL,
    provider_user_id character varying NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

-- +migrate Down
DROP TABLE "public"."connected_accounts";