// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/db_mgt/types.go
// Original timestamp : 2026.10.02 15:18:19

package db_mgt

// This structure holds the basic DB config

type DBenvironmentStruct struct {
	Host     string `json:"dbhost"`
	Port     uint   `json:"dbport"`
	DbName   string `json:"dbname"`
	User     string `json:"dbuser"`
	Password string `json:"dbpassword"`
}

const createVMStatesTableSQL = `
CREATE TABLE IF NOT EXISTS public.vm_states (
    vms_id         BIGINT PRIMARY KEY,
    vms_name       VARCHAR(20) NOT NULL,
    vms_hypervisor VARCHAR(20) NOT NULL,
    vms_last_up    TIMESTAMPTZ,
    vms_last_down  TIMESTAMPTZ,
    vms_date_created    TIMESTAMPTZ,
    vms_date_deleted  TIMESTAMPTZ,
    vms_os         VARCHAR(30)
)`

const setVMStatesOwnerSQL = "ALTER TABLE public.vm_states OWNER TO vmman"

const grantVMStatesSQL = "GRANT ALL ON TABLE public.vm_states TO vmman"
