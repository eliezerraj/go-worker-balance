# go-worker-balance

POC for test purposes.

CRUD a balance data

## Database

        CREATE TABLE balance_cdc (
            id              SERIAL PRIMARY KEY,
            account_id      varchar(200) UNIQUE NULL,
            person_id       varchar(200) NULL,
            currency        varchar(10) NULL,   
            amount          float8 NULL,
            create_at       timestamptz NULL,
            update_at       timestamptz NULL,
            tenant_id       varchar(200) null,
            user_last_update	varchar(200) NULL);