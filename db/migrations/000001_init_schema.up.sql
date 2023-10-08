CREATE TABLE "accounts" (
    "id"         bigserial primary key,
    "owner"      varchar(255) not null,
    "balance"    bigint       not null,
    "currency"   varchar(10)  not null,
    "created_at" timestamptz   not null default (now())
);

CREATE TABLE "entries" (
    "id" bigserial primary key,
    "account_id" bigint not null,
    "amount" bigint not null,
    "created_at" timestamptz   not null default (now())
);

CREATE TABLE "transfers" (
    "id" bigserial primary key,
    "from_account_id" bigint not null,
    "to_account_id" bigint not null,
    "amount" bigint not null,
    "created_at" timestamptz   not null default (now())
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

