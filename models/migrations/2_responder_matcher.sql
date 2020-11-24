-- +migrate Up
CREATE TABLE "public"."responder_matchers" (
    id uuid DEFAULT uuid_generate_v4 (),
    always_respond boolean NOT NULL,
    enabled boolean NOT NULL,
    matcher_re character varying NOT NULL,
    repsonse character varying NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

INSERT INTO public.responder_matchers(
	always_respond, enabled, matcher_re, repsonse)
	VALUES (true, true, '[uU][wW][uU]', ':skull: owo')
;

-- +migrate Down
DROP TABLE "public"."responder_matchers";