package service

const (
	// SelectPushTargetRealtimeStatus ...
	// limit 10000 작은 단위로 가져와야 성능이 좋다.
	SelectPushTargetRealtimeStatus = `
        select  push_target_seq
                ,send_status
                ,reg_dt
        from push_target_realtime_status
        where send_status in (0,4)
        order by push_target_seq
        limit 10000 `

	InsertPushTargetRealtimeStatusLog = `
        insert into push_target_realtime_status_log
        select *
        from push_target_realtime_status
        where send_status in (0,4)
        and push_target_seq in `

	InsertPushTargetRealtimeLog = `
        insert into push_target_realtime_log
        select *
        from push_target_realtime
        where push_target_seq in `

	DeletePushTargetRealtimeStatus = `
        delete from push_target_realtime_status where push_target_seq in `

	DeletePushTargetRealtime = `
        delete from push_target_realtime where push_target_seq in `
)
