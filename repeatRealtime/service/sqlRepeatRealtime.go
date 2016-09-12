package service

const (
	Select_RepeatRealtimeMsg = `
                        select   t1.push_target_seq as push_target_seq 
                                ,t2.max_send_status as max_send_status
                                ,t2.scheduler_work as scheduler_work
                                ,ifnull(t1.service_cd, '') as service_cd
                                ,ifnull(t1.push_type, '') as push_type
                                ,ifnull(t1.msg_seq, '') as msg_seq
                                ,ifnull(t1.msg_type, '') as msg_type
                                ,ifnull(t1.send_msg, '') as send_msg
                                ,ifnull(t1.send_hope_dt, '') as send_hope_dt
                                ,ifnull(t1.img_title, '') as img_title
                                ,ifnull(t1.img_file_path, '') as img_file_path
                                ,ifnull(t1.link_url, '') as link_url
                                ,ifnull(t1.user_key, '') as user_key
                                ,ifnull(t1.mobile, '') as mobile
                                ,ifnull(t1.os_cd, '') as os_cd
                                ,ifnull(t1.push_token, '') as push_token
                                ,ifnull(t1.test_yn, '') as test_yn
                        from push_target_realtime t1
                                inner join (
                                select
                                         push_target_seq
                                        ,ifnull(min(send_status), '') as min_send_status
                                        ,ifnull(max(send_status), '') as max_send_status 									 
                                        ,if (SEC_TO_TIME(unix_timestamp(now())-unix_timestamp(max(reg_dt)))+0 >= 300, 'r', '') as scheduler_work
                                        ,reg_dt 
                                from push_target_realtime_status		
                                group by push_target_seq
                                having min_send_status != 0 and max_send_status < 4
                                and scheduler_work = 'r'
                                limit ?
                            ) t2
                            on t1.push_target_seq = t2.push_target_seq `
)
