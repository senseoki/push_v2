package service

const (
	SelectPushTargetStatus = `
        select  push_target_seq
               ,send_status
               ,reg_dt
        from push_target_status
        where send_status in (0,4)
        order by push_target_seq
        limit 10000 `

	InsertPushTargetStatusLog = `
         insert into push_target_status_log 
         select *
         from push_target_status
         where send_status in (0,4)
         and push_target_seq in `

	InsertPushTargetLog = `
        insert into push_target_log
        select *
        from push_target
        where push_target_seq in `

	DeletePushTarget = `
        delete from push_target where push_target_seq in  `

	DeletePushTargetStatus = `
        delete from push_target_status where push_target_seq in `
)
