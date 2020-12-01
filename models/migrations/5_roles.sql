-- +migrate Up
CREATE TABLE "public"."roles" (
    id uuid DEFAULT uuid_generate_v4 (),
    name character varying NOT NULL UNIQUE,
    description character varying,
    system_role bool NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

CREATE TABLE "public"."user_roles" (
    user_id uuid,
    role_id uuid,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    unique (user_id, role_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE RESTRICT
)
;

INSERT INTO public.roles(
    name, description, system_role)
    VALUES ('administrator', 'MegaBot Administrator. Full access.', true),
           ('operator', 'MegaBot Operator. Access to all components and commands except administration.', true),
           ('authorized', 'Allowed to use MegaBot.', true)
;

-- +migrate Down
DROP TABLE "public"."user_roles";
DROP TABLE "public"."roles";