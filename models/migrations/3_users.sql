-- +migrate Up
CREATE TABLE "public"."users" (
    id uuid DEFAULT uuid_generate_v4 (),
    email character varying NOT NULL UNIQUE,
    password character varying,
    nick character varying,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

-- +migrate Down
DROP TABLE "public"."users";