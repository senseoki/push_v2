package service

const (
	Select_PushTargetRealtime = `
						select 
							 t1.push_target_seq	as push_target_seq
							,ifnull(t1.service_cd, '')	as service_cd
							,ifnull(t1.push_type, '')	as push_type
							,ifnull(t1.msg_seq, '')	as msg_seq
							,ifnull(t1.msg_type, '')	as msg_type
							,ifnull(t1.send_msg, '')	as send_msg
							,ifnull(t1.send_hope_dt, '')	as send_hope_dt
							,ifnull(t1.img_title, '')	as img_title
							,ifnull(t1.img_file_path, '')	as img_file_path
							,ifnull(t1.link_url, '')	as link_url
							,ifnull(t1.user_key, '')	as user_key
							,ifnull(t1.mobile, '')	as mobile
							,ifnull(t1.os_cd, '')	as os_cd
							,ifnull(t1.push_token, '')	as push_token
							,ifnull(t1.test_yn, '')	as test_yn
						FROM push_target_realtime t1
						LEFT OUTER JOIN push_target_realtime_status t2    
						ON t1.push_target_seq = t2.push_target_seq   
						WHERE t2.push_target_seq is null
                        LIMIT ? `

	Insert_PushTargetRealtimeStatus = `insert push_target_realtime_status (push_target_seq, send_status) values `
)
