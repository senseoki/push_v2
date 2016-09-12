package service

const (
	Select_PushMessage = `
                        select   service_cd
                                ,push_type
                                ,msg_seq
                                ,msg_type
                                ,send_msg
                                ,ifnull(send_status,'') as send_status
                                ,ifnull(send_hope_dt,'') as send_hope_dt
                                ,ifnull(img_title,'') as img_title
                                ,ifnull(img_file_path,'') as img_file_path
                                ,ifnull(link_url,'') as link_url
                                ,ifnull(total_cnt,'') as total_cnt
                                ,ifnull(ios_send_cnt,'') as ios_send_cnt
                                ,ifnull(android_send_cnt,'') as android_send_cnt
                                ,ifnull(reg_dt,'') as reg_dt
                                ,ifnull(send_start_dt,'') as send_start_dt
                                ,ifnull(send_end_dt,'') as send_end_dt
                                ,ifnull(del_yn,'') as del_yn
                                ,ifnull(del_dt,'') as del_dt
                                ,ifnull(test_yn,'') as test_yn
                        from push_message
						where send_status = '1001'
						limit 1 `

	Select_PushTarget = `
						select
                                t1.push_target_seq
                               ,t1.user_key
                               ,ifnull(t1.mobile,'') as mobile
                               ,t1.os_cd
                               ,t1.push_token
                               ,t1.reg_dt
						from push_target t1
                        left outer join push_target_status t2
						on t1.push_target_seq = t2.push_target_seq
						where t2.push_target_seq is null
                        and t1.service_cd = ?
                        and t1.push_type = ?
                        and t1.msg_seq = ?
						limit ? `

	Update_PushMessageSendStatus1002 = `
                          update push_message set send_status = ?, send_start_dt = now()
                          where service_cd = ? 
                          and push_type = ?
                          and msg_seq = ? `
	Update_PushMessageSendStatus1003 = `
                          update push_message set send_status = ?, send_end_dt = now()
                          where service_cd = ? 
                          and push_type = ?
                          and msg_seq = ? `

	Insert_PushMessageLog = `
							insert push_message_log
                            (service_cd,
                             push_type,
                             msg_seq,
                             msg_type,
                             send_msg,
                             send_status,
                             send_hope_dt,
                             img_title,
                             img_file_path,
                             link_url,
                             total_cnt,
                             ios_send_cnt,
                             android_send_cnt,
                             reg_dt,
                             send_start_dt,
                             send_end_dt,
                             del_yn,
                             del_dt,
                             test_yn)
							select  service_cd,
                                    push_type,
                                    msg_seq,
                                    msg_type,
                                    send_msg,
                                    send_status,
                                    send_hope_dt,
                                    ifnull(img_title,''),
                                    ifnull(img_file_path,''),
                                    ifnull(link_url,''),
                                    ifnull(total_cnt,''),
                                    ifnull(ios_send_cnt,''),
                                    ifnull(android_send_cnt,''),
                                    reg_dt,
                                    send_start_dt,
                                    send_end_dt,
                                    ifnull(del_yn,'N'),
                                    del_dt,
                                    ifnull(test_yn,'N')
                            from push_message
							where service_cd = ? 
                            and push_type = ?
                            and msg_seq = ? `

	Insert_PushTargetStatus = `insert into push_target_status (push_target_seq,send_status) values `
)
