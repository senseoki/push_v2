package service

const (
	SelectPushTargetStatus = `
        select  push_target_seq
               ,send_status
               ,reg_dt
        from push_target_status
        where send_status in (0,4)
        order by push_target_seq
        limit 1000 `

	InsertPushTargetStatusLog = `
         insert into push_target_status_log 
         select  push_target_seq
                ,send_status
                ,reg_dt
         from push_target_status
         where send_status in (0,4)
         and push_target_seq in `

	InsertPushTargetLog = `
        insert into push_target_log
        select   push_target_seq
                ,service_cd
                ,push_type
                ,msg_seq
                ,user_key
                ,mobile
                ,os_cd
                ,push_token
                ,reg_dt
        from push_target
        where push_target_seq in `

	DeletePushTarget = `
        delete from push_target where push_target_seq in  `

	DeletePushTargetStatus = `
        delete from push_target_status where push_target_seq in `

	DeletePushMessage = `
	delete from push_message where (service_cd, push_type, msg_seq) in (
                select
		        a.service_cd
                        ,a.push_type
		        ,a.msg_seq
                from (
		        select  service_cd
			        ,push_type
				,msg_seq
			from push_message
			where send_status = '1003'
			and send_end_dt <= date_add(now(), interval -1 day)
			order by send_end_dt
			limit 1000
                ) a
        )`
)
