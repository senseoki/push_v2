package service

const (
	SelectPushBatchMsg = `
                        select t2.push_target_seq as push_target_seq
							,t3.max_send_status as max_send_status
							,t3.scheduler_work as scheduler_work
							,ifnull(t2.service_cd, '') as service_cd
                            ,ifnull(t2.push_type, '') as push_type
							,ifnull(t2.msg_seq, '') as msg_seq
							,ifnull(t1.msg_type, '') as msg_type
							,ifnull(t1.send_msg, '') as send_msg
							,ifnull(t1.img_title, '') as img_title
							,ifnull(t1.img_file_path, '') as img_file_path
							,ifnull(t1.link_url, '') as link_url
							,ifnull(t2.user_key, '') as user_key
							,ifnull(t2.mobile, '') as mobile
							,ifnull(t2.os_cd, '') as os_cd
							,ifnull(t2.push_token, '') as push_token
							,ifnull(t1.test_yn, '') as test_yn
						from push_target t2
							inner join (
								select
									 push_target_seq
									,ifnull(MIN(send_status), '') as min_send_status
									,ifnull(MAX(send_status), '') as max_send_status 		             
									,if (SEC_TO_TIME(unix_timestamp(NOW())-unix_timestamp(MAX(reg_dt)))+0 >= 300, 'r', '') as scheduler_work
									,reg_dt 
								from push_target_status		
								group by push_target_seq
								having min_send_status != 0 and max_send_status < 4
								and scheduler_work = 'r'
								limit ?
							) t3
							on t2.push_target_seq = t3.push_target_seq   
							inner join push_message t1
							on t1.msg_seq = t2.msg_seq
							and t1.service_cd = t2.service_cd
							and t1.push_type = t2.push_type`
)
