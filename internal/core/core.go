package core

import (
	"time"

)

type DatabaseRDS struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	Schema				string `json:"schema"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	Db_timeout			int	`json:"db_timeout"`
	Postgres_Driver		string `json:"postgres_driver"`
}

type WorkerAppServer struct {
	InfoPod 	*InfoPod 		`json:"info_pod"`
}

type InfoPod struct {
	PodName				string `json:"pod_name"`
	ApiVersion			string `json:"version"`
	OSPID				string `json:"os_pid"`
	IPAddress			string `json:"ip_address"`
	AvailabilityZone 	string `json:"availabilityZone"`
	Database			*DatabaseRDS
	Kafka				*KafkaConfig
}

type Balance struct {
	ID				int		`json:"id,omitempty"`
	AccountID		string	`json:"account_id,omitempty"`
	PersonID		string  `json:"person_id,omitempty"`
	Currency		string  `json:"currency,omitempty"`
	Amount			float64 `json:"amount,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
	TenantID		string  `json:"tenant_id,omitempty"`
	UserLastUpdate	*string  `json:"user_last_update,omitempty"`
}