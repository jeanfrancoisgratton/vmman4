// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/db_mgt/init.go
// Original timestamp : 2026.10.02 15:40:22

package db_mgt

import ce "github.com/jeanfrancoisgratton/customError/v3"

/*

vmStates table : main table


CREATE TABLE IF NOT EXISTS public."vmStates"
(
    "vmsID" bigint NOT NULL,
    "vmsName" character varying(20) COLLATE pg_catalog."default" NOT NULL,
    "vmsHypervisor" character varying(20) COLLATE pg_catalog."default" NOT NULL,
    "vmsLastUp" time without time zone,
    "vmsLastDown" time without time zone,
    "vmsOS" character varying(30) COLLATE pg_catalog."default",
    CONSTRAINT vmstates_pkey PRIMARY KEY ("vmsID")
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."vmStates"
    OWNER to vmman;
GRANT ALL ON TABLE public."vmStates" TO vmman;

*/

func DBinit() *ce.CustomError {
	return nil
}
