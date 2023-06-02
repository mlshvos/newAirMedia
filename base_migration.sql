CREATE TABLE accounts (
    account_uuid uuid NOT NULL PRIMARY KEY,
    money float8 NOT NULL DEFAULT 0,
    create_time timestamp with time zone DEFAULT now() NOT NULL,
    update_time timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TYPE status_type AS ENUM ('failed', 'in_progress', 'success');

CREATE TABLE requests (
    request_id uuid NOT NULL PRIMARY KEY,
    account_uuid uuid NOT NULL,
    money_transfer float8 NOT NULL,
    status status_type NOT NULL DEFAULT 'in_progress',
    create_time timestamp with time zone DEFAULT now() NOT NULL,
    update_time timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(account_uuid) REFERENCES accounts(account_uuid)
);

CREATE INDEX requests_idx ON requests(account_uuid);