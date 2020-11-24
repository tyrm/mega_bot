-- +migrate Up
CREATE TABLE "public"."users" (
    id uuid DEFAULT uuid_generate_v4 (),
    email character varying NOT NULL,
    password character varying,
    nick character varying,
    admin boolean NOT NULL DEFAULT false,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

-- +migrate Down
DROP TABLE "public"."users";