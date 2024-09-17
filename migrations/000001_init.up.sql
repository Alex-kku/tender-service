CREATE TABLE IF NOT EXISTS employee
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP

);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE organization
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_id         UUID REFERENCES employee (id) ON DELETE CASCADE
);

INSERT INTO employee(id, username, first_name, last_name, created_at, updated_at)
VALUES ('550e8400-e29b-41d4-a716-446655440001', 'user1', 'First1', 'Last1', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440002', 'user2', 'First2', 'Last2', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440003', 'user3', 'First3', 'Last3', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440004', 'user4', 'First4', 'Last4', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440005', 'user5', 'First5', 'Last5', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440006', 'user6', 'First6', 'Last6', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440007', 'user7', 'First7', 'Last7', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440008', 'user8', 'First8', 'Last8', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440009', 'user9', 'First9', 'Last9', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544000a', 'user10', 'First10', 'Last10', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544000b', 'user11', 'First11', 'Last11', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544000c', 'user12', 'First12', 'Last12', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544000d', 'user13', 'First13', 'Last13', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544000e', 'user14', 'First14', 'Last14', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544000f', 'user15', 'First15', 'Last15', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440010', 'user16', 'First16', 'Last16', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440011', 'user17', 'First17', 'Last17', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440012', 'user18', 'First18', 'Last18', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440013', 'user19', 'First19', 'Last19', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440014', 'user20', 'First20', 'Last20', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440015', 'user21', 'First21', 'Last21', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440016', 'user22', 'First22', 'Last22', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440017', 'user23', 'First23', 'Last23', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440018', 'user24', 'First24', 'Last24', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-446655440019', 'user25', 'First25', 'Last25', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544001a', 'user26', 'First26', 'Last26', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544001b', 'user27', 'First27', 'Last27', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544001c', 'user28', 'First28', 'Last28', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544001d', 'user29', 'First29', 'Last29', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554'),
       ('550e8400-e29b-41d4-a716-44665544001e', 'user30', 'First30', 'Last30', '2024-09-14 14:06:05.045554',
        '2024-09-14 14:06:05.045554');

INSERT INTO organization (id, name, description, type, created_at, updated_at)
VALUES ('550e8400-e29b-41d4-a716-446655440020', 'Organization 1', 'Description 1', 'LLC', '2024-09-14 14:06:05.062332',
        '2024-09-14 14:06:05.062332'),
       ('550e8400-e29b-41d4-a716-446655440021', 'Organization 2', 'Description 2', 'IE', '2024-09-14 14:06:05.062332',
        '2024-09-14 14:06:05.062332'),
       ('550e8400-e29b-41d4-a716-446655440022', 'Organization 3', 'Description 3', 'JSC', '2024-09-14 14:06:05.062332',
        '2024-09-14 14:06:05.062332'),
       ('550e8400-e29b-41d4-a716-446655440023', 'Organization 4', 'Description 4', 'LLC', '2024-09-14 14:06:05.062332',
        '2024-09-14 14:06:05.062332');

INSERT INTO organization_responsible (id, organization_id, user_id)
VALUES ('550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440001'),
       ('550e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440032', '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440003'),
       ('550e8400-e29b-41d4-a716-446655440033', '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440004'),
       ('550e8400-e29b-41d4-a716-446655440034', '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440005'),
       ('550e8400-e29b-41d4-a716-446655440035', '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440006'),
       ('550e8400-e29b-41d4-a716-446655440036', '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440007'),
       ('550e8400-e29b-41d4-a716-446655440037', '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440008'),
       ('550e8400-e29b-41d4-a716-446655440038', '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440009'),
       ('550e8400-e29b-41d4-a716-446655440039', '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-44665544000a'),
       ('550e8400-e29b-41d4-a716-44665544003a', '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-44665544000b'),
       ('550e8400-e29b-41d4-a716-44665544003b', '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-44665544000c');

CREATE TYPE tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');

CREATE TABLE IF NOT EXISTS tenders
(
    id              UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    name            VARCHAR(100)        NOT NULL,
    description     VARCHAR(500)        NOT NULL,
    service_type    tender_service_type NOT NULL,
    status          tender_status       NOT NULL DEFAULT 'Created',
    organization_id UUID                NOT NULL,
    version         INT                 NOT NULL DEFAULT 1,
    created_at      TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP,
    creator_id      UUID                NOT NULL,
    updated_at      TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tenders_old_version
(
    id              UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    tender_id       UUID                NOT NULL,
    name            VARCHAR(100)        NOT NULL,
    description     VARCHAR(500)        NOT NULL,
    service_type    tender_service_type NOT NULL,
    status          tender_status       NOT NULL DEFAULT 'Created',
    organization_id UUID                NOT NULL,
    version         INT                 NOT NULL DEFAULT 1,
    created_at      TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP,
    creator_id      UUID                NOT NULL,
    updated_at      TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled');
CREATE TYPE bid_authorType AS ENUM ('Organization', 'User');
CREATE TYPE bid_decision AS ENUM ('Approved', 'Rejected');

CREATE TABLE IF NOT EXISTS bids
(
    id          UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    name        VARCHAR(100)   NOT NULL,
    description VARCHAR(500)   NOT NULL,
    status      bid_status     NOT NULL DEFAULT 'Created',
    tender_id   UUID           NOT NULL,
    author_type bid_authorType NOT NULL,
    author_id   UUID           NOT NULL,
    version     INT            NOT NULL DEFAULT 1,
    created_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    decision    bid_decision,
    FOREIGN KEY (tender_id) REFERENCES tenders (id)
);

CREATE TABLE IF NOT EXISTS bids_old_version
(
    id          UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    bid_id      UUID           NOT NULL,
    name        VARCHAR(100)   NOT NULL,
    description VARCHAR(500)   NOT NULL,
    status      bid_status     NOT NULL DEFAULT 'Created',
    tender_id   UUID           NOT NULL,
    author_type bid_authorType NOT NULL,
    author_id   UUID           NOT NULL,
    version     INT            NOT NULL DEFAULT 1,
    created_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    decision    bid_decision,
    FOREIGN KEY (tender_id) REFERENCES tenders (id)
);