##
## Blocks
##

create table pvm_blocks
(
    id            varchar(50)      not null primary key,
    type          smallint         not null,
    parent_id     varchar(50)      not null,
    chain_id      varchar(50)      not null,
    serialization varbinary(32000) not null,
    created_at    timestamp        not null default current_timestamp
);

##
## Transactions
##

create table pvm_transactions
(
    id         varchar(50)     not null primary key,
    type       smallint        not null,
    block_id   varchar(50)     not null,
    nonce      bigint unsigned not null,
    signature  binary(65),
    created_at timestamp       not null default current_timestamp
);


##
## Allychains
##

create table pvm_allychains
(
    id         varchar(50)  not null primary key,
    network_id int unsigned not null,
    chain_id   varchar(50)  not null,
    threshold  int unsigned not null,
    created_at timestamp    not null default current_timestamp
);

create table pvm_allychain_control_keys
(
    allychain_id  varchar(50) not null,
    address    varchar(50) not null,
    public_key binary(33)  null

);
create unique index pvm_allychain_control_keys_allychain_id_address_idx ON pvm_allychain_control_keys (allychain_id, address);

##
## Validators
##

create table pvm_validators
(
    transaction_id varchar(50)     not null,

    # Validator
    node_id        varchar(50)     not null,
    weight         bigint unsigned not null,

    # Duration validator
    start_time     datetime        not null,
    end_time       datetime        not null,

    # Default allychain validator
    destination    varchar(50)     not null,
    shares         int unsigned    not null,

    # Allychain validator
    allychain_id      varchar(50)     not null
);
create index pvm_validators_node_id_idx ON pvm_validators (node_id);
create index pvm_validators_allychain_id_idx ON pvm_validators (allychain_id);
create unique index pvm_validators_tx_id_idx ON pvm_validators (transaction_id);

##
## Chains
##

create table pvm_chains
(
    id           varchar(50)      not null primary key,
    network_id   int unsigned     not null,
    allychain_id    varchar(50)      not null,

    name         varchar(255)     not null,
    vm_id        varchar(50)      not null,
    genesis_data varbinary(16384) not null
);

create table pvm_chains_control_signatures
(
    chain_id  varchar(50) not null,
    signature binary(65)  not null
);
create index pvm_chains_control_sigs_chain_id_idx ON pvm_chains_control_signatures (chain_id);


create table pvm_chains_fx_ids
(
    chain_id varchar(50) not null,
    fx_id    varchar(50) not null

);
create index pvm_chains_fix_ids_chain_id_idx ON pvm_chains_fx_ids (chain_id);
create unique index pvm_chains_fx_ids_chain_id_fix_id_idx ON pvm_chains_fx_ids (chain_id, fx_id);

