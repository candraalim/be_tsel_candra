CREATE SCHEMA referral;

CREATE TABLE IF NOT EXISTS referral.referral_code
(
    id SERIAL,
    msisdn character varying(20) NOT NULL,
    code character varying(20) NOT NULL,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    status integer NOT NULL DEFAULT 1,
    CONSTRAINT referral_code_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX referral_code_msisdn_idx ON referral.referral_code(msisdn);
CREATE UNIQUE INDEX referral_code_idx ON referral.referral_code(code);


CREATE TABLE IF NOT EXISTS referral.reward
(
    id SERIAL,
    total_referral integer NOT NULL,
    reward_description character varying(100) NOT NULL,
    status integer NOT NULL DEFAULT 1,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    updated_date timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT reward_pkey PRIMARY KEY (id)
);
INSERT INTO referral.reward (total_referral, reward_description) VALUES
(1, 'bonus 2 GB'),
(5, 'bonus 12 GB'),
(6, 'bonus 20 GB');


CREATE TABLE IF NOT EXISTS referral.referral_history
(
    id SERIAL,
    msisdn character varying(20) NOT NULL,
    code character varying(20) NOT NULL,
    msisdn_referee character varying(20) NOT NULL,
    referral_date character varying(10) NOT NULL,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT referral_history_pkey PRIMARY KEY (id)
);
CREATE INDEX referral_history_msisdn_idx ON referral.referral_history(msisdn);
CREATE INDEX referral_history_code_idx ON referral.referral_history(code);
CREATE INDEX referral_history_msisdn_referee_idx ON referral.referral_history(msisdn_referee);
